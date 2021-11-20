<script>
	import ThemeSwitcherComponent from './ThemeSwitcher'
</script>

<div class="report-header">
	<div class="report-header-item">
		<ThemeSwitcherComponent />
	</div>
	<div class="report-header-item">
		<a
			class="report-header-github-link"
			href="https://github.com/le0pard/vmail"
			target="_blank"
			rel="noopener noreferrer"
		>
			<svg class="report-header-github">
				<use href="#github-svg-icon" />
			</svg>
		</a>
	</div>
</div>

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
		align-items: center;
		justify-content: center;
		padding-top: 0;
		padding-bottom: 0;
	}

	.report-header-github-link {
		color: var(--mutedButtonColor);
		margin: 0 1rem;
	}

	.report-header-github {
		fill: var(--mutedButtonColor);
		width: 2rem;
		height: 2rem;
		padding-top: 3px;
	}

	.report-header-github-link:hover .report-header-github {
		fill: var(--mutedButtonHoverColor);
	}
</style>
