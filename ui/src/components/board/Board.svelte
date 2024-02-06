<script>
	import { Canvas, Layer } from 'svelte-canvas';
	import { onMount } from 'svelte';
	import Hex from './Hex.svelte';
	import Unit from './Unit.svelte';
	import event from './store';
	import Highlight from './Highlight.svelte';

	export let id;

	let canvas;

	let cells = [
	];
	let highlights = [];

	let gamedata = { PlayerAUnits: [], PlayerBUnits: [] };
	let currentOpts = { AllowedCells: [] };
	let cols = 0;
	let rows = 0;

	onMount(() => {
		fetch('/mekstrike/api/battlefield/' + id + '/method/GetBoardCells')
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				cells = data;
			});

		fetch('/mekstrike/api/battlefield/' + id + '/method/GetNumberOfCols')
			.then((response) => {
				return response.text();
			})
			.then((data) => {
				cols = parseInt(data);
			});
		fetch('/mekstrike/api/battlefield/' + id + '/method/GetNumberOfRows')
			.then((response) => {
				return response.text();
			})
			.then((data) => {
				rows = parseInt(data);
			});

		refreshGameData();
	});

	function refreshGameData() {
		fetch('/mekstrike/api/gamemaster/games/' + id)
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				gamedata = data;
			});
		fetch('/mekstrike/api/gamemaster/games/' + id + '/currentOpts')
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				currentOpts = data;
				highlights = [...currentOpts.AllowedCells];
				console.log("highlights");
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
			<Hex row={cell.coordinates.Row + 1} col={cell.coordinates.Col + 1} />
		{/each}
		{#each highlights as highlight}
			<Highlight row={highlight.Row + 1} col={highlight.Col + 1} />
		{/each}
		{#each gamedata.PlayerAUnits as unitID}
			<Unit id={unitID} />
		{/each}
		{#each gamedata.PlayerBUnits as unitID}
			<Unit id={unitID} />
		{/each}
	</Canvas>
</main>
