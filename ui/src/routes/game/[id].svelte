<script>
	import ScrollableBoard from '../../components/board/ScrollableBoard.svelte';
	import PlayerUnitList from '../../components/game/PlayerUnitList.svelte';
	import GameMessageArea from '../../components/game/GameMessageArea.svelte';
	import UnitStatsPanel from '../../components/game/UnitStatsPanel.svelte';
	import { page } from '$app/stores';
	import 'terminal.css';
</script>

<div class="game-layout">
	<div class="left-panel">
		<GameMessageArea gameId={$page.params.id} />
	</div>
	
	<div class="center-panel">
		<ScrollableBoard id={$page.params.id} />
	</div>
	
	<div class="right-panel">
		<UnitStatsPanel />
	</div>
	
	<div class="bottom-panel">
		<PlayerUnitList id={$page.params.id} />
	</div>
</div>

<style>
	.game-layout {
		display: grid;
		grid-template-columns: 20vw 1fr 20vw;
		grid-template-rows: 1fr 155px;
		grid-template-areas: 
			"left center right"
			"bottom bottom bottom";
		height: 100vh;
		width: 100vw;
		position: fixed;
		top: 0;
		left: 0;
		background-color: #111;
		gap: 0;
	}

	.left-panel {
		grid-area: left;
		border-right: 2px solid #333;
		background-color: #1a1a1a;
		overflow: hidden;
	}

	.center-panel {
		grid-area: center;
		background-color: #111;
		position: relative;
		overflow: hidden;
	}

	.right-panel {
		grid-area: right;
		border-left: 2px solid #333;
		background-color: #1a1a1a;
		overflow: hidden;
	}

	.bottom-panel {
		grid-area: bottom;
		border-top: 2px solid #333;
		background-color: #222;
		height: 155px;
		overflow: hidden;
	}
</style>
