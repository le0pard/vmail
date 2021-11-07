<script>
	import {onDestroy} from 'svelte'
	import {report, reportLoading, reportError} from 'stores/report'
	import EditorViewComponent from './EditorView'
	import ReportListComponent from './ReportList'

	export let parserFunction

	onDestroy(() => {
		report.reset()
		reportLoading.reset()
		reportError.reset()
	})
</script>

<style>
	.parser-view {
		display: flex;
		position: absolute;
		bottom: 0;
    left: 0;
    right: 0;
    top: 0;
	}

	.parser-editor {
		display: flex;
		flex-direction: column;
		height: 100%;
		flex: 1;
		position: relative;
	}

	.parser-editor-header {
		background-color: var(--headColor);
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
	}

	.parser-editor-header {
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
	}

	.parser-resize {
		display: flex;
		width: 30px;
	}

	.parser-report {
		height: 100%;
		flex: 1;
		position: relative;
		display: flex;
		flex-direction: column;
	}

	.parser-report-header {
		background-color: var(--headColor);
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
	}

	.parser-report-area {
		flex-grow: 1;
		position: relative;
		overflow: scroll;
	}
</style>

<div class="parser-view">
	<div class="parser-editor">
		<div class="parser-editor-header">
			<a href="/">Home</a>
			<a href="/faq.html">FAQ</a>
		</div>
		<EditorViewComponent parserFunction={parserFunction} />
	</div>
	<div class="parser-resize"></div>
	<div class="parser-report">
		<div class="parser-report-header">
			<p>Header</p>
		</div>
		<div class="parser-report-area">
			{#if Object.keys($report).length > 0}
				<ReportListComponent />
			{:else}
				<p>Submit your HTML</p>
			{/if}
		</div>
	</div>
</div>
