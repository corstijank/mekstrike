<script>
	import { addGameMessage } from '../../stores/gameStores.js';

	export let gameId;
	export let availableActions;
	export let selectedUnit;

	let isAdvancing = false;

	const phases = ['Movement', 'Combat', 'End'];

	async function advanceTurn() {
		isAdvancing = true;
		try {
			const response = await fetch(`/mekstrike/api/gamemaster/games/${gameId}/advanceTurn`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				}
			});
			
			if (response.ok) {
				addGameMessage('system', 'Turn advanced');
			} else {
				addGameMessage('error', 'Failed to advance turn');
			}
		} catch (error) {
			addGameMessage('error', `Error advancing turn: ${error.message}`);
		} finally {
			isAdvancing = false;
		}
	}

	function getActionStatus() {
		if (!availableActions.CurrentUnitID) return 'No active unit';
		if (availableActions.UnitOwner === 'CPU') return 'CPU turn';
		return `${availableActions.UnitOwner} - ${phases[availableActions.CurrentPhase] || 'Unknown'} phase`;
	}

	function canAdvanceTurn() {
		return availableActions.UnitOwner && 
			   availableActions.UnitOwner !== 'CPU' && 
			   !isAdvancing;
	}
</script>

<div class="action-controls">
	<div class="terminal-card">
		<header>Actions</header>
		
		<div class="status-section">
			<div class="status-label">Status:</div>
			<div class="status-value">{getActionStatus()}</div>
		</div>

		{#if availableActions.CurrentUnitID}
			<div class="active-unit-section">
				<div class="status-label">Active Unit:</div>
				<div class="status-value">{availableActions.CurrentUnitID}</div>
			</div>
		{/if}

		{#if selectedUnit}
			<div class="selection-section">
				<div class="status-label">Selected:</div>
				<div class="status-value">{selectedUnit}</div>
			</div>
		{/if}

		<div class="button-section">
			<button 
				class="btn btn-primary" 
				disabled={!canAdvanceTurn()}
				on:click={advanceTurn}>
				{#if isAdvancing}
					Advancing...
				{:else}
					End Turn
				{/if}
			</button>

			{#if availableActions.UnitOwner === 'CPU'}
				<div class="cpu-indicator">
					<div class="spinner"></div>
					CPU is thinking...
				</div>
			{/if}
		</div>

		<div class="help-section">
			<div class="help-text">
				• Click units to select/deselect<br>
				• Green circle = selected unit<br>
				• Yellow highlight = available moves<br>
				• Red circle = active unit position
			</div>
		</div>
	</div>
</div>

<style>
	.action-controls {
		height: 100%;
	}

	.status-section, .active-unit-section, .selection-section {
		margin: 10px 0;
		padding: 8px;
		background-color: rgba(0, 0, 0, 0.2);
		border-radius: 4px;
	}

	.status-label {
		font-size: 11px;
		color: #aaa;
		margin-bottom: 4px;
	}

	.status-value {
		font-size: 12px;
		color: #fff;
		font-family: 'Courier New', monospace;
		font-weight: bold;
	}

	.button-section {
		margin: 15px 0;
		text-align: center;
	}

	.btn {
		width: 100%;
		margin-bottom: 10px;
	}

	.cpu-indicator {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		color: #aaa;
		font-size: 11px;
		padding: 8px;
		background-color: rgba(255, 255, 0, 0.1);
		border-radius: 4px;
	}

	.spinner {
		width: 12px;
		height: 12px;
		border: 2px solid #333;
		border-top: 2px solid #ffff00;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	.help-section {
		margin-top: 15px;
		padding: 8px;
		background-color: rgba(0, 0, 0, 0.1);
		border-radius: 4px;
		border: 1px solid #333;
	}

	.help-text {
		font-size: 10px;
		color: #888;
		line-height: 1.3;
	}
</style>