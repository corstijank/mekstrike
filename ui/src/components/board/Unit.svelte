<script>
	import { onMount } from 'svelte';
	import { colToCenterX, hexSize, rowToCenterY } from "./board";
	import { selectedUnit, selectUnit, clearSelection } from '../../stores/gameStores.js';

	export let game;
	export let id;

	let col = 0;
	let row = 0;
	let heading = 0;
	let name = '';
	let model = '';
	let owner = '';
	let active= false;
	let unitData = null;
	const spriteSize=hexSize*1.75
	$: x = colToCenterX(col);
    $: y = rowToCenterY(row,col);
	$: isSelected = $selectedUnit === id;

	onMount(() => {
		fetch('/mekstrike/api/gamemaster/games/' + game + '/units/' + id)
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				col = data.location.position.x;
				row = data.location.position.y;
				heading = data.location.heading;
				name = data.stats.name;
				model = data.stats.model;
				owner = data.owner;
				active = data.active;
				unitData = data;
			});
	});

	function handleUnitClick(event) {
		event.stopPropagation();
		
		if ($selectedUnit === id) {
			clearSelection();
		} else if (unitData) {
			selectUnit(id, unitData);
		}
	}
</script>

<g on:click={handleUnitClick} style="cursor: pointer;">
	{#if isSelected}
		<circle cx={x} cy={y} r={spriteSize*0.6} fill="none" stroke="#00ff00" stroke-width="3" opacity="0.8"/>
	{/if}
	<image transform="rotate({heading * 60}, {x}, {y})"  x="{x-(0.5*spriteSize)}" y="{y-(0.5*spriteSize)}" width="{spriteSize}"  href="/mekstrike/media/sprites/{name}"/>
	<text x="{x}" y="{y+20-hexSize}" font-size="10" text-anchor="middle">
		{model}
	</text>
</g>