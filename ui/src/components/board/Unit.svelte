<script>
	import { onMount } from 'svelte';
	import { Layer } from 'svelte-canvas';
	import event from './store';

	export let radius = 50;
	export let game;
	export let id;

	const TO_RADIANS = Math.PI / 180;
	const a = (2 * Math.PI) / 6;

	let col = 0;
	let row = 0;
	let heading = 0;
	let name = '';
	let model = '';
	let owner = '';
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
	let active = false;

	onMount(() => {
		fetch('/mekstrike/api/gamemaster/games/' + game + '/units/' + id)
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				console.log(data);
				col = data.location.position.x;
				row = data.location.position.y;
				heading = data.location.heading;
				name = data.stats.name;
				model = data.stats.model;
				owner = data.owner;
				active = data.active;

				if (active) {
					console.log('Unit ' + name + ' ' + model + ' from ' + owner + ' is active');
				}
			});
		event.subscribe((value) => {
			if (
				value.x >= topleft.x &&
				value.x <= bottomright.x &&
				value.y >= topleft.y &&
				value.y <= bottomright.y
			) {
				handleClick();
			}
		});
	});

	$: render = ({ context, width, height }) => {
		if (name != '') {
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
			img.src = '/mekstrike/media/sprites/' + name;
			rotateAndPaintImage(context, img, TO_RADIANS * (heading * 60), x, y, img.width, img.height);

			context.fillStyle = '#0000FF';

			if (owner == 'CPU') {
				context.fillStyle = '#FF0000';
			}
			if (active) {
				// console.log("Drawing active unit " + model + " from " + owner);
				context.strokeStyle = '#0000FF';
				context.lineWidth = 4;
				context.beginPath();

				// This could be done in a neat loop, but I want to keep track of points for simplehitboxing later.
				// This is much easier to read.
				context.lineTo(midright.x, midright.y);
				context.lineTo(bottomright.x, bottomright.y);
				context.lineTo(bottomleft.x, bottomleft.y);
				context.lineTo(midleft.x, midleft.y);
				context.lineTo(topleft.x, topleft.y);
				context.lineTo(topright.x, topright.y);

				context.closePath();
				context.stroke();
				context.strokeStyle = '#000000';
			}

			context.textAlign = 'center';
			context.fillText(model, x, topleft.y + 9);
			context.fillStyle = '#000000';
		}
	};

	function rotateAndPaintImage(context, image, angleInRad, positionX, positionY, width, height) {
		context.translate(positionX, positionY);
		context.rotate(angleInRad);
		context.drawImage(image, -(width / 2), -(height / 2));
		context.rotate(-angleInRad);
		context.translate(-positionX, -positionY);
	}
	function handleClick() {
		console.log('Clicked unit ' + name + ' ' + model + ' from ' + owner + ' active is ' + active);
	}
</script>

<Layer {render} />
