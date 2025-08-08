import { writable } from 'svelte/store';

export const selectedUnit = writable(null);

export const gameMessages = writable([]);

export const selectedUnitData = writable(null);

export const scrollToUnit = writable(null);

export function addGameMessage(type, message, timestamp = new Date()) {
    gameMessages.update(messages => [...messages, { type, message, timestamp }]);
}

export function selectUnit(unitId, unitData = null) {
    selectedUnit.set(unitId);
    selectedUnitData.set(unitData);
}

export function selectUnitAndScroll(unitId, unitData = null) {
    selectedUnit.set(unitId);
    selectedUnitData.set(unitData);
    if (unitData && unitData.location) {
        scrollToUnit.set({
            x: unitData.location.position.x,
            y: unitData.location.position.y
        });
    }
}

export function clearSelection() {
    selectedUnit.set(null);
    selectedUnitData.set(null);
}