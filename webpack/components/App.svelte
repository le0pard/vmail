<svelte:options immutable="{true}" />

<script>
  import {wrap, releaseProxy} from 'comlink'
  import {onMount, onDestroy, setContext} from 'svelte'
  import {report, reportLoading, reportError} from 'stores/report'
  import {splitState} from 'stores/split'
  import EditorHeaderComponent from './EditorViewHeader'
  import EditorViewComponent from './EditorView'
  import SplitViewComponent from './SplitView'
  import ReportHeaderComponent from './ReportViewHeader'
  import ReportViewComponent from './ReportView'

  export let workerURL

  let workerObject = null

  setContext('ww', {
    getWebWorker: () => workerObject
  })

  onMount(() => {
    const worker = new Worker(workerURL)
    workerObject = wrap(worker)

    return () => {
      if (workerObject) {
        workerObject[releaseProxy]()
      }
      if (worker?.terminate) {
        worker.terminate()
      }
    }
  })

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

<div class="parser-view">
  <div class="parser-editor" class:parser-editor-hidden="{$splitState.visible === 'right'}">
    <EditorHeaderComponent />
    <EditorViewComponent />
  </div>
  <SplitViewComponent />
  <div class="parser-report" class:parser-report-hidden="{$splitState.visible === 'left'}">
    <ReportHeaderComponent />
    <ReportViewComponent />
  </div>
</div>
