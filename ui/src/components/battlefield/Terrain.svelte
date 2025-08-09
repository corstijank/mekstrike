<!-- Hexagon.svelte -->
<script>
	import { colToCenterX, hexSize, rowToCenterY } from '../../utils/coordinates.js';

	export let row = 0;
	export let col = 0;
	export let image = '/grass_plains_1.gif';
	$: x = colToCenterX(col);
	$: y = rowToCenterY(row, col);
	$: points = [
		[hexSize, 0],
		[hexSize / 2, (hexSize * Math.sqrt(3)) / 2],
		[-hexSize / 2, (hexSize * Math.sqrt(3)) / 2],
		[-hexSize, 0],
		[-hexSize / 2, (-hexSize * Math.sqrt(3)) / 2],
		[hexSize / 2, (-hexSize * Math.sqrt(3)) / 2]
	]
		.map(([px, py]) => `${px + x},${py + y}`)
		.join(' ');
</script>

<clipPath id="{col}-{row}">
	<!-- Define the points of the polygon -->
	<polygon {points} />
</clipPath>
<image
	href={image}
	x={x - hexSize}
	y={y - hexSize}
	width={hexSize * 2}
	height={hexSize * 2}
	clip-path="url(#{col}-{row})"
/>
<text x="{x}" y="{y+hexSize-12}" font-size="10" text-anchor="middle">
  {String(col).padStart(3, '0')}|{String(row).padStart(3, '0')}
</text>
