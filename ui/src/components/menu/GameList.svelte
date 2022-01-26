<script>
	import { DateTime } from 'luxon';

	export let games = [{ StartTime: DateTime.now(), PlayerA: '', PlayerB: '',BattlefieldID:"BattlefieldID" }];	
</script>

<main>
	<div class="terminal-card">
		<header>Active games</header>
		<div class="background-transparent">
			<nav>
				<ul>
					{#each games as game}
						<li>
							<a sveltekit:prefetch href="/game/{game.BattlefieldID}">	{DateTime.fromISO(game.StartTime).toLocaleString(DateTime.DATETIME_SHORT)}: {game.PlayerA}({game.PlayerAValue})
								vs. {game.PlayerB}({game.PlayerBValue})</a>
						</li>
					{/each}
				</ul>
			</nav>
		</div>
	</div>
</main>

<style>
	.foreground-blue {
		color: #62c4ff;
	}
	.background-transparent {
		background: rgba(0, 0, 0, 0.75);
	}
</style>
