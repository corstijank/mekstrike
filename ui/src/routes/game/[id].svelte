<script>
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	
	// Import refactored components
	import ScrollableBoard from '../../components/battlefield/ScrollableBoard.svelte';
	import UnitList from '../../components/units/UnitList.svelte';
	import GameMessageArea from '../../components/game/GameMessageArea.svelte';
	import UnitStatsPanel from '../../components/units/UnitStatsPanel.svelte';
	
	// Import UI components
	import ErrorBoundary from '../../components/ui/ErrorBoundary.svelte';
	import LoadingSpinner from '../../components/ui/LoadingSpinner.svelte';
	
	// Import stores and services
	import { initializeGame, cleanup, gameState } from '../../stores/gameStore.js';
	import { initializeBoard } from '../../stores/battlefieldStore.js';
	
	import 'terminal.css';

	let gameError = null;
	let gameId = '';
	let initialized = false;
	
	// Reactive game ID from route params
	$: gameId = $page.params.id;
	$: isLoading = $gameState.loading;

	// Only initialize once when gameId is available and not already initialized
	$: if (gameId && !initialized) {
		initializeGameData();
	}

	async function initializeGameData() {
		if (initialized) return;
		
		try {
			initialized = true;
			// Initialize game state and board
			await Promise.all([
				initializeGame(gameId),
				initializeBoard(gameId)
			]);
		} catch (error) {
			console.error('Failed to initialize game:', error);
			gameError = error;
			initialized = false; // Allow retry
		}
	}

	onMount(() => {
		// This will trigger the reactive initialization above
	});

	onDestroy(() => {
		// Clean up stores and polling when leaving the route
		cleanup();
		initialized = false; // Reset for next time
	});

	function handleRetry() {
		gameError = null;
		initialized = false; // Reset initialization flag to allow retry
		// This will trigger the reactive initialization
	}
</script>

<svelte:head>
	<title>Mekstrike - Game {gameId}</title>
</svelte:head>

{#if gameError}
	<div class="error-container">
		<ErrorBoundary error={gameError} on:retry={handleRetry} />
	</div>
{:else if isLoading}
	<div class="loading-container">
		<LoadingSpinner size="large" text="Loading game..." />
	</div>
{:else}
	<div class="game-layout">
		<div class="left-panel">
			<GameMessageArea gameId={gameId} />
		</div>
		
		<div class="center-panel">
			<ScrollableBoard id={gameId} />
		</div>
		
		<div class="right-panel">
			<div class="right-content">
				<div class="stats-section">
					<UnitStatsPanel />
				</div>
			</div>
		</div>
		
		<div class="bottom-panel">
			<UnitList id={gameId} />
		</div>
	</div>
{/if}

<style>
	.error-container, .loading-container {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100vh;
		width: 100vw;
		background-color: #111;
	}

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
		display: flex;
		flex-direction: column;
	}

	.right-content {
		height: 100%;
		display: flex;
		flex-direction: column;
	}

	.stats-section {
		flex: 1;
		min-height: 0;
		overflow: hidden;
	}

	.actions-section {
		flex-shrink: 0;
		border-top: 1px solid #333;
		max-height: 200px;
	}

	.bottom-panel {
		grid-area: bottom;
		border-top: 2px solid #333;
		background-color: #222;
		height: 155px;
		overflow: hidden;
	}

	/* Responsive adjustments for smaller screens */
	@media (max-width: 1200px) {
		.game-layout {
			grid-template-columns: 18vw 1fr 18vw;
		}
	}

	@media (max-width: 900px) {
		.game-layout {
			grid-template-columns: 15vw 1fr 15vw;
		}
		
		.actions-section {
			max-height: 150px;
		}
	}
</style>