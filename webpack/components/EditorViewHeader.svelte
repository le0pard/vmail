<script>
	import IconComponent from './Icon'
</script>

<style>
  .editor-header {
		display: flex;
		align-items: center;
		background-color: var(--headBgColor);
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
		box-shadow: 0 10px 30px 0 rgb(82 63 105 / 8%);
		border-bottom: 1px solid var(--headBorderColor);
	}

	.editor-header-item {
		display: flex;
		align-items: stretch;
		padding-top: 0;
		padding-bottom: 0;
	}

	.editor-header-link {
		color: var(--headColor);
		background-color: transparent;
		border-bottom: 2px solid transparent;
		transition: color 0.2s ease,background-color 0.2s ease;
		padding: 0.5rem 0;
		margin: 0 1rem;
		text-decoration: none;
	}

	.editor-header-link-active {
		border-bottom: 2px solid var(--headHoverColor);
		color: var(--headHoverColor);
	}

	.editor-header-link:hover, .editor-header-link:active {
		border-bottom: 2px solid var(--headHoverColor);
	}

	.editor-header-logo-container {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.editor-header-logo-link {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 2rem;
		height: 2rem;
		margin: 0 0.5rem;
	}
</style>

<div class="editor-header">
	<div class="editor-header-logo-container">
		<a class="editor-header-logo-link" href="/">
			<IconComponent />
		</a>
	</div>
	<div class="editor-header-item">
		<a class="editor-header-link editor-header-link-active" href="/">
			VMail
		</a>
	</div>
  <div class="editor-header-item">
		<a class="editor-header-link" href="/faq.html">FAQ</a>
	</div>
</div>
