<script>
	import Board from '../../components/board/Board.svelte';
	import PlayerUnitList from '../../components/game/PlayerUnitList.svelte';
	import { page } from '$app/stores';
	import 'terminal.css';
</script>

<div class="map"><Board id={$page.params.id} /></div>
<div class="bottombar">
	<PlayerUnitList id={$page.params.id} />
</div>

<style>
	.map {
		width: 100%;
		text-align: center;
	}
	.bottombar {
		position: fixed; /* Set the navbar to fixed position */
		bottom: 2vh; /* Position the navbar at the bottom of the page */
		width: 100%; /* Full width */
		z-index: 999;
	}
</style>
