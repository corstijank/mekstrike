<script>
	import { createEventDispatcher } from 'svelte';
	let playername;
	let dispatch = createEventDispatcher();
	function createGame() {
		console.log('Creating new game for ' + playername);
		fetch('/mekstrike/api/gamemaster/games', {
			method: 'POST', // or 'PUT'
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ PlayerName: playername })
		}).then((response) => {
			console.log('Game created');
			dispatch('NewGame');
		});
	}
</script>

<main>
	<div class="terminal-card">
		<header>Create New game</header>
		<div class="background-transparent">
			<input
				id="playername"
				name="playername"
				type="text"
				minlength="4"
				placeholder="playername"
				bind:value={playername}
			/>
			<button class="btn btn-primary" on:click={createGame}>Create</button>
		</div>
	</div>
</main>

<style>
	.background-transparent {
		background: rgba(0, 0, 0, 0.75);
	}
</style>
