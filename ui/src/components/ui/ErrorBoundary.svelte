<script>
	export let error = null;
	export let showDetails = false;
	export let fallback = null;

	function toggleDetails() {
		showDetails = !showDetails;
	}

	function handleRetry() {
		error = null;
		// Dispatch retry event for parent to handle
		const event = new CustomEvent('retry');
		dispatchEvent(event);
	}
</script>

{#if error}
	<div class="error-boundary">
		<div class="error-container">
			<div class="error-icon">⚠️</div>
			<div class="error-content">
				<h3 class="error-title">Something went wrong</h3>
				<p class="error-message">
					{error.message || 'An unexpected error occurred'}
				</p>
				
				<div class="error-actions">
					<button class="btn btn-primary" on:click={handleRetry}>
						Try Again
					</button>
					{#if error.stack}
						<button class="btn btn-secondary" on:click={toggleDetails}>
							{showDetails ? 'Hide' : 'Show'} Details
						</button>
					{/if}
				</div>
				
				{#if showDetails && error.stack}
					<details class="error-details">
						<summary>Error Details</summary>
						<pre class="error-stack">{error.stack}</pre>
					</details>
				{/if}
			</div>
		</div>
	</div>
{:else if fallback}
	{@html fallback}
{:else}
	<slot />
{/if}

<style>
	.error-boundary {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 200px;
		padding: 20px;
		background-color: rgba(255, 68, 68, 0.1);
		border: 1px solid rgba(255, 68, 68, 0.3);
		border-radius: 8px;
		margin: 10px;
	}

	.error-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		max-width: 500px;
	}

	.error-icon {
		font-size: 48px;
		margin-bottom: 16px;
	}

	.error-content {
		width: 100%;
	}

	.error-title {
		margin: 0 0 8px 0;
		font-size: 20px;
		color: #ff4444;
		font-weight: bold;
	}

	.error-message {
		margin: 0 0 20px 0;
		color: #ccc;
		font-size: 14px;
		line-height: 1.4;
	}

	.error-actions {
		display: flex;
		gap: 12px;
		justify-content: center;
		margin-bottom: 20px;
	}

	.btn {
		padding: 8px 16px;
		border-radius: 4px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		border: 1px solid;
		transition: all 0.2s ease;
	}

	.btn-primary {
		background-color: #62c4ff;
		border-color: #62c4ff;
		color: #000;
	}

	.btn-primary:hover {
		background-color: #4db8ff;
	}

	.btn-secondary {
		background-color: transparent;
		border-color: #666;
		color: #ccc;
	}

	.btn-secondary:hover {
		background-color: #333;
		border-color: #777;
		color: #fff;
	}

	.error-details {
		text-align: left;
		width: 100%;
		margin-top: 16px;
	}

	.error-details summary {
		cursor: pointer;
		color: #ccc;
		font-size: 12px;
		margin-bottom: 8px;
	}

	.error-stack {
		background-color: #1a1a1a;
		border: 1px solid #333;
		border-radius: 4px;
		padding: 12px;
		font-size: 11px;
		color: #999;
		overflow-x: auto;
		white-space: pre-wrap;
		word-wrap: break-word;
	}
</style>