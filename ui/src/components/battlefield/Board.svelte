<script>
	import { onMount } from 'svelte';
	import Unit from '../units/Unit.svelte';
	import Terrain from './Terrain.svelte';
	import Highlight from './Highlight.svelte';
	import ActiveUnitHighlight from './ActiveUnitHighlight.svelte';
	
	// Import stores
	import { gameState, availableActions } from '../../stores/gameStore.js';
	import { boardState, activeUnitPosition, initializeBoard, updateMovementHighlights, updateActiveUnitPosition } from '../../stores/battlefieldStore.js';

	export let id;

	// Subscribe to stores
	$: cells = $boardState.cells;
	$: cols = $boardState.cols;
	$: rows = $boardState.rows;
	$: boardWidth = $boardState.width;
	$: boardHeight = $boardState.height;
	$: gamedata = $gameState;

	// Direct highlight filtering - avoid store timing issues
	$: highlights = $availableActions?.UnitOwner && $availableActions.UnitOwner !== 'CPU' 
		? [...($availableActions.AllowedCoordinates || [])]
		: [];

	// Watch for changes to update battlefield state - only when CurrentUnitID changes
	let previousActiveUnitID = '';
	$: if ($availableActions) {
		// Still update store for other components that might need it
		updateMovementHighlights($availableActions);
		// Only update active unit position when the active unit actually changes
		if ($availableActions.CurrentUnitID !== previousActiveUnitID) {
			updateActiveUnitPosition(id, $availableActions);
			previousActiveUnitID = $availableActions.CurrentUnitID;
		}
	}

	onMount(() => {
		// Initialize board data
		initializeBoard(id);
	});
</script>

<main>
	<div class="hex-board" style="width: {boardWidth}px; height: {boardHeight}px;">
		<svg width={boardWidth} height={boardHeight}>
			{#each cells as cell}
				<Terrain row={cell.coordinates.y} col={cell.coordinates.x} />
			{/each}
			{#if $activeUnitPosition}
				<ActiveUnitHighlight row={$activeUnitPosition.y} col={$activeUnitPosition.x} />
			{/if}
			{#each highlights as highlight}
				<Highlight row={highlight.y} col={highlight.x} />
			{/each}
			{#each gamedata.PlayerAUnits as unitID}
				<Unit game={id} id={unitID} />
			{/each}
			{#each gamedata.PlayerBUnits as unitID}
				<Unit game={id} id={unitID} />
			{/each}
		</svg>
	</div>
</main>

<style>
	.hex-board {
		position: relative;
		width: 100%;
		height: 100%;
		min-width: fit-content;
		min-height: fit-content;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 50px;
	}
	svg {
		display: block;
		margin: 0 auto;
	}
</style>