<style>
  .report-header {
		display: flex;
		background-color: var(--headBgColor);
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
		box-shadow: 0 2px 2px rgb(0 0 0 / 3%), 0 1px 0 rgb(0 0 0 / 3%);
	}

	.report-header-item {
		display: flex;
		align-items: stretch;
		padding-top: 0;
		padding-bottom: 0;
	}

	.report-header-link {
		color: var(--headColor);
		background-color: transparent;
		border-bottom: 2px solid transparent;
		transition: color 0.2s ease,background-color 0.2s ease;
		padding: 0.5rem 0;
		margin: 0 1rem;
		text-decoration: none;
	}

	.report-header-link-active {
		border-bottom: 2px solid var(--headHoverColor);
		color: var(--headHoverColor);
	}

	.report-header-link:hover, .report-header-link:active {
		border-bottom: 2px solid var(--headHoverColor);
	}
</style>

<div class="report-header">
	<div class="report-header-item">
		<a class="report-header-link report-header-link-active" href="/">Home</a>
	</div>
  <div class="report-header-item">
		<a class="report-header-link" href="/faq.html">FAQ</a>
	</div>
</div>
