<script>
    import GameList from "./GameList.svelte";
    import NewGame from "./NewGame.svelte";

    import { onMount } from 'svelte';
	import { DateTime } from 'luxon';

	let games = [{ StartTime: DateTime.now(), PlayerA: '', PlayerB: '',BattlefieldID:"BattlefieldID" }];
	onMount(() => {
		loadGames();
	});
    function loadGames(){
        fetch('/mekstrike/api/gamemaster/games')
			.then((response) => {
				return response.json();
			})
			.then((data) => {
				console.log(games);
				games = [...data];
			});
    }
</script>

<main>
    <img src="/bg.jpg" alt="bgimg" class="bg"/>
    <div class="mainmenu">
        <img src="/header.png" alt="header"/>
        <GameList bind:games={games} />
        <br/>
        <NewGame on:NewGame={loadGames}/>
    </div>
</main>

<style>
    .mainmenu {
        margin: auto !important;
        width: 40vw;
        height: 60vh;
    }
    .bg {
        position: absolute;
        top: 0;
        left: 0;
        width: 100vw;
        height: 100vh;
        object-fit: cover;
        object-position: top;
        z-index: -1;
    }
</style>
