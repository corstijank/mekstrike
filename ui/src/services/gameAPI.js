/**
 * Centralized API service for game-related requests
 * Handles all communication with gamemaster service
 */

const API_BASE = '/mekstrike/api';
const GAMEMASTER_BASE = `${API_BASE}/gamemaster`;
const UNIT_BASE = `${API_BASE}/unit`;

/**
 * Generic fetch wrapper with error handling and retry logic
 */
async function apiRequest(url, options = {}) {
    const config = {
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
        ...options,
    };

    try {
        const response = await fetch(url, config);
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error(`API request failed for ${url}:`, error);
        throw error;
    }
}

/**
 * Game management endpoints
 */
export const gameAPI = {
    /**
     * Get game state including current round, phase, active player
     */
    async getGame(gameId) {
        return apiRequest(`${GAMEMASTER_BASE}/games/${gameId}`);
    },

    /**
     * Get available actions for current active unit
     */
    async getAvailableActions(gameId) {
        return apiRequest(`${GAMEMASTER_BASE}/games/${gameId}/availableActions`);
    },

    /**
     * Get current movement options (replaces availableActions for movement)
     */
    async getCurrentOptions(gameId) {
        return apiRequest(`${GAMEMASTER_BASE}/games/${gameId}/currentOpts`);
    },

    /**
     * Get battlefield board state
     */
    async getBoard(gameId) {
        return apiRequest(`${GAMEMASTER_BASE}/games/${gameId}/board`);
    },

    /**
     * Get unit data by unit ID
     */
    async getUnit(gameId, unitId) {
        return apiRequest(`${GAMEMASTER_BASE}/games/${gameId}/units/${unitId}`);
    },

    /**
     * Get game logs/events
     */
    async getGameLogs(gameId) {
        return apiRequest(`${GAMEMASTER_BASE}/games/${gameId}/logs`);
    },

    /**
     * Move unit to specified coordinates with heading
     */
    async moveUnit(gameId, unitId, x, y, heading) {
        return apiRequest(`${GAMEMASTER_BASE}/games/${gameId}/move`, {
            method: 'POST',
            body: JSON.stringify({
                unitId,
                x,
                y,
                heading
            }),
        });
    },

    /**
     * Advance turn (end current unit's turn)
     */
    async advanceTurn(gameId) {
        return apiRequest(`${GAMEMASTER_BASE}/games/${gameId}/advanceTurn`, {
            method: 'POST',
        });
    },

    /**
     * Create a new game
     */
    async createGame(gameData) {
        return apiRequest(`${GAMEMASTER_BASE}/games`, {
            method: 'POST',
            body: JSON.stringify(gameData),
        });
    },

    /**
     * Get list of games
     */
    async getGames() {
        return apiRequest(`${GAMEMASTER_BASE}/games`);
    },
};

/**
 * Unit-specific endpoints (Dapr actor calls)
 */
export const unitAPI = {
    /**
     * Get unit data via Dapr actor
     */
    async getUnitData(unitId) {
        return apiRequest(`${UNIT_BASE}/${unitId}/method/GetData`);
    },

    /**
     * Deploy unit to battlefield
     */
    async deployUnit(unitId, position) {
        return apiRequest(`${UNIT_BASE}/${unitId}/method/deploy`, {
            method: 'POST',
            body: JSON.stringify(position),
        });
    },

    /**
     * Set unit active/inactive
     */
    async setUnitActive(unitId, active) {
        return apiRequest(`${UNIT_BASE}/${unitId}/method/setActive`, {
            method: 'POST',
            body: JSON.stringify({ active }),
        });
    },
};

/**
 * Battlefield-specific endpoints (Dapr actor calls)
 */
export const battlefieldAPI = {
    /**
     * Get board cells from battlefield actor
     */
    async getBoardCells(battlefieldId) {
        return apiRequest(`${API_BASE}/battlefield/${battlefieldId}/method/getBoardCells`);
    },

    /**
     * Get movement options for a unit
     */
    async getMovementOptions(battlefieldId, unitPosition) {
        return apiRequest(`${API_BASE}/battlefield/${battlefieldId}/method/getMovementOptions`, {
            method: 'POST',
            body: JSON.stringify(unitPosition),
        });
    },

    /**
     * Check if a cell is blocked
     */
    async isCellBlocked(battlefieldId, position) {
        return apiRequest(`${API_BASE}/battlefield/${battlefieldId}/method/isCellBlocked`, {
            method: 'POST',
            body: JSON.stringify(position),
        });
    },
};