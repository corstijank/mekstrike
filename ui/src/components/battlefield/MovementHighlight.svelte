<script>
	import { colToCenterX, hexSize, rowToCenterY } from '../../utils/coordinates.js';
	import { createEventDispatcher } from 'svelte';

	export let row = 0;
	export let col = 0;
	export let isSelected = false;
	export let showHeadingSelector = false;
	
	const dispatch = createEventDispatcher();
	
	$: x = colToCenterX(col);
	$: y = rowToCenterY(row, col);
	$: hexPoints = [
		[hexSize, 0],
		[hexSize / 2, (hexSize * Math.sqrt(3)) / 2],
		[-hexSize / 2, (hexSize * Math.sqrt(3)) / 2],
		[-hexSize, 0],
		[-hexSize / 2, (-hexSize * Math.sqrt(3)) / 2],
		[hexSize / 2, (-hexSize * Math.sqrt(3)) / 2]
	]
		.map(([px, py]) => `${px + x},${py + y}`)
		.join(' ');

	// Heading directions (0-5 for hex directions)
	const headingDirections = [
		{ heading: 0, angle: 0, label: 'N' },      // North
		{ heading: 1, angle: 60, label: 'NE' },   // Northeast  
		{ heading: 2, angle: 120, label: 'SE' },  // Southeast
		{ heading: 3, angle: 180, label: 'S' },   // South
		{ heading: 4, angle: 240, label: 'SW' },  // Southwest
		{ heading: 5, angle: 300, label: 'NW' }   // Northwest
	];

	function handleClick() {
		if (!showHeadingSelector) {
			dispatch('click', { row, col });
		}
	}

	function handleHeadingSelect(heading) {
		dispatch('headingSelected', { row, col, heading });
	}

	function calculateArrowPosition(angle) {
		const radius = hexSize * 0.7; // Position arrows near edge of hex
		const radians = (angle - 90) * (Math.PI / 180); // -90 to make 0° point up
		return {
			x: x + radius * Math.cos(radians),
			y: y + radius * Math.sin(radians)
		};
	}
</script>

<!-- Main movement highlight hex -->
<polygon 
	points={hexPoints}
	fill={isSelected ? "rgba(255, 215, 0, 0.4)" : "rgba(0, 150, 255, 0.3)"} 
	stroke={isSelected ? "rgba(255, 215, 0, 0.8)" : "rgba(0, 150, 255, 0.6)"} 
	stroke-width="2"
	class="movement-highlight"
	on:click={handleClick}
/>

<!-- Heading selector arrows (shown when selected) -->
{#if showHeadingSelector}
	{#each headingDirections as direction}
		{@const arrowPos = calculateArrowPosition(direction.angle)}
		<g class="heading-selector">
			<!-- Arrow background circle -->
			<circle 
				cx={arrowPos.x} 
				cy={arrowPos.y} 
				r="12" 
				fill="rgba(255, 255, 255, 0.9)"
				stroke="rgba(0, 100, 200, 0.8)"
				stroke-width="2"
				class="heading-option"
				on:click={() => handleHeadingSelect(direction.heading)}
			/>
			<!-- Arrow icon -->
			<text 
				x={arrowPos.x} 
				y={arrowPos.y + 2} 
				text-anchor="middle" 
				dominant-baseline="middle"
				font-family="monospace"
				font-size="10"
				font-weight="bold"
				fill="rgba(0, 100, 200, 0.9)"
				class="heading-text"
				transform="rotate({direction.angle} {arrowPos.x} {arrowPos.y})"
				on:click={() => handleHeadingSelect(direction.heading)}
			>
				▲
			</text>
			<!-- Heading label -->
			<text 
				x={arrowPos.x} 
				y={arrowPos.y + 20} 
				text-anchor="middle" 
				font-family="sans-serif"
				font-size="8"
				fill="rgba(0, 0, 0, 0.7)"
				class="heading-label"
			>
				{direction.label}
			</text>
		</g>
	{/each}
	
	<!-- Center cancel button -->
	<circle 
		cx={x} 
		cy={y} 
		r="8" 
		fill="rgba(255, 100, 100, 0.9)"
		stroke="rgba(200, 0, 0, 0.8)"
		stroke-width="1"
		class="cancel-button"
		on:click={() => dispatch('cancel')}
	/>
	<text 
		x={x} 
		y={y + 2} 
		text-anchor="middle" 
		dominant-baseline="middle"
		font-family="sans-serif"
		font-size="10"
		fill="white"
		class="cancel-text"
		on:click={() => dispatch('cancel')}
	>
		×
	</text>
{/if}

<style>
	.movement-highlight {
		cursor: pointer;
		transition: all 0.2s ease;
	}
	
	.movement-highlight:hover {
		filter: brightness(1.2);
		stroke-width: 3;
	}
	
	.heading-option {
		cursor: pointer;
		transition: all 0.15s ease;
	}
	
	.heading-option:hover {
		fill: rgba(255, 255, 255, 1);
		stroke-width: 3;
		r: 14;
	}
	
	.heading-text {
		cursor: pointer;
		pointer-events: none;
	}
	
	.heading-label {
		pointer-events: none;
		font-weight: bold;
	}
	
	.cancel-button {
		cursor: pointer;
		transition: all 0.15s ease;
	}
	
	.cancel-button:hover {
		fill: rgba(255, 150, 150, 1);
		r: 10;
	}
	
	.cancel-text {
		cursor: pointer;
		pointer-events: none;
		font-weight: bold;
	}
</style>