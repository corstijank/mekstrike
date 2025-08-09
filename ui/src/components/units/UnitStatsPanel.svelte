<script>
	import { selectedUnitData, clearSelection } from '../../stores/unitStore.js';

	function getDamageRange(shortdmg, meddmg, longdmg, ovhdmg) {
		return `${shortdmg}/${meddmg}/${longdmg}/${ovhdmg}`;
	}

	function formatUnitType(type, role) {
		if (type && role) {
			return `${type} • ${role}`;
		}
		return type || role || 'Unknown';
	}

	function formatSize(size) {
		const sizeMap = { 1: 'Light', 2: 'Medium', 3: 'Heavy', 4: 'Assault' };
		return sizeMap[size] || size || 'Unknown';
	}

	function formatFacing(facing) {
		const facingMap = { 0: 'N', 1: 'NE', 2: 'SE', 3: 'S', 4: 'SW', 5: 'NW' };
		return facingMap[facing] || facing;
	}

	function formatMovement(movement) {
		if (!movement) return { hexes: 0, canJump: false };
		
		const movementStr = String(movement);
		const canJump = movementStr.includes('j');
		const inches = parseInt(movementStr.replace(/[^0-9]/g, '')) || 0;
		const hexes = Math.floor(inches / 2);
		
		return { hexes, canJump };
	}

	function hasSpecials(specials) {
		return specials && Array.isArray(specials) && specials.length > 0;
	}
</script>

<div class="stats-area">
	<div class="terminal-card">
		<header>
			Unit Stats
		</header>
		
		<div class="stats-container">
			{#if !$selectedUnitData}
				<div class="no-selection">
					<p>Click a unit to view its stats</p>
				</div>
			{:else}
				<div class="unit-stats">
					<div class="unit-art-section">
						{#if $selectedUnitData.stats.image}
							<img 
								src={$selectedUnitData.stats.image} 
								alt={$selectedUnitData.stats.model} 
								class="unit-art"
								on:error={(e) => { e.target.style.display = 'none'; }}
							/>
						{:else}
							<div class="no-image">
								<span>No Image Available</span>
							</div>
						{/if}
					</div>
					
					<div class="unit-header">
						<h3>{$selectedUnitData.stats.model || 'Unknown Model'}</h3>
						<p class="unit-type">{formatUnitType($selectedUnitData.stats.type, $selectedUnitData.stats.role)}</p>
						<p class="unit-owner">Owner: {$selectedUnitData.owner || 'Unknown'}</p>
					</div>
					
					<div class="stats-grid">
						<div class="stat-row">
							<span class="stat-label">Structure:</span>
							<span class="stat-value">{$selectedUnitData.stats.struct || 0}</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Armor:</span>
							<span class="stat-value">{$selectedUnitData.stats.armor || 0}</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Movement:</span>
							<span class="movement-compact">
								<span class="stat-value">{formatMovement($selectedUnitData.stats.movement).hexes}</span>
								{#if formatMovement($selectedUnitData.stats.movement).canJump}
									<span class="jump-indicator">⬆</span>
								{/if}
							</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Damage:</span>
							<span class="damage-compact">
								<span class="dmg-item">S:{$selectedUnitData.stats.shortdmg || 0}</span>
								<span class="dmg-item">M:{$selectedUnitData.stats.meddmg || 0}</span>
								<span class="dmg-item">L:{$selectedUnitData.stats.longdmg || 0}</span>
								<span class="dmg-item">OVH:{$selectedUnitData.stats.ovhdmg || 0}</span>
							</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Point Value:</span>
							<span class="stat-value">{$selectedUnitData.stats.pointvalue || 0}</span>
						</div>
						<div class="stat-row">
							<span class="stat-label">Size:</span>
							<span class="stat-value">{formatSize($selectedUnitData.stats.size)}</span>
						</div>
					</div>
					
					{#if hasSpecials($selectedUnitData.stats.specials)}
						<div class="specials-section">
							<div class="stat-label">Special Abilities:</div>
							<div class="specials-list">
								{#each $selectedUnitData.stats.specials as special}
									<span class="special-tag">{special}</span>
								{/each}
							</div>
						</div>
					{/if}

					{#if $selectedUnitData.location}
						<div class="location-info">
							<div class="stat-label">Position:</div>
							<div class="location-value">
								({$selectedUnitData.location.position.x}, {$selectedUnitData.location.position.y})
								{#if $selectedUnitData.location.heading !== undefined}
									• {formatFacing($selectedUnitData.location.heading)}
								{/if}
							</div>
						</div>
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

	.no-image {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 150px;
		height: 100px;
		margin: 0 auto;
		background-color: #222;
		border: 2px dashed #555;
		border-radius: 8px;
		color: #777;
		font-size: 12px;
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
		gap: 6px;
		margin-bottom: 10px;
		flex-grow: 1;
	}

	.stat-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 4px 0;
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

	.damage-compact {
		display: flex;
		gap: 8px;
		font-family: 'Courier New', monospace;
		font-size: 11px;
		color: #fff;
	}

	.dmg-item {
		font-weight: bold;
	}

	.movement-compact {
		display: flex;
		align-items: center;
		gap: 4px;
	}

	.jump-indicator {
		color: #4CAF50;
		font-size: 12px;
		font-weight: bold;
	}

	.specials-section {
		margin-top: 8px;
		flex-shrink: 0;
	}

	.specials-list {
		display: flex;
		flex-wrap: wrap;
		gap: 4px;
		margin-top: 4px;
	}

	.special-tag {
		background-color: #444;
		color: #fff;
		padding: 2px 6px;
		border-radius: 3px;
		font-size: 9px;
		border: 1px solid #666;
	}

	.location-info {
		margin-top: 8px;
		flex-shrink: 0;
		padding: 4px 6px;
		background-color: rgba(0, 0, 0, 0.2);
		border-radius: 3px;
	}

	.location-value {
		color: #fff;
		font-family: 'Courier New', monospace;
		font-size: 12px;
		margin-top: 4px;
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
		cursor: pointer;
		border-radius: 3px;
		transition: all 0.2s ease;
	}

	.btn-clear:hover {
		background-color: #333;
		color: #fff;
		border-color: #999;
	}
</style>