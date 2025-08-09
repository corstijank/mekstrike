/**
 * Battlefield service layer
 * Handles battlefield-specific business logic and board operations
 */

import { battlefieldAPI, gameAPI } from './gameAPI.js';
import { getHexCenter, colToCenterX, rowToCenterY, hexWidth, hexHeight } from '../utils/coordinates.js';

/**
 * Calculate board dimensions with padding
 */
export function calculateBoardDimensions(cols, rows) {
    return {
        width: (cols + 1) * hexWidth,
        height: (rows + 1) * hexHeight,
    };
}

/**
 * Get board data with calculated dimensions
 */
export async function getBoardWithDimensions(gameId) {
    try {
        const boardData = await gameAPI.getBoard(gameId);
        const cols = parseInt(boardData.cols);
        const rows = parseInt(boardData.rows);
        const dimensions = calculateBoardDimensions(cols, rows);
        
        return {
            ...boardData,
            cols,
            rows,
            ...dimensions,
        };
    } catch (error) {
        console.error(`Failed to get board data for game ${gameId}:`, error);
        throw error;
    }
}

/**
 * Check if coordinates are valid for the board
 */
export function areValidCoordinates(x, y, cols, rows) {
    return x >= 0 && x < cols && y >= 0 && y < rows;
}

/**
 * Get hex center coordinates for scrolling
 */
export function getScrollPosition(col, row, viewportWidth, viewportHeight) {
    const x = colToCenterX(col);
    const y = rowToCenterY(row, col);
    
    return {
        scrollX: Math.max(0, x - (viewportWidth / 2) + 50),
        scrollY: Math.max(0, y - (viewportHeight / 2) + 50),
    };
}

/**
 * Filter movement highlights for human players only
 */
export function filterMovementHighlights(availableActions) {
    if (!availableActions?.UnitOwner || availableActions.UnitOwner === 'CPU') {
        return [];
    }
    return [...(availableActions.AllowedCoordinates || [])];
}

/**
 * Get active unit position from available actions
 */
export async function getActiveUnitPosition(gameId, availableActions) {
    if (!availableActions?.CurrentUnitID) {
        return null;
    }
    
    try {
        const unitData = await gameAPI.getUnit(gameId, availableActions.CurrentUnitID);
        return {
            x: unitData.location.position.x,
            y: unitData.location.position.y,
        };
    } catch (error) {
        console.error(`Failed to get active unit position:`, error);
        return null;
    }
}

/**
 * Check if position matches coordinates
 */
export function positionMatches(position1, position2) {
    if (!position1 || !position2) return false;
    return position1.x === position2.x && position1.y === position2.y;
}

/**
 * Get terrain type for a cell
 */
export function getTerrainType(cell) {
    return cell?.terrain || 'default';
}

/**
 * Check if cell has unit
 */
export function cellHasUnit(cell, units = []) {
    if (!cell || !units.length) return false;
    
    return units.some(unitId => {
        // This would need to be enhanced with unit position data
        // For now, returning false as units are handled separately
        return false;
    });
}