<script>
  import { onDestroy, setContext } from 'svelte'
  import { report, reportLoading, reportError } from '@stores/report'
  import { splitState } from '@stores/split'
  import EditorHeaderComponent from '@components/EditorViewHeader.svelte'
  import EditorViewComponent from '@components/EditorView.svelte'
  import SplitViewComponent from '@components/SplitView.svelte'
  import ReportHeaderComponent from '@components/ReportViewHeader.svelte'
  import ReportViewComponent from '@components/ReportView.svelte'

  let { webWorkerObject, githubSvgIcon } = $props()

  setContext('ww', {
    getWebWorker: () => webWorkerObject
  })

  const resetAllStates = () => {
    report.reset()
    reportLoading.reset()
    reportError.reset()
    splitState.reset()
  }

  onDestroy(resetAllStates)
</script>

<div class="parser-view">
  <div class="parser-editor" class:parser-editor-hidden={$splitState.visible === 'right'}>
    <EditorHeaderComponent />
    <EditorViewComponent />
  </div>
  <SplitViewComponent />
  <div class="parser-report" class:parser-report-hidden={$splitState.visible === 'left'}>
    <ReportHeaderComponent>
      {#snippet githubIcon()}
        {@render githubSvgIcon()}
      {/snippet}
    </ReportHeaderComponent>
    <ReportViewComponent />
  </div>
</div>

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

  .parser-report {
    height: 100%;
    flex: 1;
    position: relative;
    display: flex;
    flex-direction: column;
  }

  .parser-editor-hidden,
  .parser-report-hidden {
    display: none;
  }
</style>
