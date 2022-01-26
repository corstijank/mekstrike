<script>
	import { onMount } from 'svelte';
	import { Layer } from 'svelte-canvas';
	import event from './store';

	export let radius = 50;
	export let id;

	const a = (2 * Math.PI) / 6;

	let col = 0;
	let row = 0;
	let heading = 0;

	// Points:
	// 0 is mid right point,
	// 1 is bottom right,
	// 2 bottom left,
	// 3 = mid left,
	// 4 = top left,
	// 5 = top right
	let midright = {},
		bottomright = {},
		bottomleft = {},
		midleft = {},
		topleft = {},
		topright = {};

	onMount(() => {
		fetch('/mekstrike/api/unit/' + id + '/method/GetLocation')
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				col = data.position.Col + 1;
				row = data.position.Row + 1;
				heading = data.heading;
			});
		event.subscribe((value) => {
			// A very simple hitbox, we use topleft as min x,y,
			// and bottomright as max x,y
			if (
				value.x >= topleft.x &&
				value.x <= bottomright.x &&
				value.y >= topleft.y &&
				value.y <= bottomright.y
			) {
				console.log('Clicked unit ' + row + ' col ' + col);
			}
		});
	});

	$: render = ({ context, width, height }) => {
		var x = col * (radius * (1 + Math.cos(a)));
		var y = row * (radius * Math.sin(a) + radius * Math.sin(a));
		if (col % 2 == 0) {
			y += 1 ** (row + 1) * radius * Math.sin(a);
		}
		// This could be done in a neat loop, but I want to keep track of points for simplehitboxing later.
		// This is much easier to read.
		midright = { x: x + radius * Math.cos(a * 0), y: y + radius * Math.sin(a * 0) };
		bottomright = { x: x + radius * Math.cos(a * 1), y: y + radius * Math.sin(a * 1) };
		bottomleft = { x: x + radius * Math.cos(a * 2), y: y + radius * Math.sin(a * 2) };
		midleft = { x: x + radius * Math.cos(a * 3), y: y + radius * Math.sin(a * 3) };
		topleft = { x: x + radius * Math.cos(a * 4), y: y + radius * Math.sin(a * 4) };
		topright = { x: x + radius * Math.cos(a * 5), y: y + radius * Math.sin(a * 5) };

		var img = new Image(); // Create new img element
		img.src = '/unit.png';
		context.drawImage(img, x - radius, y - radius + 7.5, radius * 2, 0.85 * (radius * 2));
	};
</script>

<Layer {render} />
