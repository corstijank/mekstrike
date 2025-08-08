<script>
	import { gameMessages } from '../../stores/gameStores.js';
	import { onMount } from 'svelte';

	export let gameId;
	
	let messageContainer;

	$: if (messageContainer && $gameMessages.length > 0) {
		messageContainer.scrollTop = messageContainer.scrollHeight;
	}

	onMount(() => {
		gameMessages.update(() => [
			{ type: 'system', message: 'Game started', timestamp: new Date() },
			{ type: 'info', message: 'Waiting for player actions...', timestamp: new Date() }
		]);
	});

	function formatTime(timestamp) {
		return timestamp.toLocaleTimeString('en-US', { 
			hour12: false, 
			hour: '2-digit', 
			minute: '2-digit',
			second: '2-digit'
		});
	}

	function getMessageClass(type) {
		switch(type) {
			case 'combat': return 'message-combat';
			case 'movement': return 'message-movement';
			case 'system': return 'message-system';
			case 'error': return 'message-error';
			default: return 'message-info';
		}
	}
</script>

<div class="message-area">
	<div class="terminal-card">
		<header>Combat Log</header>
		<div class="message-container" bind:this={messageContainer}>
			{#each $gameMessages as message}
				<div class="message-entry {getMessageClass(message.type)}">
					<span class="timestamp">{formatTime(message.timestamp)}</span>
					<span class="message-text">{message.message}</span>
				</div>
			{/each}
		</div>
	</div>
</div>

<style>
	.message-area {
		height: 100%;
		padding: 10px;
		overflow: hidden;
	}

	.message-container {
		height: 60vh;
		overflow-y: auto;
		padding: 10px;
		font-family: 'Courier New', monospace;
		font-size: 12px;
		line-height: 1.4;
	}

	.message-entry {
		margin-bottom: 8px;
		padding: 4px 0;
		border-left: 3px solid transparent;
		padding-left: 8px;
	}

	.timestamp {
		color: #666;
		margin-right: 8px;
		font-size: 10px;
	}

	.message-text {
		word-wrap: break-word;
	}

	.message-combat {
		border-left-color: #ff4444;
		background-color: rgba(255, 68, 68, 0.1);
	}

	.message-movement {
		border-left-color: #4488ff;
		background-color: rgba(68, 136, 255, 0.1);
	}

	.message-system {
		border-left-color: #44ff44;
		background-color: rgba(68, 255, 68, 0.1);
	}

	.message-error {
		border-left-color: #ff8844;
		background-color: rgba(255, 136, 68, 0.1);
	}

	.message-info {
		border-left-color: #888;
		background-color: rgba(136, 136, 136, 0.1);
	}
</style>