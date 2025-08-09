<script>
	import { gameMessages } from '../../stores/gameStore.js';
	
	export let gameId;
	
	let messageContainer;

	// Auto-scroll to bottom when new messages arrive
	$: if (messageContainer && $gameMessages.length > 0) {
		setTimeout(() => {
			messageContainer.scrollTop = messageContainer.scrollHeight;
		}, 50);
	}

	function formatTime(timestamp) {
		return timestamp.toLocaleTimeString('en-US', { 
			hour12: false, 
			hour: '2-digit', 
			minute: '2-digit',
			second: '2-digit'
		});
	}

	function getMessageClass(type, data) {
		switch(type) {
			case 'com.dapr.event.sent':
				// For Dapr events, determine color by phase
				if (data && data.Phase === 'Movement') return 'message-movement';
				if (data && data.Phase === 'Combat') return 'message-combat';
				if (data && data.Phase === 'End') return 'message-info';
				return 'message-info';
			case 'combat': 
			case 'unit.attack.completed': 
			case 'unit-attack-completed': 
				return 'message-combat';
			case 'movement': 
			case 'unit.movement.completed': 
			case 'unit-movement-completed': 
				return 'message-movement';
			case 'unit.end.completed':
			case 'unit-end-phase-completed':
				return 'message-info';
			case 'system': 
				return 'message-system';
			case 'error': 
				return 'message-error';
			default: 
				return 'message-info';
		}
	}
</script>

<div class="message-area">
	<div class="terminal-card">
		<header>Combat Log</header>
		<div class="message-container" bind:this={messageContainer}>
			{#if $gameMessages.length === 0}
				<div class="no-messages">
					<p>No game events yet...</p>
				</div>
			{:else}
				{#each $gameMessages as message}
					<div class="message-entry {getMessageClass(message.type, message.rawData)}">
						<span class="timestamp">{formatTime(message.timestamp)}</span>
						<span class="message-text">{message.message}</span>
					</div>
				{/each}
			{/if}
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

	.no-messages {
		text-align: center;
		color: #666;
		padding: 40px 20px;
		font-style: italic;
	}

	.message-entry {
		margin-bottom: 8px;
		padding: 4px 0;
		border-left: 3px solid transparent;
		padding-left: 8px;
		word-wrap: break-word;
	}

	.timestamp {
		color: #666;
		margin-right: 8px;
		font-size: 10px;
	}

	.message-text {
		word-wrap: break-word;
		color: #ffffff;
	}

	.message-combat {
		border-left-color: #ff4444;
		background-color: rgba(255, 68, 68, 0.2);
		color: #ffcccc;
	}

	.message-movement {
		border-left-color: #4488ff;
		background-color: rgba(68, 136, 255, 0.2);
		color: #ccddff;
	}

	.message-system {
		border-left-color: #44ff44;
		background-color: rgba(68, 255, 68, 0.2);
		color: #ccffcc;
	}

	.message-error {
		border-left-color: #ff8844;
		background-color: rgba(255, 136, 68, 0.2);
		color: #ffddcc;
	}

	.message-info {
		border-left-color: #888;
		background-color: rgba(136, 136, 136, 0.2);
		color: #dddddd;
	}
</style>