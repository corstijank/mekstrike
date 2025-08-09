/**
 * Unified game state management
 * Central store for game state, turn management, and polling coordination
 */

import { writable, derived, get } from 'svelte/store';
import { gameAPI } from '../services/gameAPI.js';

// Core game state
export const gameState = writable({
    gameId: null,
    PlayerAUnits: [],
    PlayerBUnits: [],
    CurrentRound: 0,
    CurrentPhase: 0,
    ActivePlayer: 0,
    loading: false,
    error: null,
});

// Available actions state  
export const availableActions = writable({
    UnitOwner: '',
    CurrentUnitID: '',
    ActionType: '',
    AllowedCoordinates: [],
    CurrentPhase: 0,
    loading: false,
    error: null,
});

// Game messages/logs
export const gameMessages = writable([]);

// Polling state
let gamePollingInterval = null;
let actionsPollingInterval = null;
let messagesPollingInterval = null;

/**
 * Derived stores for computed values
 */
export const currentGamePhase = derived(gameState, ($gameState) => {
    const phases = ['Movement', 'Combat', 'End'];
    return phases[$gameState.CurrentPhase] || 'Unknown';
});

export const currentPlayer = derived(gameState, ($gameState) => {
    const players = ['Player A', 'Player B'];
    return players[$gameState.ActivePlayer] || 'Unknown';
});

export const isHumanTurn = derived(availableActions, ($availableActions) => {
    return $availableActions.UnitOwner && $availableActions.UnitOwner !== 'CPU';
});

export const isCpuTurn = derived(availableActions, ($availableActions) => {
    return $availableActions.UnitOwner === 'CPU';
});

/**
 * Initialize game state and start polling
 */
export async function initializeGame(gameId) {
    if (!gameId) return;

    // Set game ID
    gameState.update(state => ({ ...state, gameId, loading: true, error: null }));

    try {
        // Initial data load
        await Promise.all([
            refreshGameState(gameId),
            refreshAvailableActions(gameId),
            refreshGameMessages(gameId),
        ]);

        // Start polling
        startPolling(gameId);
    } catch (error) {
        console.error('Failed to initialize game:', error);
        gameState.update(state => ({ ...state, error: error.message, loading: false }));
    }
}

/**
 * Refresh game state
 */
export async function refreshGameState(gameId) {
    if (!gameId) return;

    try {
        const data = await gameAPI.getGame(gameId);
        gameState.update(state => ({ 
            ...state, 
            PlayerAUnits: data.PlayerAUnits || state.PlayerAUnits,
            PlayerBUnits: data.PlayerBUnits || state.PlayerBUnits,
            CurrentRound: data.CurrentRound ?? state.CurrentRound,
            CurrentPhase: data.CurrentPhase ?? state.CurrentPhase,
            ActivePlayer: data.ActivePlayer ?? state.ActivePlayer,
            loading: false, 
            error: null 
        }));
    } catch (error) {
        console.error('Failed to refresh game state:', error);
        gameState.update(state => ({ ...state, error: error.message, loading: false }));
    }
}

/**
 * Refresh available actions
 */
export async function refreshAvailableActions(gameId) {
    if (!gameId) return;

    try {
        const data = await gameAPI.getAvailableActions(gameId);
        availableActions.update(state => ({ 
            ...state,
            UnitOwner: data.UnitOwner,
            CurrentUnitID: data.CurrentUnitID,
            ActionType: data.ActionType,
            AllowedCoordinates: data.AllowedCoordinates,
            CurrentPhase: data.CurrentPhase,
            loading: false, 
            error: null 
        }));
    } catch (error) {
        console.error('Failed to refresh available actions:', error);
        availableActions.update(state => ({ ...state, error: error.message, loading: false }));
    }
}

/**
 * Refresh game messages/logs
 */
export async function refreshGameMessages(gameId) {
    if (!gameId) return;

    try {
        const logs = await gameAPI.getGameLogs(gameId);
        const formattedMessages = logs.map(log => {
            let parsedData = log.message;
            try {
                parsedData = JSON.parse(log.message);
            } catch (e) {
                // Keep as string if not JSON
            }

            return {
                type: log.type,
                message: formatEventMessage(parsedData, log.type),
                timestamp: new Date(log.timestamp),
                rawData: parsedData,
            };
        });
        
        // Only update if the messages have actually changed
        gameMessages.update(currentMessages => {
            if (currentMessages.length !== formattedMessages.length) {
                return formattedMessages;
            }
            
            // Quick check if anything changed (compare last message timestamp)
            const lastCurrent = currentMessages[currentMessages.length - 1];
            const lastNew = formattedMessages[formattedMessages.length - 1];
            
            if (!lastCurrent || !lastNew || lastCurrent.timestamp.getTime() !== lastNew.timestamp.getTime()) {
                return formattedMessages;
            }
            
            return currentMessages; // No changes, keep current
        });
    } catch (error) {
        console.error('Failed to refresh game messages:', error);
    }
}

