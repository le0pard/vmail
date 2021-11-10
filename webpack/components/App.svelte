<script>
	import {onDestroy} from 'svelte'
	import {report, reportLoading, reportError} from 'stores/report'
	import {splitState} from 'stores/split'
	import EditorViewComponent from './EditorView'
	import SplitViewComponent from './SplitView'
	import ReportViewComponent from './ReportView'

	export let parserFunction

	onDestroy(() => {
		report.reset()
		reportLoading.reset()
		reportError.reset()
		splitState.reset()
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
		box-shadow: 0 2px 2px rgb(0 0 0 / 3%), 0 1px 0 rgb(0 0 0 / 3%);
		padding: 0.3rem 0;
	}

	.parser-editor-header {
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
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
		box-shadow: 0 2px 2px rgb(0 0 0 / 3%), 0 1px 0 rgb(0 0 0 / 3%);
		padding: 0.3rem 0;
	}

	.parser-editor-hidden {
		display: none;
	}

	.parser-report-hidden {
		display: none;
	}
</style>

<div class="parser-view">
	<div class="parser-editor" class:parser-editor-hidden="{$splitState.visible === 'right'}">
		<div class="parser-editor-header">
			<a href="/">Home</a>
			<a href="/faq.html">FAQ</a>
		</div>
		<EditorViewComponent parserFunction={parserFunction} />
	</div>
	<SplitViewComponent />
	<div class="parser-report"  class:parser-report-hidden="{$splitState.visible === 'left'}">
		<div class="parser-report-header">
			<a href="/">Home</a>
			<a href="/faq.html">FAQ</a>
		</div>
		<ReportViewComponent />
	</div>
</div>
