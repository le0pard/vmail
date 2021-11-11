<script>
  import {onMount} from 'svelte'
  import {report} from 'stores/report'
  import {
    MULTI_LEVEL_REPORT_KEYS,
    SINGLE_LEVEL_REPORT_KEYS,
    REPORT_CSS_VARIABLES,
    EVENT_LINE_TO_EDITOR,
    EVENT_LINE_TO_REPORT
  } from 'lib/constants'
  import ReportMultiItemComponent from './ReportMultiItem'
  import ReportSingleItemComponent from './ReportSingleItem'

  const handleLineClick = (line) => {
    window.dispatchEvent(new window.CustomEvent(EVENT_LINE_TO_EDITOR, {detail: {line}}))
  }

  const handleEditorLineClickEvent = (e) => {
    if (!e.detail?.line) {
      return
    }

    const {line} = e.detail
    console.log('Scroll to line: ', line)
  }

  onMount(() => {
    window.addEventListener(EVENT_LINE_TO_REPORT, handleEditorLineClickEvent)
    return () => window.removeEventListener(EVENT_LINE_TO_REPORT, handleEditorLineClickEvent)
  })
</script>

<style>
  .report-list {
    list-style-type: none;
    padding: 0;
    margin: 0.5rem 0.5rem 1rem 0.5rem;
  }
</style>

<ul class="report-list">
  {#each MULTI_LEVEL_REPORT_KEYS as reportInfo (reportInfo.key)}
    {#if $report[reportInfo.key]}
      {#each Object.keys($report[reportInfo.key]).sort() as itemName (itemName)}
        {#each Object.keys($report[reportInfo.key][itemName]).sort() as itemVal (itemVal)}
          <ReportMultiItemComponent
            reportInfo={reportInfo}
            itemName={itemName}
            itemVal={itemVal}
            report={$report[reportInfo.key][itemName][itemVal]}
            handleLineClick={handleLineClick}
          />
        {/each}
      {/each}
    {/if}
  {/each}

  {#if $report[REPORT_CSS_VARIABLES.key] && $report[REPORT_CSS_VARIABLES.key].lines.length > 0}
    <ReportSingleItemComponent
      reportInfo={REPORT_CSS_VARIABLES}
      itemName={'css_vars'}
      report={$report[REPORT_CSS_VARIABLES.key]}
      handleLineClick={handleLineClick}
    />
  {/if}

  {#each SINGLE_LEVEL_REPORT_KEYS as reportInfo (reportInfo.key)}
    {#if $report[reportInfo.key]}
      {#each Object.keys($report[reportInfo.key]).sort() as itemName (itemName)}
        <ReportSingleItemComponent
          reportInfo={reportInfo}
          itemName={itemName}
          report={$report[reportInfo.key][itemName]}
          handleLineClick={handleLineClick}
        />
      {/each}
    {/if}
  {/each}
</ul>
