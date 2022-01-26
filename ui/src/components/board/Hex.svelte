<script>
	import { Layer } from 'svelte-canvas';
	import { onMount } from 'svelte';

	import event from './store';

	const a = (2 * Math.PI) / 6;
	export let radius = 50;
	export let row;
	export let col;

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
		event.subscribe((value) => {
			// A very simple hitbox, we use topleft as min x,y,
			// and bottomright as max x,y
			if (
				value.x >= topleft.x &&
				value.x <= bottomright.x &&
				value.y >= topleft.y &&
				value.y <= bottomright.y
			) {
				console.log('Clicked row ' + row + ' col ' + col);
			}
		});
	});

	$: render = ({ context, width, height }) => {
		var x = col * (radius * (1 + Math.cos(a)));
		var y = row * (radius * Math.sin(a) + radius * Math.sin(a));
		if (col % 2 == 0) {
			y += 1 ** (row + 1) * radius * Math.sin(a);
		}

		context.lineWidth = 2;
		context.beginPath();

		// This could be done in a neat loop, but I want to keep track of points for simplehitboxing later.
		// This is much easier to read.
		midright = { x: x + radius * Math.cos(a * 0), y: y + radius * Math.sin(a * 0) };
		context.lineTo(midright.x, midright.y);
		bottomright = { x: x + radius * Math.cos(a * 1), y: y + radius * Math.sin(a * 1) };
		context.lineTo(bottomright.x, bottomright.y);
		bottomleft = { x: x + radius * Math.cos(a * 2), y: y + radius * Math.sin(a * 2) };
		context.lineTo(bottomleft.x, bottomleft.y);
		midleft = { x: x + radius * Math.cos(a * 3), y: y + radius * Math.sin(a * 3) };
		context.lineTo(midleft.x, midleft.y);
		topleft = { x: x + radius * Math.cos(a * 4), y: y + radius * Math.sin(a * 4) };
		context.lineTo(topleft.x, topleft.y);
		topright = { x: x + radius * Math.cos(a * 5), y: y + radius * Math.sin(a * 5) };
		context.lineTo(topright.x, topright.y);

		context.closePath();

		var img = new Image(); // Create new img element
		img.src = '/grass_plains_1.gif';
		context.drawImage(img, x - radius, y - radius + 7.5, radius * 2, 0.85 * (radius * 2));

		context.fillText(row + ',' + col, x - 5, y + radius - 8);
		context.stroke();
	};
</script>

<Layer {render} />
