<script>
	export let title = '';
	export let variant = 'default'; // 'default', 'terminal'
	export let padding = 'normal'; // 'none', 'small', 'normal', 'large'
	export let headerActions = false;
</script>

<div class="card card-{variant} padding-{padding}">
	{#if title || headerActions || $$slots.header}
		<header class="card-header">
			{#if $$slots.header}
				<slot name="header" />
			{:else}
				<h3 class="card-title">{title}</h3>
			{/if}
			{#if headerActions && $$slots.actions}
				<div class="card-actions">
					<slot name="actions" />
				</div>
			{/if}
		</header>
	{/if}
	
	<div class="card-content">
		<slot />
	</div>
</div>

<style>
	.card {
		background-color: #222;
		border-radius: 4px;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		height: 100%;
	}

	.card-default {
		border: 1px solid #333;
	}

	.card-terminal {
		border: 1px solid #333;
		background-color: #1a1a1a;
		font-family: 'Courier New', monospace;
	}

	.card-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 12px 16px;
		border-bottom: 1px solid #333;
		background-color: rgba(0, 0, 0, 0.2);
		flex-shrink: 0;
	}

	.card-title {
		margin: 0;
		font-size: 14px;
		font-weight: bold;
		color: #fff;
	}

	.card-actions {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.card-content {
		flex: 1;
		overflow: hidden;
		display: flex;
		flex-direction: column;
	}

	/* Padding variants */
	.padding-none .card-content {
		padding: 0;
	}

	.padding-small .card-content {
		padding: 8px;
	}

	.padding-normal .card-content {
		padding: 16px;
	}

	.padding-large .card-content {
		padding: 24px;
	}

	/* Terminal variant adjustments */
	.card-terminal .card-header {
		background-color: #333;
		color: #00ff00;
		font-weight: bold;
		font-size: 12px;
	}

	.card-terminal .card-title {
		color: #00ff00;
		text-transform: uppercase;
	}
</style>