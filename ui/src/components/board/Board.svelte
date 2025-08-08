<script>
	import { onMount } from 'svelte';
	import {  hexHeight, hexWidth } from './board.js';
	import Unit from './Unit.svelte';
	import event from './store';
	import Terrain from './Terrain.svelte';

	export let id;

	let cells = [];
	let highlights = [];

	let gamedata = { PlayerAUnits: [], PlayerBUnits: [] };
	let currentOpts = { AllowedCells: [] };
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

	function refreshGameData() {
		fetch('/mekstrike/api/gamemaster/games/' + id)
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				console.log('gamedata');
				console.log(data);
				gamedata = data;
			});
		fetch('/mekstrike/api/gamemaster/games/' + id + '/currentOpts')
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				currentOpts = data;
				highlights = [...currentOpts.AllowedCoordinates];
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
			{#each gamedata.PlayerAUnits as unitID}
				<Unit game={id} id={unitID} />
			{/each}
			{#each gamedata.PlayerBUnits as unitID}
				<Unit game={id} id={unitID} />
			{/each}
		</svg>
	</div>
	<!-- TODO: This width/height is probably magic and won't scale for bigger boards. Test and fix-->
	<!-- <Canvas
		width={cols * 80}
		height={rows * 95}
		on:click={(e) => handleClick(e)}
		style="display:inline"
	>
		{#each cells as cell}
			<Hex row={cell.coordinates.y} col={cell.coordinates.x} />
		{/each}
		{#each highlights as highlight}
			<Highlight row={highlight.y} col={highlight.x} />
		{/each}
		{#each gamedata.PlayerAUnits as unitID}
			<Unit game={id} id={unitID} />
		{/each}
		{#each gamedata.PlayerBUnits as unitID}
			<Unit game={id} id={unitID} />
		{/each}
	</Canvas> -->
</main>

<style>
	.hex-board {
		position: relative;
		overflow: hidden;
	}
	svg {
		position: absolute;
		top: 10px;
		margin-left: auto;
		margin-right: auto;
		transform: translateX(-50%);
	}
</style>
