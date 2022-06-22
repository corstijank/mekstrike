<script>
	import { onMount } from 'svelte';
	import SmallUnitIcon from './SmallUnitIcon.svelte';

	export let id;

	let gamedata = { PlayerAUnits: [], PlayerBUnits: [] };

	onMount(() => {
		fetch('/mekstrike/api/gamemaster/games/' + id)
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				gamedata = data;
			});
	});
</script>

<main>
	<div class="unitlist">
		{#each gamedata.PlayerAUnits as unitID}
			<SmallUnitIcon {unitID} />
			<!-- <Unit id={unitID} /> -->
		{/each}
	</div>
</main>

<style>
	.unitlist {
        display: flex; 
		flex-direction: row;
        justify-content: center;
        column-gap: 10px;
	}
</style>
