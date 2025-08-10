<script>
	import { onMount } from 'svelte';
	import { colToCenterX, hexSize, rowToCenterY } from '../../utils/coordinates.js';
	import { selectedUnit, selectUnit, clearSelection, loadUnitData } from '../../stores/unitStore.js';
	import { gameState } from '../../stores/gameStore.js';
	import { boardVersion } from '../../stores/battlefieldStore.js';
	import { gameAPI } from '../../services/gameAPI.js';

	export let game;
	export let id;

	let col = 0;
	let row = 0;
	let heading = 0;
	let name = '';
	let model = '';
	let owner = '';
	let active = false;
	let unitData = null;
	let loading = true;
	let error = null;

	const spriteSize = hexSize * 1.75;
	
	$: x = colToCenterX(col);
	$: y = rowToCenterY(row, col);
	$: isSelected = $selectedUnit === id;

	// Reactive function to load unit data
	async function loadUnit() {
		try {
			loading = true;
			const data = await gameAPI.getUnit(game, id);
			
			const oldCol = col;
			const oldRow = row;
			
			col = data.location.position.x;
			row = data.location.position.y;
			heading = data.location.heading;
			name = data.stats.name;
			model = data.stats.model;
			owner = data.owner;
			active = data.active;
			unitData = data;
			error = null;
			
		} catch (err) {
			console.error(`Failed to load unit ${id}:`, err);
			error = err.message;
		} finally {
			loading = false;
		}
	}

	// Load unit data on mount
	onMount(() => {
		loadUnit();
	});

	// Track last known state to detect changes
	let lastRound = -1;
	let lastPhase = -1;
	let lastBoardVersion = -1;
	
	// Reactively reload unit data when round or phase changes (indicates unit movement/turn advance)
	$: if ($gameState && ($gameState.CurrentRound !== lastRound || $gameState.CurrentPhase !== lastPhase)) {
		// Only reload if we've seen an initial load and state actually changed
		if (lastRound >= 0 && !loading) {
			loadUnit();
		}
		lastRound = $gameState.CurrentRound;
		lastPhase = $gameState.CurrentPhase;
	}
	
	// Reactively reload unit data when board is refreshed (indicates WebSocket unit position updates)
	$: if ($boardVersion && $boardVersion !== lastBoardVersion) {
		// Only reload if we've seen an initial load and version actually changed
		if (lastBoardVersion >= 0 && !loading) {
			loadUnit();
		}
		lastBoardVersion = $boardVersion;
	}

	function handleUnitClick(event) {
		event.stopPropagation();
		
		if ($selectedUnit === id) {
			clearSelection();
		} else if (unitData) {
			selectUnit(id, unitData, false); // Don't scroll from board unit click
		}
	}
</script>

{#if !loading && !error}
	<g on:click={handleUnitClick} style="cursor: pointer;">
		{#if isSelected}
			<circle cx={x} cy={y} r={spriteSize * 0.6} fill="none" stroke="#00ff00" stroke-width="3" opacity="0.8"/>
		{/if}
		<image 
			transform="rotate({heading * 60}, {x}, {y})"  
			x="{x - (0.5 * spriteSize)}" 
			y="{y - (0.5 * spriteSize)}" 
			width="{spriteSize}"  
			href="/mekstrike/media/sprites/{name}"
			alt="{model}"
		/>
		<text x="{x}" y="{y + 20 - hexSize}" font-size="10" text-anchor="middle">
			{model}
		</text>
	</g>
{:else if error}
	<g>
		<circle cx={x} cy={y} r={spriteSize * 0.3} fill="#ff4444" opacity="0.7"/>
		<text x="{x}" y="{y + 4}" font-size="8" text-anchor="middle" fill="white">
			ERROR
		</text>
	</g>
{/if}