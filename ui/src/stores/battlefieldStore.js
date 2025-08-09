/**
 * Battlefield state management
 * Handles board state, highlights, and battlefield-specific data
 */

import { writable, derived } from 'svelte/store';
import { gameAPI } from '../services/gameAPI.js';
import { getBoardWithDimensions, filterMovementHighlights, getActiveUnitPosition } from '../services/battlefieldService.js';

// Board state
export const boardState = writable({
    cells: [],
    cols: 0,
    rows: 0,
    width: 0,
    height: 0,
    loading: false,
    error: null,
});

// Movement highlights
export const movementHighlights = writable([]);

// Active unit position (for highlighting)
export const activeUnitPosition = writable(null);

// Scroll position for centering on units
export const scrollToPosition = writable(null);

/**
 * Derived store for board dimensions
 */
export const boardDimensions = derived(boardState, ($boardState) => ({
    width: $boardState.width,
    height: $boardState.height,
    cols: $boardState.cols,
    rows: $boardState.rows,
}));

/**
 * Initialize board for a game
 */
export async function initializeBoard(gameId) {
    if (!gameId) return;

    try {
        boardState.update(state => ({ ...state, loading: true, error: null }));
        const boardData = await getBoardWithDimensions(gameId);
        boardState.update(state => ({ 
            ...state, 
            ...boardData, 
            loading: false, 
            error: null 
        }));
    } catch (error) {
        console.error('Failed to initialize board:', error);
        boardState.update(state => ({ 
            ...state, 
            error: error.message, 
            loading: false 
        }));
    }
}

/**
 * Update movement highlights based on available actions
 */
export function updateMovementHighlights(availableActions) {
    const highlights = filterMovementHighlights(availableActions);
    movementHighlights.set(highlights);
}

/**
 * Update active unit position for highlighting
 */
export async function updateActiveUnitPosition(gameId, availableActions) {
    try {
        const position = await getActiveUnitPosition(gameId, availableActions);
        activeUnitPosition.set(position);
    } catch (error) {
        console.error('Failed to update active unit position:', error);
        activeUnitPosition.set(null);
    }
}

/**
 * Scroll to a specific unit position
 */
export function scrollToUnit(position) {
    if (position && typeof position.x === 'number' && typeof position.y === 'number') {
        // Check if we're already trying to scroll to avoid rapid position changes
        let currentScrollRequest = null;
        scrollToPosition.subscribe(current => currentScrollRequest = current)();
        
        // Only set scroll position if it's different from current request
        if (!currentScrollRequest || 
            currentScrollRequest.x !== position.x || 
            currentScrollRequest.y !== position.y) {
            scrollToPosition.set({
                x: position.x,
                y: position.y
            });
        }
    }
}

/**
 * Clear scroll position after scrolling is complete
 */
export function clearScrollPosition() {
    scrollToPosition.set(null);
}

/**
 * Clear all highlights
 */
export function clearHighlights() {
    movementHighlights.set([]);
    activeUnitPosition.set(null);
}

/**
 * Cleanup function
 */
export function cleanup() {
    boardState.set({
        cells: [],
        cols: 0,
        rows: 0,
        width: 0,
        height: 0,
        loading: false,
        error: null,
    });
    movementHighlights.set([]);
    activeUnitPosition.set(null);
    scrollToPosition.set(null);
}