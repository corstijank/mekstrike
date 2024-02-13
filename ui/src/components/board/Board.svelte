<script>
	import { Canvas, Layer } from 'svelte-canvas';
	import { onMount } from 'svelte';
	import Hex from './Hex.svelte';
	import Unit from './Unit.svelte';
	import event from './store';
	import Highlight from './Highlight.svelte';

	export let id;

	let canvas;

	let cells = [];
	let highlights = [];

	let gamedata = { PlayerAUnits: [], PlayerBUnits: [] };
	let currentOpts = { AllowedCells: [] };
	let cols = 0;
	let rows = 0;

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
	<!-- TODO: This width/height is probably magic and won't scale for bigger boards. Test and fix-->
	<Canvas
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
	</Canvas>
</main>
