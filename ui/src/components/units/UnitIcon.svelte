<script>
	import { onMount } from 'svelte';
	import { selectedUnit, toggleUnitSelection, loadUnitData } from '../../stores/unitStore.js';
	import { getHealthBarColor, formatUnitDisplayName, calculateHealthPercent } from '../../services/unitService.js';

	export let unitID;
	export let gameId = '';
	export let isActive = false;
	
	let model = '';
	let name = '';
	let struct = 0;
	let maxStruct = 0;
	let unitData = null;
	let loading = true;
	let error = null;
	
	$: isSelected = $selectedUnit === unitID;
	$: healthPercent = calculateHealthPercent(struct, maxStruct);
	$: displayName = formatUnitDisplayName(model);

	onMount(async () => {
		try {
			loading = true;
			unitData = await loadUnitData(unitID);
			
			name = unitData.stats.name;
			model = unitData.stats.model;
			struct = unitData.stats.struct;
			maxStruct = unitData.stats.struct; // Note: API doesn't have maxStruct
			error = null;
		} catch (err) {
			console.error(`Failed to load unit icon data for ${unitID}:`, err);
			error = err.message;
		} finally {
			loading = false;
		}
	});

	function handleClick() {
		if (!loading && !error && unitData) {
			toggleUnitSelection(unitID, unitData, true); // Scroll to unit when clicked from icon
		}
	}
</script>

<button 
	class="unit-button {isActive ? 'active-unit' : ''} {isSelected ? 'selected-unit' : ''} {loading ? 'loading' : ''}" 
	on:click={handleClick}
	disabled={loading || error}
>
	{#if loading}
		<div class="unit-label">Loading...</div>
		<div class="health-bar">
			<div class="health-fill loading-fill"></div>
		</div>
	{:else if error}
		<div class="unit-label error">Error</div>
		<div class="health-bar">
			<div class="health-fill error-fill"></div>
		</div>
	{:else}
		<div class="unit-label">{displayName}</div>
		<div class="health-bar">
			<div class="health-fill" 
				 style="width: {healthPercent}%; background-color: {getHealthBarColor(healthPercent)}">
			</div>
		</div>
	{/if}
	
	{#if isActive}
		<div class="active-badge">●</div>
	{/if}
</button>

<style>
	.unit-button {
		position: relative;
		min-width: 80px;
		height: 60px;
		padding: 6px;
		background-color: #333;
		border: 2px solid #555;
		border-radius: 4px;
		color: #fff;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
	}

	.unit-button:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
		background-color: #444;
	}

	.unit-button:disabled {
		cursor: not-allowed;
		opacity: 0.7;
	}

	.selected-unit {
		border-color: #00ff00;
		box-shadow: 0 0 8px rgba(0, 255, 0, 0.4);
		background-color: #2a4a2a;
	}

	.active-unit {
		border-color: #ffff00;
		box-shadow: 0 0 8px rgba(255, 255, 0, 0.4);
		background-color: #4a4a2a;
		animation: pulse 2s infinite;
	}

	.active-unit.selected-unit {
		border-color: #00ffff;
		box-shadow: 0 0 8px rgba(0, 255, 255, 0.5);
		background-color: #2a4a4a;
	}

	.loading {
		border-color: #666;
		animation: loading-pulse 1.5s infinite;
	}

	@keyframes pulse {
		0%, 100% { box-shadow: 0 0 8px rgba(255, 255, 0, 0.4); }
		50% { box-shadow: 0 0 12px rgba(255, 255, 0, 0.7); }
	}

	@keyframes loading-pulse {
		0%, 100% { opacity: 0.7; }
		50% { opacity: 1; }
	}

	.unit-label {
		font-size: 11px;
		font-weight: bold;
		text-align: center;
		flex-grow: 1;
		display: flex;
		align-items: center;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.unit-label.error {
		color: #ff4444;
	}

	.health-bar {
		width: 100%;
		height: 4px;
		background-color: #222;
		border-radius: 2px;
		overflow: hidden;
		margin-top: 4px;
	}

	.health-fill {
		height: 100%;
		transition: width 0.3s ease, background-color 0.3s ease;
		border-radius: 2px;
	}

	.loading-fill {
		width: 100%;
		background: linear-gradient(90deg, #333, #666, #333);
		animation: loading-shimmer 1.5s infinite;
	}

	.error-fill {
		width: 100%;
		background-color: #ff4444;
	}

	@keyframes loading-shimmer {
		0% { background-position: -200px 0; }
		100% { background-position: 200px 0; }
	}

	.active-badge {
		position: absolute;
		top: 2px;
		right: 2px;
		color: #ffff00;
		font-size: 10px;
		font-weight: bold;
	}
</style>