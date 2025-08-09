<script>
	import { colToCenterX, hexSize, rowToCenterY } from '../../utils/coordinates.js';

	export let row = 0;
	export let col = 0;
	
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

<polygon 
	{points} 
	fill="rgba(255, 165, 0, 0.4)" 
	stroke="rgba(255, 165, 0, 0.8)" 
	stroke-width="3"
/>