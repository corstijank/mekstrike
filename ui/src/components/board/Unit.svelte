<script>
	import { onMount } from 'svelte';
	import { colToCenterX, hexSize, rowToCenterY } from "./board";

	export let game;
	export let id;

	let col = 0;
	let row = 0;
	let heading = 0;
	let name = '';
	let model = '';
	let owner = '';
	let active= false;
	const spriteSize=hexSize*1.75
	$: x = colToCenterX(col);
    $: y = rowToCenterY(row,col);

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
			});

	});
</script>

<image transform="rotate({heading * 60}, {x}, {y})"  x="{x-(0.5*spriteSize)}" y="{y-(0.5*spriteSize)}" width="{spriteSize}"  href="/mekstrike/media/sprites/{name}"/>
<text x="{x}" y="{y+20-hexSize}" font-size="10" text-anchor="middle">
	{model}
</text>