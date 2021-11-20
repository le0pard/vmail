<script>
  import {onMount} from 'svelte'
  import {report, linesAndSelectors} from 'stores/report'
  import {splitState} from 'stores/split'
  import {camelize} from 'lib/reportHelpers'
  import {
    MULTI_LEVEL_REPORT_KEYS,
    SINGLE_LEVEL_REPORT_KEYS,
    REPORT_CSS_VARIABLES,
    EVENT_LINE_TO_EDITOR,
    EVENT_LINE_TO_REPORT
  } from 'lib/constants'
  import ReportItemComponent from './ReportItem'

  const genElementID = ([reportInfo, itemName, itemVal]) =>
    camelize(['item', reportInfo.key, itemName, itemVal].join('_')).replace(/_/g, '')

  const handleLineClick = (line) => {
    splitState.switchToLeftOnMobile()
    window.dispatchEvent(new window.CustomEvent(EVENT_LINE_TO_EDITOR, {detail: {line}}))
  }

  const handleEditorLineClickEvent = (e) => {
    if (!e.detail?.line) {
      return
    }

    const {line} = e.detail
    if (!$linesAndSelectors[line] || $linesAndSelectors[line].length === 0) {
      return
    }

    const scrollElementID = genElementID($linesAndSelectors[line][0])
    const scrollElement = document.getElementById(scrollElementID)
    if (!scrollElement) {
      return
    }

    setTimeout(() => {
      scrollElement.scrollIntoView({
        behavior: 'auto',
        block: 'center',
        inline: 'nearest'
      })
    }, 0)
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
          <ReportItemComponent
            reportInfo="{reportInfo}"
            itemName="{itemName}"
            itemVal="{itemVal}"
            elementID="{genElementID([reportInfo, itemName, itemVal])}"
            report="{$report[reportInfo.key][itemName][itemVal]}"
            handleLineClick="{handleLineClick}"
          />
        {/each}
      {/each}
    {/if}
  {/each}

  {#if $report[REPORT_CSS_VARIABLES.key] && $report[REPORT_CSS_VARIABLES.key].lines.length > 0}
    <ReportItemComponent
      reportInfo="{REPORT_CSS_VARIABLES}"
      itemName="{''}"
      itemVal="{''}"
      elementID="{genElementID([REPORT_CSS_VARIABLES, '', ''])}"
      report="{$report[REPORT_CSS_VARIABLES.key]}"
      handleLineClick="{handleLineClick}"
    />
  {/if}

  {#each SINGLE_LEVEL_REPORT_KEYS as reportInfo (reportInfo.key)}
    {#if $report[reportInfo.key]}
      {#each Object.keys($report[reportInfo.key]).sort() as itemName (itemName)}
        <ReportItemComponent
          reportInfo="{reportInfo}"
          itemName="{itemName}"
          itemVal="{''}"
          elementID="{genElementID([reportInfo, itemName, ''])}"
          report="{$report[reportInfo.key][itemName]}"
          handleLineClick="{handleLineClick}"
        />
      {/each}
    {/if}
  {/each}
</ul>
