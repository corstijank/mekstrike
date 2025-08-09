/**
 * Unit state management
 * Handles unit selection, unit data caching, and unit-specific state
 */

import { writable, derived } from 'svelte/store';
import { getUnitWithHealth, getGameUnit } from '../services/unitService.js';

// Selected unit state
export const selectedUnit = writable(null);
export const selectedUnitData = writable(null);

// Unit data cache (to avoid repeated API calls)
export const unitCache = writable(new Map());

// Loading states for units
export const unitLoadingStates = writable(new Map());

/**
 * Derived store for selected unit info
 */
export const selectedUnitInfo = derived(
    [selectedUnit, selectedUnitData], 
    ([$selectedUnit, $selectedUnitData]) => {
        if (!$selectedUnit || !$selectedUnitData) return null;
        
        return {
            unitId: $selectedUnit,
            data: $selectedUnitData,
            position: $selectedUnitData?.location?.position,
            displayName: $selectedUnitData?.stats?.model?.split(' ')[0] || 'Unknown',
            healthPercent: calculateHealthPercent($selectedUnitData),
        };
    }
);

/**
 * Select a unit and optionally scroll to it
 */
export function selectUnit(unitId, unitData = null, shouldScroll = false) {
    selectedUnit.set(unitId);
    selectedUnitData.set(unitData);
    
    // Import scrollToUnit from battlefieldStore to avoid circular dependency
    if (shouldScroll && unitData?.location?.position) {
        import('./battlefieldStore.js').then(({ scrollToUnit }) => {
            scrollToUnit({
                x: unitData.location.position.x,
                y: unitData.location.position.y,
            });
        });
    }
}

/**
 * Clear unit selection
 */
export function clearSelection() {
    selectedUnit.set(null);
    selectedUnitData.set(null);
}

/**
 * Toggle unit selection (select if not selected, clear if already selected)
 */
export function toggleUnitSelection(unitId, unitData = null, shouldScroll = false) {
    const currentSelected = getSelectedUnitId();
    
    if (currentSelected === unitId) {
        clearSelection();
    } else {
        selectUnit(unitId, unitData, shouldScroll);
    }
}

/**
 * Get currently selected unit ID
 */
function getSelectedUnitId() {
    let currentSelected = null;
    selectedUnit.subscribe(value => currentSelected = value)();
    return currentSelected;
}

/**
 * Load unit data with caching
 */
export async function loadUnitData(unitId, gameId = null) {
    // Check cache first
    const cache = getCurrentCache();
    if (cache.has(unitId)) {
        return cache.get(unitId);
    }

    // Set loading state
    setUnitLoading(unitId, true);

    try {
        let unitData;
        if (gameId) {
            unitData = await getGameUnit(gameId, unitId);
        } else {
            unitData = await getUnitWithHealth(unitId);
        }

        // Cache the result
        updateUnitCache(unitId, unitData);
        setUnitLoading(unitId, false);
        
        return unitData;
    } catch (error) {
        console.error(`Failed to load unit data for ${unitId}:`, error);
        setUnitLoading(unitId, false);
        throw error;
    }
}

/**
 * Preload unit data for a list of unit IDs
 */
export async function preloadUnits(unitIds, gameId = null) {
    const loadPromises = unitIds.map(unitId => 
        loadUnitData(unitId, gameId).catch(error => {
            console.error(`Failed to preload unit ${unitId}:`, error);
            return null;
        })
    );

    const results = await Promise.all(loadPromises);
    return results.filter(result => result !== null);
}

/**
 * Get unit data from cache or load it
 */
export async function getUnitData(unitId, gameId = null) {
    const cache = getCurrentCache();
    if (cache.has(unitId)) {
        return cache.get(unitId);
    }
    return loadUnitData(unitId, gameId);
}

/**
 * Update unit in cache
 */
export function updateUnitCache(unitId, unitData) {
    unitCache.update(cache => {
        const newCache = new Map(cache);
        newCache.set(unitId, unitData);
        return newCache;
    });
}

/**
 * Set loading state for a unit
 */
function setUnitLoading(unitId, loading) {
    unitLoadingStates.update(states => {
        const newStates = new Map(states);
        if (loading) {
            newStates.set(unitId, true);
        } else {
            newStates.delete(unitId);
        }
        return newStates;
    });
}

/**
 * Get current cache (helper function)
 */
function getCurrentCache() {
    let currentCache = new Map();
    unitCache.subscribe(cache => currentCache = cache)();
    return currentCache;
}

/**
 * Check if unit is currently loading
 */
export function isUnitLoading(unitId) {
    let loadingStates = new Map();
    unitLoadingStates.subscribe(states => loadingStates = states)();
    return loadingStates.has(unitId);
}

/**
 * Calculate health percentage from unit data
 */
function calculateHealthPercent(unitData) {
    if (!unitData?.stats) return 100;
    
    const current = unitData.stats.struct || 0;
    const max = unitData.stats.maxStruct || unitData.stats.struct || 1;
    
    return max > 0 ? (current / max) * 100 : 100;
}

/**
 * Clear unit cache
 */
export function clearUnitCache() {
    unitCache.set(new Map());
    unitLoadingStates.set(new Map());
}

/**
 * Cleanup function
 */
export function cleanup() {
    selectedUnit.set(null);
    selectedUnitData.set(null);
    clearUnitCache();
}