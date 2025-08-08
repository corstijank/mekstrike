<script>
	import Board from './Board.svelte';
	import { clearSelection, scrollToUnit } from '../../stores/gameStores.js';
	import { colToCenterX, rowToCenterY } from './board.js';

	export let id;

	let boardViewport;

	function handleBoardClick(event) {
		if (event.target.closest('g[style*="cursor: pointer"]')) {
			return;
		}
		clearSelection();
	}

	$: if ($scrollToUnit && boardViewport) {
		centerOnUnit($scrollToUnit.x, $scrollToUnit.y);
		scrollToUnit.set(null);
	}

	function centerOnUnit(col, row) {
		if (!boardViewport) return;
		
		const x = colToCenterX(col);
		const y = rowToCenterY(row, col);
		
		const viewportWidth = boardViewport.clientWidth;
		const viewportHeight = boardViewport.clientHeight;
		
		const scrollX = x - (viewportWidth / 2) + 50;
		const scrollY = y - (viewportHeight / 2) + 50;
		
		boardViewport.scrollTo({
			left: Math.max(0, scrollX),
			top: Math.max(0, scrollY),
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