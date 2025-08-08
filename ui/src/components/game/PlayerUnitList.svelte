<script>
	import { onMount } from 'svelte';
	import SmallUnitIcon from './SmallUnitIcon.svelte';
	import { selectedUnit } from '../../stores/gameStores.js';

	export let id;

	let gamedata = { PlayerAUnits: [], PlayerBUnits: [], CurrentRound: 0, CurrentPhase: 0, ActivePlayer: 0 };
	let availableActions = { UnitOwner: '', CurrentUnitID: '', ActionType: '' };

	function refreshData() {
		fetch('/mekstrike/api/gamemaster/games/' + id)
			.then((response) => response.json())
			.then((data) => {
				gamedata = data;
			});
		
		fetch('/mekstrike/api/gamemaster/games/' + id + '/availableActions')
			.then((response) => response.json())
			.then((data) => {
				availableActions = data;
			});
	}

	onMount(() => {
		refreshData();
		const interval = setInterval(refreshData, 2000);
		return () => clearInterval(interval);
	});

	const phases = ['Movement', 'Combat', 'End'];
	const players = ['Player A', 'Player B'];
</script>

<div class="bottom-container">
	<div class="roster-header">
		<h4>Army Overview</h4>
		<div class="game-status">
			Round {gamedata.CurrentRound} • {phases[gamedata.CurrentPhase]} Phase • {players[gamedata.ActivePlayer]}'s Turn
		</div>
	</div>
	
	<div class="unit-roster">
		<div class="player-units">
			<span class="player-label">Player A:</span>
			<div class="unitlist">
				{#each gamedata.PlayerAUnits as unitID}
					<SmallUnitIcon {unitID} gameId={id} isActive={availableActions.CurrentUnitID === unitID} />
				{/each}
			</div>
		</div>
		
		<div class="player-units">
			<span class="player-label">Player B:</span>
			<div class="unitlist">
				{#each gamedata.PlayerBUnits as unitID}
					<SmallUnitIcon {unitID} gameId={id} isActive={availableActions.CurrentUnitID === unitID} />
				{/each}
			</div>
		</div>
	</div>
</div>

<style>
	.bottom-container {
		height: 100%;
		background-color: #222;
		padding: 10px;
		display: flex;
		flex-direction: column;
	}

	.roster-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 10px;
		padding-bottom: 5px;
		border-bottom: 1px solid #333;
	}

	.roster-header h4 {
		margin: 0;
		color: #fff;
		font-size: 14px;
	}

	.game-status {
		font-size: 11px;
		color: #aaa;
		font-family: 'Courier New', monospace;
	}

	.unit-roster {
		display: flex;
		gap: 20px;
		align-items: flex-start;
		flex: 1;
		overflow: hidden;
	}

	.player-units {
		display: flex;
		flex-direction: column;
		gap: 8px;
		flex: 1;
		min-width: 0;
	}

	.player-label {
		font-size: 12px;
		font-weight: bold;
		color: #ccc;
		margin-bottom: 5px;
	}

	.unitlist {
        display: flex; 
		flex-direction: row;
        justify-content: flex-start;
        column-gap: 6px;
		flex-wrap: nowrap;
		overflow-x: auto;
		overflow-y: hidden;
		padding-bottom: 5px;
	}

	.unitlist::-webkit-scrollbar {
		height: 6px;
	}

	.unitlist::-webkit-scrollbar-track {
		background: #333;
		border-radius: 3px;
	}

	.unitlist::-webkit-scrollbar-thumb {
		background: #555;
		border-radius: 3px;
	}
</style>
