<script>
	export let variant = 'default'; // 'default', 'primary', 'secondary', 'danger'
	export let size = 'medium'; // 'small', 'medium', 'large'
	export let disabled = false;
	export let loading = false;
	export let fullWidth = false;

	// Handle click event
	function handleClick(event) {
		if (disabled || loading) {
			event.preventDefault();
			return;
		}
	}
</script>

<button
	class="btn btn-{variant} btn-{size} {fullWidth ? 'btn-full-width' : ''} {loading ? 'btn-loading' : ''}"
	{disabled}
	on:click={handleClick}
	on:click
	{...$$restProps}
>
	{#if loading}
		<div class="spinner"></div>
	{/if}
	<slot />
</button>

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
		border: 1px solid;
		border-radius: 4px;
		font-family: inherit;
		font-weight: 500;
		text-align: center;
		text-decoration: none;
		cursor: pointer;
		transition: all 0.2s ease;
		position: relative;
		white-space: nowrap;
		user-select: none;
	}

	.btn:disabled {
		cursor: not-allowed;
		opacity: 0.6;
	}

	/* Sizes */
	.btn-small {
		padding: 4px 8px;
		font-size: 12px;
		min-height: 24px;
	}

	.btn-medium {
		padding: 8px 16px;
		font-size: 14px;
		min-height: 32px;
	}

	.btn-large {
		padding: 12px 20px;
		font-size: 16px;
		min-height: 40px;
	}

	.btn-full-width {
		width: 100%;
	}

	/* Variants */
	.btn-default {
		background-color: #333;
		border-color: #555;
		color: #fff;
	}

	.btn-default:hover:not(:disabled) {
		background-color: #444;
		border-color: #666;
	}

	.btn-primary {
		background-color: #62c4ff;
		border-color: #62c4ff;
		color: #000;
	}

	.btn-primary:hover:not(:disabled) {
		background-color: #4db8ff;
		border-color: #4db8ff;
	}

	.btn-secondary {
		background-color: transparent;
		border-color: #666;
		color: #ccc;
	}

	.btn-secondary:hover:not(:disabled) {
		background-color: #333;
		border-color: #777;
		color: #fff;
	}

	.btn-danger {
		background-color: #ff3c74;
		border-color: #ff3c74;
		color: #fff;
	}

	.btn-danger:hover:not(:disabled) {
		background-color: #e6355f;
		border-color: #e6355f;
	}

	/* Loading state */
	.btn-loading {
		color: transparent;
	}

	.spinner {
		position: absolute;
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top: 2px solid currentColor;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}
</style>