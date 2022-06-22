<script>
	import { onMount } from 'svelte';

	export let unitID;
	let model = '';
	let name = '';
	let image = '';
	onMount(() => {
		fetch('/mekstrike/api/unit/' + unitID + '/method/GetData')
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				// col = data.location.position.Col + 1;
				// row = data.location.position.Row + 1;
				// heading = data.location.heading;
				name = data.stats.name;
				model = data.stats.model;
				image = data.stats.image;
				// owner = data.owner;
			});
	});
</script>

<div class="terminal-card">
	<header>{model.split(" ")[0]}</header> <!-- Splitting model on space to allow for short cards -->
	<div class="unitCard">
		<img class="unitImage" src={image} alt={model} />
	</div>
</div>

<style>
	.unitCard {
		padding: 0 !important;
	}
	.unitImage {
		height: 10vh;
        display:block;
	}
</style>
