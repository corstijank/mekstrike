<script>
	import { onMount } from 'svelte';
	import Unit from '../units/Unit.svelte';
	import Terrain from './Terrain.svelte';
	import Highlight from './Highlight.svelte';
	import MovementHighlight from './MovementHighlight.svelte';
	import ActiveUnitHighlight from './ActiveUnitHighlight.svelte';
	
	// Import stores and API
	import { gameState, availableActions } from '../../stores/gameStore.js';
	import { boardState, activeUnitPosition, initializeBoard, updateMovementHighlights, updateActiveUnitPosition } from '../../stores/battlefieldStore.js';
	import { gameAPI } from '../../services/gameAPI.js';

	export let id;

	// Subscribe to stores
	$: cells = $boardState.cells;
	$: cols = $boardState.cols;
	$: rows = $boardState.rows;
	$: boardWidth = $boardState.width;
	$: boardHeight = $boardState.height;
	$: gamedata = $gameState;

	// Movement selection state
	let selectedMovement = null; // {row, col} of selected movement hex
	let isMoving = false; // Flag to prevent multiple moves

	// Direct highlight filtering - avoid store timing issues
	$: highlights = $availableActions?.UnitOwner && $availableActions.UnitOwner !== 'CPU' 
		? [...($availableActions.AllowedCoordinates || [])]
		: [];
	
	// Clear selection when available actions change (new turn, etc.)
	$: if ($availableActions) {
		selectedMovement = null;
	}

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

	// Movement selection handlers
	function handleMovementClick(event) {
		if (isMoving) return;
		const { row, col } = event.detail;
		selectedMovement = { row, col };
	}

	function handleHeadingSelected(event) {
		if (isMoving) return;
		const { row, col, heading } = event.detail;
		executeMove(row, col, heading);
	}

	function handleMovementCancel() {
		selectedMovement = null;
	}

	async function executeMove(row, col, heading) {
		if (!$availableActions?.CurrentUnitID || isMoving) return;
		
		isMoving = true;
		try {
			console.log(`Moving unit ${$availableActions.CurrentUnitID} to (${col}, ${row}) with heading ${heading}`);
			
			await gameAPI.moveUnit(
				id, 
				$availableActions.CurrentUnitID, 
				col, 
				row, 
				heading
			);
			
			// Clear selection after successful move
			selectedMovement = null;
			
			// The game state will update via the normal polling mechanism
			console.log('Move completed successfully');
			
		} catch (error) {
			console.error('Failed to execute move:', error);
			// You could show an error message to the user here
		} finally {
			isMoving = false;
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
				<MovementHighlight 
					row={highlight.y} 
					col={highlight.x}
					isSelected={selectedMovement?.row === highlight.y && selectedMovement?.col === highlight.x}
					showHeadingSelector={selectedMovement?.row === highlight.y && selectedMovement?.col === highlight.x}
					on:click={handleMovementClick}
					on:headingSelected={handleHeadingSelected}
					on:cancel={handleMovementCancel}
				/>
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