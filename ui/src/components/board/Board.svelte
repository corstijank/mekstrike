<script>
	import { onMount } from 'svelte';
	import {  hexHeight, hexWidth } from './board.js';
	import Unit from './Unit.svelte';
	import event from './store';
	import Terrain from './Terrain.svelte';
	import Highlight from './Highlight.svelte';
	import ActiveUnitHighlight from './ActiveUnitHighlight.svelte';
	import { addGameMessage } from '../../stores/gameStores.js';

	export let id;

	let cells = [];
	let highlights = [];

	let gamedata = { PlayerAUnits: [], PlayerBUnits: [] };
	let availableActions = { UnitOwner: '', AllowedCoordinates: [], CurrentUnitID: '' };
	let activeUnitPosition = null;
	let cols = 0;
	let rows = 0;

	// Slightly overdimension board.
	let boardWidth = 0;
	let boardHeight = 0;

	onMount(() => {
		fetch('/mekstrike/api/gamemaster/games/' + id + '/board')
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				console.log('board');
				console.log(data);
				cells = data.cells;
				cols = parseInt(data.cols);
				rows = parseInt(data.rows);
				// This could be slightly cleaner as we know just oversize the board. We should account for the offset if we want to do it cutely
				boardWidth = (cols + 1) * hexWidth;
				boardHeight = (rows + 1) * hexHeight;
			});
		refreshGameData();
	});

	let previousPhase = -1;
	let previousActiveUnit = '';

	function refreshGameData() {
		fetch('/mekstrike/api/gamemaster/games/' + id)
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				console.log('gamedata');
				console.log(data);
				
				if (previousPhase !== -1 && previousPhase !== data.CurrentPhase) {
					const phases = ['Movement', 'Combat', 'End'];
					addGameMessage('system', `Phase changed to: ${phases[data.CurrentPhase]}`);
				}
				previousPhase = data.CurrentPhase;
				
				gamedata = data;
			});
		fetch('/mekstrike/api/gamemaster/games/' + id + '/availableActions')
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				availableActions = data;
				
				if (previousActiveUnit !== '' && previousActiveUnit !== availableActions.CurrentUnitID) {
					addGameMessage('info', `New active unit: ${availableActions.CurrentUnitID}`);
				}
				previousActiveUnit = availableActions.CurrentUnitID;
				
				// Only show highlights for player units (not CPU)
				if (availableActions.UnitOwner && availableActions.UnitOwner !== 'CPU') {
					highlights = [...availableActions.AllowedCoordinates];
					if (availableActions.UnitOwner !== 'CPU') {
						addGameMessage('movement', `Your turn to move unit ${availableActions.CurrentUnitID}`);
					}
				} else {
					highlights = [];
					if (availableActions.UnitOwner === 'CPU') {
						addGameMessage('info', `CPU is thinking...`);
					}
				}
				
				// Get active unit position for highlighting
				if (availableActions.CurrentUnitID) {
					fetch('/mekstrike/api/gamemaster/games/' + id + '/units/' + availableActions.CurrentUnitID)
						.then((response) => response.json())
						.then((unitData) => {
							activeUnitPosition = {
								x: unitData.location.position.x,
								y: unitData.location.position.y
							};
						});
				}
				
				console.log('highlights');
				console.log(highlights);
			});

		cells = [...cells];
	}
	function handleClick(e) {
		// console.log({ x: e.offsetX, y: e.offsetY });
		event.set({ x: e.offsetX, y: e.offsetY });
	}
</script>

<main>
	<div class="hex-board" style="width: {boardWidth}px; height: {boardHeight}px;">
		<svg width={boardWidth} height={boardHeight}>
			{#each cells as cell}
				<Terrain row={cell.coordinates.y} col={cell.coordinates.x} />
			{/each}
			{#if activeUnitPosition}
				<ActiveUnitHighlight row={activeUnitPosition.y} col={activeUnitPosition.x} />
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
