/**
 * WebSocket service for real-time game updates
 * Handles connection management and event routing to stores
 */

import { get } from 'svelte/store';
import { gameMessages, refreshGameState, refreshAvailableActions } from '../stores/gameStore.js';
import { refreshBoard } from '../stores/battlefieldStore.js';

class WebSocketService {
    constructor() {
        this.socket = null;
        this.gameId = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000; // Start with 1 second
        this.isConnecting = false;
        this.shouldReconnect = true;
        this.onViewportPreserve = null; // Callback to preserve viewport before refresh
        this.onViewportRestore = null;  // Callback to restore viewport after refresh
    }

    /**
     * Connect to WebSocket for a specific game
     */
    connect(gameId, viewportPreserveFn = null, viewportRestoreFn = null) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN && this.gameId === gameId) {
            console.log('WebSocket already connected for game:', gameId);
            return;
        }

        if (this.isConnecting) {
            console.log('WebSocket connection already in progress');
            return;
        }

        this.gameId = gameId;
        this.shouldReconnect = true;
        this.isConnecting = true;
        this.onViewportPreserve = viewportPreserveFn;
        this.onViewportRestore = viewportRestoreFn;

        // Close existing connection if any
        this.disconnect();

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/mekstrike/api/gamemaster/ws/games/${gameId}`;

        console.log('Connecting to WebSocket:', wsUrl);

        try {
            this.socket = new WebSocket(wsUrl);
            
            this.socket.onopen = (event) => {
                console.log('WebSocket connected for game:', gameId);
                this.isConnecting = false;
                this.reconnectAttempts = 0;
                this.reconnectDelay = 1000;
            };

            this.socket.onmessage = async (event) => {
                await this.handleMessage(event.data);
            };

            this.socket.onclose = (event) => {
                console.log('WebSocket connection closed:', event.code, event.reason);
                this.isConnecting = false;
                
                if (this.shouldReconnect && this.reconnectAttempts < this.maxReconnectAttempts) {
                    this.scheduleReconnect();
                } else {
                    console.log('WebSocket reconnection stopped');
                }
            };

            this.socket.onerror = (event) => {
                console.error('WebSocket error:', event);
                this.isConnecting = false;
            };

        } catch (error) {
            console.error('Failed to create WebSocket connection:', error);
            this.isConnecting = false;
            this.scheduleReconnect();
        }
    }

    /**
     * Handle incoming WebSocket messages (CloudEvents)
     */
    async handleMessage(data) {
        try {
            const cloudEvent = JSON.parse(data);
            
            // Add message to game logs immediately
            this.addCloudEventMessage(cloudEvent);

            // Preserve viewport position before refresh
            let scrollPosition = null;
            if (this.onViewportPreserve) {
                scrollPosition = this.onViewportPreserve();
            }

            // Refresh game state and board based on CloudEvent
            await this.handleCloudEventUpdate(cloudEvent, scrollPosition);

        } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
        }
    }

    /**
     * Add CloudEvent as a formatted message
     */
    addCloudEventMessage(cloudEvent) {
        let parsedData = cloudEvent.data;
        try {
            parsedData = JSON.parse(cloudEvent.data);
        } catch (e) {
            // Keep as string if not JSON
        }

        const message = {
            type: cloudEvent.type,
            message: this.formatEventMessage(parsedData, cloudEvent.type),
            timestamp: new Date(cloudEvent.time || Date.now()),
            rawData: parsedData,
        };

        // Add to gameMessages store
        gameMessages.update(messages => [...messages, message]);
    }

    /**
     * Handle game state updates based on CloudEvent
     */
    async handleCloudEventUpdate(cloudEvent, scrollPosition) {
        // Extract GameId from the CloudEvent data
        let gameId = this.gameId;
        if (!gameId && cloudEvent.data) {
            try {
                const eventData = JSON.parse(cloudEvent.data);
                gameId = eventData.GameId;
            } catch (e) {
                console.error('Failed to parse CloudEvent data for GameId:', e);
            }
        }
        
        if (!gameId) {
            return;
        }

        try {
            // Refresh game state, available actions, and board data
            await Promise.all([
                refreshGameState(gameId),
                refreshAvailableActions(gameId),
                refreshBoard(gameId)
            ]);

            // Restore viewport position after refresh
            if (scrollPosition && this.onViewportRestore) {
                // Small delay to ensure DOM is updated
                setTimeout(() => {
                    this.onViewportRestore(scrollPosition);
                }, 50);
            }

        } catch (error) {
            console.error('Failed to refresh game state after WebSocket event:', error);
        }
    }

    /**
     * Format event message for display (reuse logic from gameStore)
     */
    formatEventMessage(data, type) {
        if (typeof data !== 'object') {
            return data;
        }

        // Handle different event types
        switch (type) {
            case 'com.dapr.event.sent':
                if (data.Phase === 'Movement') return this.formatMovementEvent(data);
                if (data.Phase === 'Combat') return this.formatAttackEvent(data);
                if (data.Phase === 'End') return this.formatEndPhaseEvent(data);
                return `Unit ${data.UnitId || 'unknown'} completed ${data.Phase || 'unknown'} phase`;

            case 'unit.movement.completed':
            case 'unit-movement-completed':
                return this.formatMovementEvent(data);

            case 'unit.attack.completed':
            case 'unit-attack-completed':
                return this.formatAttackEvent(data);

            case 'unit.end.completed':
            case 'unit-end-phase-completed':
                return this.formatEndPhaseEvent(data);

            case 'movement':
            case 'combat':
            case 'system':
                return data;

            default:
                return typeof data === 'string' ? data : JSON.stringify(data);
        }
    }

    formatMovementEvent(data) {
        if (data.Unit && data.SourceLocation && data.TargetLocation) {
            const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
            const player = data.Unit.Owner || 'Unknown player';
            const source = `(${data.SourceLocation.x}, ${data.SourceLocation.y})`;
            const target = `(${data.TargetLocation.x}, ${data.TargetLocation.y})`;
            return `${player} moved ${unitName} from ${source} to ${target}`;
        }
        return `Unit ${data.UnitId || 'unknown'} completed movement phase`;
    }

    formatAttackEvent(data) {
        if (data.Unit && data.SourceLocation && data.TargetId) {
            const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
            const player = data.Unit.Owner || 'Unknown player';
            const source = `(${data.SourceLocation.x}, ${data.SourceLocation.y})`;
            return `${player}'s ${unitName} at ${source} attacked target ${data.TargetId}`;
        }
        return `Unit ${data.UnitId || 'unknown'} completed combat phase`;
    }

    formatEndPhaseEvent(data) {
        if (data.Unit) {
            const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
            const player = data.Unit.Owner || 'Unknown player';
            const location = data.SourceLocation ? 
                `(${data.SourceLocation.x}, ${data.SourceLocation.y})` : 'unknown location';
            return `${player}'s ${unitName} at ${location} ended its turn`;
        }
        return `Unit ${data.UnitId || 'unknown'} completed end phase`;
    }

    /**
     * Schedule reconnection with exponential backoff
     */
    scheduleReconnect() {
        if (!this.shouldReconnect || this.reconnectAttempts >= this.maxReconnectAttempts) {
            return;
        }

        this.reconnectAttempts++;
        const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
        
        console.log(`Scheduling WebSocket reconnection in ${delay}ms (attempt ${this.reconnectAttempts})`);

        setTimeout(() => {
            if (this.shouldReconnect && this.gameId) {
                this.connect(this.gameId, this.onViewportPreserve, this.onViewportRestore);
            }
        }, delay);
    }

    /**
     * Disconnect WebSocket
     */
    disconnect() {
        this.shouldReconnect = false;
        
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }

        this.gameId = null;
        this.isConnecting = false;
        this.reconnectAttempts = 0;
        this.onViewportPreserve = null;
        this.onViewportRestore = null;
    }

    /**
     * Check if WebSocket is connected
     */
    isConnected() {
        return this.socket && this.socket.readyState === WebSocket.OPEN;
    }

    /**
     * Get connection state
     */
    getConnectionState() {
        if (!this.socket) return 'disconnected';
        
        switch (this.socket.readyState) {
            case WebSocket.CONNECTING: return 'connecting';
            case WebSocket.OPEN: return 'connected';
            case WebSocket.CLOSING: return 'closing';
            case WebSocket.CLOSED: return 'closed';
            default: return 'unknown';
        }
    }
}

// Export singleton instance
export const websocketService = new WebSocketService();