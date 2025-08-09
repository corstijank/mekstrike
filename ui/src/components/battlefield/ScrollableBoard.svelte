<script>
	import Board from './Board.svelte';
	import { clearSelection } from '../../stores/unitStore.js';
	import { scrollToPosition, clearScrollPosition } from '../../stores/battlefieldStore.js';
	import { colToCenterX, rowToCenterY } from '../../utils/coordinates.js';
	import { getScrollPosition } from '../../services/battlefieldService.js';

	export let id;

	let boardViewport;

	function handleBoardClick(event) {
		if (event.target.closest('g[style*="cursor: pointer"]')) {
			return;
		}
		clearSelection();
	}

	// Watch for scroll position changes
	$: if ($scrollToPosition && boardViewport) {
		centerOnUnit($scrollToPosition.x, $scrollToPosition.y);
		clearScrollPosition();
	}

	function centerOnUnit(col, row) {
		if (!boardViewport) return;
		
		const viewportWidth = boardViewport.clientWidth;
		const viewportHeight = boardViewport.clientHeight;
		const { scrollX, scrollY } = getScrollPosition(col, row, viewportWidth, viewportHeight);
		
		boardViewport.scrollTo({
			left: scrollX,
			top: scrollY,
			behavior: 'smooth'
		});
	}
</script>

<div class="scrollable-board-container" on:click={handleBoardClick}>
	<div class="board-viewport" bind:this={boardViewport}>
		<Board {id} />
	</div>
</div>

<style>
	.scrollable-board-container {
		width: 100%;
		height: 100%;
		position: relative;
		overflow: hidden;
		background-color: #1a1a1a;
		border: 2px solid #333;
		border-radius: 4px;
	}

	.board-viewport {
		width: 100%;
		height: 100%;
		overflow: auto;
		position: relative;
		box-sizing: border-box;
	}

	.board-viewport::-webkit-scrollbar {
		width: 12px;
		height: 12px;
	}

	.board-viewport::-webkit-scrollbar-track {
		background: #2a2a2a;
		border-radius: 6px;
	}

	.board-viewport::-webkit-scrollbar-thumb {
		background: #555;
		border-radius: 6px;
		border: 2px solid #2a2a2a;
	}

	.board-viewport::-webkit-scrollbar-thumb:hover {
		background: #777;
	}

	.board-viewport::-webkit-scrollbar-corner {
		background: #2a2a2a;
	}
</style>