/**
 * Advance turn
 */
export async function advanceTurn() {
    const currentState = get(gameState);
    if (!currentState.gameId) return;

    try {
        await gameAPI.advanceTurn(currentState.gameId);
        // Refresh state after advancing turn
        await Promise.all([
            refreshGameState(currentState.gameId),
            refreshAvailableActions(currentState.gameId),
        ]);
    } catch (error) {
        console.error('Failed to advance turn:', error);
        throw error;
    }
}

/**
 * Start polling for game updates
 */
function startPolling(gameId) {
    if (!gameId) return;

    stopPolling(); // Clear any existing intervals

    // Poll game state every 3 seconds (reduced frequency)
    gamePollingInterval = setInterval(() => {
        refreshGameState(gameId);
    }, 3000);

    // Poll available actions every 3 seconds (reduced frequency)
    actionsPollingInterval = setInterval(() => {
        refreshAvailableActions(gameId);
    }, 3000);

    // Poll messages every 5 seconds (reduced frequency)
    messagesPollingInterval = setInterval(() => {
        refreshGameMessages(gameId);
    }, 5000);
}

/**
 * Stop all polling
 */
export function stopPolling() {
    if (gamePollingInterval) {
        clearInterval(gamePollingInterval);
        gamePollingInterval = null;
    }
    if (actionsPollingInterval) {
        clearInterval(actionsPollingInterval);
        actionsPollingInterval = null;
    }
    if (messagesPollingInterval) {
        clearInterval(messagesPollingInterval);
        messagesPollingInterval = null;
    }
}

/**
 * Add a new game message
 */
export function addGameMessage(type, message, timestamp = new Date()) {
    gameMessages.update(messages => [...messages, { type, message, timestamp }]);
}

/**
 * Format event message for display (extracted from GameMessageArea)
 */
function formatEventMessage(data, type) {
    if (typeof data !== 'object') {
        return data;
    }

    // Handle different event types
    switch (type) {
        case 'com.dapr.event.sent':
            if (data.Phase === 'Movement') return formatMovementEvent(data);
            if (data.Phase === 'Combat') return formatAttackEvent(data);
            if (data.Phase === 'End') return formatEndPhaseEvent(data);
            return `Unit ${data.UnitId || 'unknown'} completed ${data.Phase || 'unknown'} phase`;

        case 'unit.movement.completed':
        case 'unit-movement-completed':
            return formatMovementEvent(data);

        case 'unit.attack.completed':
        case 'unit-attack-completed':
            return formatAttackEvent(data);

        case 'unit.end.completed':
        case 'unit-end-phase-completed':
            return formatEndPhaseEvent(data);

        case 'movement':
        case 'combat':
        case 'system':
            return data;

        default:
            return typeof data === 'string' ? data : JSON.stringify(data);
    }
}

function formatMovementEvent(data) {
    if (data.Unit && data.SourceLocation && data.TargetLocation) {
        const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
        const player = data.Unit.Owner || 'Unknown player';
        const source = `(${data.SourceLocation.x}, ${data.SourceLocation.y})`;
        const target = `(${data.TargetLocation.x}, ${data.TargetLocation.y})`;
        return `${player} moved ${unitName} from ${source} to ${target}`;
    }
    return `Unit ${data.UnitId || 'unknown'} completed movement phase`;
}

function formatAttackEvent(data) {
    if (data.Unit && data.SourceLocation && data.TargetId) {
        const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
        const player = data.Unit.Owner || 'Unknown player';
        const source = `(${data.SourceLocation.x}, ${data.SourceLocation.y})`;
        return `${player}'s ${unitName} at ${source} attacked target ${data.TargetId}`;
    }
    return `Unit ${data.UnitId || 'unknown'} completed combat phase`;
}

function formatEndPhaseEvent(data) {
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
 * Cleanup function to call when component is destroyed
 */
export function cleanup() {
    stopPolling();
    gameState.set({
        gameId: null,
        PlayerAUnits: [],
        PlayerBUnits: [],
        CurrentRound: 0,
        CurrentPhase: 0,
        ActivePlayer: 0,
        loading: false,
        error: null,
    });
    availableActions.set({
        UnitOwner: '',
        CurrentUnitID: '',
        ActionType: '',
        AllowedCoordinates: [],
        CurrentPhase: 0,
        loading: false,
        error: null,
    });
    gameMessages.set([]);
}