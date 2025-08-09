<script>
	export let size = 'medium'; // 'small', 'medium', 'large'
	export let color = 'primary'; // 'primary', 'secondary', 'white'
	export let text = '';
	export let inline = false;
</script>

<div class="spinner-container {inline ? 'spinner-inline' : 'spinner-block'}">
	<div class="spinner spinner-{size} spinner-{color}"></div>
	{#if text}
		<div class="spinner-text">{text}</div>
	{/if}
</div>

<style>
	.spinner-container {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
	}

	.spinner-block {
		flex-direction: column;
		padding: 20px;
	}

	.spinner-inline {
		flex-direction: row;
	}

	.spinner {
		border-radius: 50%;
		animation: spin 1s linear infinite;
		border-style: solid;
	}

	/* Sizes */
	.spinner-small {
		width: 16px;
		height: 16px;
		border-width: 2px;
	}

	.spinner-medium {
		width: 24px;
		height: 24px;
		border-width: 3px;
	}

	.spinner-large {
		width: 32px;
		height: 32px;
		border-width: 4px;
	}

	/* Colors */
	.spinner-primary {
		border-color: rgba(98, 196, 255, 0.3);
		border-top-color: #62c4ff;
	}

	.spinner-secondary {
		border-color: rgba(170, 170, 170, 0.3);
		border-top-color: #aaa;
	}

	.spinner-white {
		border-color: rgba(255, 255, 255, 0.3);
		border-top-color: #fff;
	}

	.spinner-text {
		font-size: 14px;
		color: #ccc;
		text-align: center;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}
</style>