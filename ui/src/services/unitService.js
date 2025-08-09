/**
 * Unit service layer
 * Handles unit-specific business logic and data transformations
 */

import { unitAPI, gameAPI } from './gameAPI.js';

/**
 * Calculate health percentage for unit
 */
export function calculateHealthPercent(currentStruct, maxStruct = null) {
    // If no maxStruct provided, use current as max (for cases where API doesn't provide max)
    const max = maxStruct || currentStruct || 1;
    if (max <= 0) return 100;
    return (currentStruct / max) * 100;
}

/**
 * Get health bar color based on percentage
 */
export function getHealthBarColor(percent) {
    if (percent > 66) return '#00ff00'; // Green
    if (percent > 33) return '#ffff00'; // Yellow
    return '#ff0000'; // Red
}

/**
 * Format unit name for display (first word of model)
 */
export function formatUnitDisplayName(model) {
    if (!model) return 'Unknown';
    return model.split(' ')[0];
}

/**
 * Get unit data with health calculations
 */
export async function getUnitWithHealth(unitId) {
    try {
        const unitData = await unitAPI.getUnitData(unitId);
        const healthPercent = calculateHealthPercent(
            unitData.stats.struct, 
            unitData.stats.struct // Note: API doesn't seem to have maxStruct, using current as max
        );
        
        return {
            ...unitData,
            healthPercent,
            healthColor: getHealthBarColor(healthPercent),
            displayName: formatUnitDisplayName(unitData.stats.model),
        };
    } catch (error) {
        console.error(`Failed to get unit data for ${unitId}:`, error);
        throw error;
    }
}

/**
 * Get enhanced unit data from game context
 */
export async function getGameUnit(gameId, unitId) {
    try {
        const unitData = await gameAPI.getUnit(gameId, unitId);
        return {
            ...unitData,
            displayName: formatUnitDisplayName(unitData.stats?.model),
        };
    } catch (error) {
        console.error(`Failed to get game unit data for ${unitId}:`, error);
        throw error;
    }
}

/**
 * Check if unit is active in current turn
 */
export function isUnitActive(unitId, availableActions) {
    return availableActions?.CurrentUnitID === unitId;
}

/**
 * Check if unit belongs to human player (not CPU)
 */
export function isHumanUnit(availableActions) {
    return availableActions?.UnitOwner && availableActions.UnitOwner !== 'CPU';
}

/**
 * Get unit position from unit data
 */
export function getUnitPosition(unitData) {
    return unitData?.location?.position;
}

/**
 * Format unit position for display
 */
export function formatUnitPosition(position) {
    if (!position) return 'Unknown position';
    return `(${position.x}, ${position.y})`;
}