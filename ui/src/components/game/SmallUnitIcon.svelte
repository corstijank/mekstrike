<script>
	import { onMount } from 'svelte';
	import { selectedUnit, selectUnitAndScroll, clearSelection } from '../../stores/gameStores.js';

	export let unitID;
	export let gameId = '';
	export let isActive = false;
	
	let model = '';
	let name = '';
	let struct = 0;
	let maxStruct = 0;
	let unitData = null;
	
	$: isSelected = $selectedUnit === unitID;
	$: healthPercent = maxStruct > 0 ? (struct / maxStruct) * 100 : 100;

	onMount(() => {
		fetch('/mekstrike/api/unit/' + unitID + '/method/GetData')
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				name = data.stats.name;
				model = data.stats.model;
				struct = data.stats.struct;
				maxStruct = data.stats.struct;
				unitData = data;
			});
	});

	function handleClick() {
		if ($selectedUnit === unitID) {
			clearSelection();
		} else if (unitData) {
			selectUnitAndScroll(unitID, unitData);
		}
	}

	function getHealthBarColor(percent) {
		if (percent > 66) return '#00ff00';
		if (percent > 33) return '#ffff00';
		return '#ff0000';
	}
</script>

<button class="unit-button {isActive ? 'active-unit' : ''} {isSelected ? 'selected-unit' : ''}" 
        on:click={handleClick}>
	<div class="unit-label">{model.split(" ")[0]}</div>
	<div class="health-bar">
		<div class="health-fill" 
			 style="width: {healthPercent}%; background-color: {getHealthBarColor(healthPercent)}">
		</div>
	</div>
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

	.unit-button:hover {
		transform: translateY(-1px);
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
		background-color: #444;
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

	@keyframes pulse {
		0%, 100% { box-shadow: 0 0 8px rgba(255, 255, 0, 0.4); }
		50% { box-shadow: 0 0 12px rgba(255, 255, 0, 0.7); }
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

	.active-badge {
		position: absolute;
		top: 2px;
		right: 2px;
		color: #ffff00;
		font-size: 10px;
		font-weight: bold;
	}
</style>
