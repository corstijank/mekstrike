<script>
	import { Canvas, Layer } from 'svelte-canvas';
	import { onMount } from 'svelte';
	import Hex from './Hex.svelte';
	import Unit from './Unit.svelte';
	import event from './store';

	export let id;

	let canvas;

	let cells = [
		{ Col: 1, Row: 1, TerrainTypeID: 1 },
		{ Col: 1, Row: 2, TerrainTypeID: 1 },
		{ Col: 2, Row: 1, TerrainTypeID: 1 },
		{ Col: 2, Row: 2, TerrainTypeID: 1 }
	];

	let gamedata = { PlayerAUnits: [], PlayerBUnits: [] };
	let cols = 0;
	let rows = 0;

	onMount(() => {
		fetch('/mekstrike/api/gamemaster/games/' + id)
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				gamedata = data;
			});
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
	});
	function handleClick(e) {
		// console.log({ x: e.offsetX, y: e.offsetY });
		event.set({ x: e.offsetX, y: e.offsetY });
	}
</script>

<main>
	<!-- TODO: This width/height is probably magic and won't scale for bigger boards. Test and fix-->
	<Canvas width={cols * 80} height={rows * 95} on:click={(e) => handleClick(e)}>
		{#each cells as cell}
			<Hex row={cell.Row + 1} col={cell.Col + 1} />
		{/each}
		{#each gamedata.PlayerAUnits as unitID}
			<Unit id={unitID} />
		{/each}
		{#each gamedata.PlayerBUnits as unitID}
			<Unit id={unitID} />
		{/each}
	</Canvas>
</main>
