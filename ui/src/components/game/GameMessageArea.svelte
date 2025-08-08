<script>
	import { gameMessages } from '../../stores/gameStores.js';
	import { onMount } from 'svelte';

	export let gameId;
	
	let messageContainer;

	$: if (messageContainer && $gameMessages.length > 0) {
		messageContainer.scrollTop = messageContainer.scrollHeight;
	}

	async function fetchGameLogs() {
		if (!gameId) return;
		
		try {
			const response = await fetch(`/mekstrike/api/gamemaster/games/${gameId}/logs`);
			if (response.ok) {
				const logs = await response.json();
				gameMessages.set(logs.map(log => {
					// Try to parse the message as JSON to see if it's an enriched event
					let parsedData = log.message;
					try {
						parsedData = JSON.parse(log.message);
					} catch (e) {
						// If it's not JSON, keep as is
					}
					
					// Format the message based on event type
					const formattedMessage = formatEventMessage(parsedData, log.type);
					
					return {
						type: log.type,
						message: formattedMessage,
						timestamp: new Date(log.timestamp),
						rawData: parsedData
					};
				}));
			}
		} catch (error) {
			console.error('Error fetching game logs:', error);
		}
	}

	function formatMovementEvent(data) {
		if (data.Unit && data.SourceLocation && data.TargetLocation) {
			const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
			const player = data.Unit.Owner || 'Unknown player';
			const source = `(${data.SourceLocation.x}, ${data.SourceLocation.y})`;
			const target = `(${data.TargetLocation.x}, ${data.TargetLocation.y})`;
			return `${player} moved ${unitName} from ${source} to ${target}`;
		}
		return `Unit ${data.UnitId || 'unknown'} completed movement phase`;
	}
	
	function formatAttackEvent(data) {
		if (data.Unit && data.SourceLocation && data.TargetId) {
			const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
			const player = data.Unit.Owner || 'Unknown player';
			const source = `(${data.SourceLocation.x}, ${data.SourceLocation.y})`;
			return `${player}'s ${unitName} at ${source} attacked target ${data.TargetId}`;
		}
		return `Unit ${data.UnitId || 'unknown'} completed combat phase`;
	}
	
	function formatEndPhaseEvent(data) {
		if (data.Unit) {
			const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
			const player = data.Unit.Owner || 'Unknown player';
			const location = data.SourceLocation ? `(${data.SourceLocation.x}, ${data.SourceLocation.y})` : 'unknown location';
			return `${player}'s ${unitName} at ${location} ended its turn`;
		}
		return `Unit ${data.UnitId || 'unknown'} completed end phase`;
	}

	function formatEventMessage(data, type) {
		// If it's not a parsed object, return as is
		if (typeof data !== 'object') {
			return data;
		}

		// Handle different event types with stylized messages
		// Also check for Dapr-style topic names
		switch (type) {
			case 'com.dapr.event.sent':
				// For Dapr events, determine the actual event type from the data
				if (data.Phase === 'Movement') {
					return formatMovementEvent(data);
				} else if (data.Phase === 'Combat') {
					return formatAttackEvent(data);
				} else if (data.Phase === 'End') {
					return formatEndPhaseEvent(data);
				}
				// Fallback for unknown phases
				return `Unit ${data.UnitId || 'unknown'} completed ${data.Phase || 'unknown'} phase`;
			case 'unit.movement.completed':
			case 'unit-movement-completed':
				if (data.Unit && data.SourceLocation && data.TargetLocation) {
					const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
					const player = data.Unit.Owner || 'Unknown player';
					const source = `(${data.SourceLocation.x}, ${data.SourceLocation.y})`;
					const target = `(${data.TargetLocation.x}, ${data.TargetLocation.y})`;
					return `${player} moved ${unitName} from ${source} to ${target}`;
				}
				return `Unit ${data.UnitId || 'unknown'} completed movement phase`;
				
			case 'unit.attack.completed':
			case 'unit-attack-completed':
				if (data.Unit && data.SourceLocation && data.TargetId) {
					const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
					const player = data.Unit.Owner || 'Unknown player';
					const source = `(${data.SourceLocation.x}, ${data.SourceLocation.y})`;
					return `${player}'s ${unitName} at ${source} attacked target ${data.TargetId}`;
				}
				return `Unit ${data.UnitId || 'unknown'} completed combat phase`;
				
			case 'unit.end.completed':
			case 'unit-end-phase-completed':
				if (data.Unit) {
					const unitName = data.Unit.Model || data.Unit.Id || 'Unknown unit';
					const player = data.Unit.Owner || 'Unknown player';
					const location = data.SourceLocation ? `(${data.SourceLocation.x}, ${data.SourceLocation.y})` : 'unknown location';
					return `${player}'s ${unitName} at ${location} ended its turn`;
				}
				return `Unit ${data.UnitId || 'unknown'} completed end phase`;
				
			case 'movement':
				// Legacy format
				return data;
				
			case 'combat':
				// Legacy format
				return data;
				
			case 'system':
				// System messages
				return data;
				
			default:
				// For any other type, try to create a meaningful message
				if (typeof data === 'string') {
					return data;
				}
				// If it's an object but we don't know how to format it, show a simplified version
				return JSON.stringify(data);
		}
	}

	onMount(() => {
		fetchGameLogs();
		
		// Poll for updates every 3 seconds
		const interval = setInterval(fetchGameLogs, 3000);
		
		return () => clearInterval(interval);
	});

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
			{#each $gameMessages as message}
				<div class="message-entry {getMessageClass(message.type, message.rawData)}">
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