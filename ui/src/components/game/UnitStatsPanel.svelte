<script>
	import { selectedUnitData, clearSelection } from '../../stores/gameStores.js';

	function getDamageRange(shortdmg, meddmg, longdmg, ovhdmg) {
		return `${shortdmg}/${meddmg}/${longdmg}/${ovhdmg}`;
	}
</script>

<div class="stats-area">
	<div class="terminal-card">
		<header>
			Unit Stats
			{#if $selectedUnitData}
				<button class="btn btn-clear" on:click={clearSelection}>✕</button>
			{/if}
		</header>
		
		<div class="stats-container">
			{#if !$selectedUnitData}
				<div class="no-selection">
					<p>Click a unit to view its stats</p>
				</div>
			{:else}
				<div class="unit-stats">
					<div class="unit-art-section">
						<img src={$selectedUnitData.stats.image} alt={$selectedUnitData.stats.model} class="unit-art"/>
					</div>
					
					<div class="unit-header">
						<h3>{$selectedUnitData.stats.model}</h3>
						<p class="unit-type">{$selectedUnitData.stats.type} • {$selectedUnitData.stats.role}</p>
						<p class="unit-owner">Owner: {$selectedUnitData.owner}</p>
					</div>
					
					<div class="stats-grid">
						<div class="stat-row">
							<span class="stat-label">Structure:</span>
							<span class="stat-value">{$selectedUnitData.stats.struct}</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Armor:</span>
							<span class="stat-value">{$selectedUnitData.stats.armor}</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Movement:</span>
							<span class="stat-value">{$selectedUnitData.stats.movement}</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Damage:</span>
							<span class="stat-value">{getDamageRange($selectedUnitData.stats.shortdmg, $selectedUnitData.stats.meddmg, $selectedUnitData.stats.longdmg, $selectedUnitData.stats.ovhdmg)}</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Point Value:</span>
							<span class="stat-value">{$selectedUnitData.stats.pointvalue}</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Size:</span>
							<span class="stat-value">{$selectedUnitData.stats.size}</span>
						</div>
					</div>
					
					{#if $selectedUnitData.stats.specials && $selectedUnitData.stats.specials.length > 0}
						<div class="specials-section">
							<div class="stat-label">Special Abilities:</div>
							<div class="specials-list">
								{#each $selectedUnitData.stats.specials as special}
									<span class="special-tag">{special}</span>
								{/each}
							</div>
						</div>
					{/if}
					
					{#if $selectedUnitData.active}
						<div class="active-indicator">ACTIVE UNIT</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.stats-area {
		height: 100%;
		padding: 10px;
		overflow: hidden;
	}

	.stats-container {
		height: calc(100vh - 200px);
		overflow-y: auto;
		padding: 10px;
	}

	.no-selection {
		text-align: center;
		color: #666;
		padding: 40px 20px;
		font-style: italic;
	}

	.unit-stats {
		height: 100%;
		display: flex;
		flex-direction: column;
	}

	.unit-art-section {
		text-align: center;
		margin-bottom: 20px;
		flex-shrink: 0;
	}

	.unit-art {
		max-width: 90%;
		max-height: 200px;
		height: auto;
		width: auto;
		object-fit: contain;
		border: 3px solid #333;
		border-radius: 8px;
		background-color: #000;
	}

	.unit-header {
		text-align: center;
		margin-bottom: 20px;
		flex-shrink: 0;
	}

	.unit-header h3 {
		margin: 0 0 8px 0;
		font-size: 16px;
		color: #fff;
		word-wrap: break-word;
	}

	.unit-type {
		margin: 4px 0;
		font-size: 12px;
		color: #aaa;
	}

	.unit-owner {
		margin: 4px 0;
		font-size: 11px;
		color: #888;
	}

	.stats-grid {
		display: grid;
		gap: 12px;
		margin-bottom: 15px;
		flex-grow: 1;
	}

	.stat-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 8px 0;
		border-bottom: 1px solid #333;
	}

	.stat-label {
		font-weight: bold;
		color: #ccc;
		font-size: 13px;
	}

	.stat-value {
		color: #fff;
		font-family: 'Courier New', monospace;
		font-size: 13px;
		font-weight: bold;
	}

	.specials-section {
		margin-top: 15px;
		flex-shrink: 0;
	}

	.specials-list {
		display: flex;
		flex-wrap: wrap;
		gap: 6px;
		margin-top: 8px;
	}

	.special-tag {
		background-color: #444;
		color: #fff;
		padding: 4px 8px;
		border-radius: 4px;
		font-size: 10px;
		border: 1px solid #666;
	}

	.active-indicator {
		text-align: center;
		background-color: #00ff00;
		color: #000;
		padding: 8px;
		margin-top: 15px;
		font-weight: bold;
		border-radius: 4px;
		animation: pulse 2s infinite;
		flex-shrink: 0;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.6; }
	}

	.btn-clear {
		float: right;
		font-size: 12px;
		padding: 4px 8px;
		margin: -4px;
		background: none;
		border: 1px solid #666;
		color: #ccc;
	}

	.btn-clear:hover {
		background-color: #333;
		color: #fff;
	}
</style>