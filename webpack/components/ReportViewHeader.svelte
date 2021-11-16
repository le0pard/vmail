<style>
  .report-header {
		display: flex;
		background-color: var(--headBgColor);
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
		box-shadow: 0 10px 30px 0 rgb(82 63 105 / 8%);
		border-bottom: 1px solid var(--headBorderColor);
		justify-content: flex-end;
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
