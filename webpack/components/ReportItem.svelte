<script>
  import {getContext} from 'svelte'
  import {normalizeItemVal, reportStats, clientsList} from 'lib/report-helpers'

  export let reportInfo
  export let itemName
  export let itemVal
  export let report
  export let handleLineClick

  let normalizedClientsList = null
  let itemStats = null

  const {getWebWorker} = getContext('ww')
  const webWorker = getWebWorker()

  if (webWorker?.getStatsAndClients) {
    webWorker.getStatsAndClients(report.rules).then(([resStats, resClients]) => {
      itemStats = resStats
      normalizedClientsList = resClients
    })
  } else {
    itemStats = reportStats(report.rules)
    normalizedClientsList = clientsList(report.rules)
  }
</script>

<style>
  .report-item {
    background-color: var(--headBgColor);
    padding: 0.5rem 1rem;
    border: 1px solid var(--borderColor);
  }

  .report-item:first-child {
    border-top-color: var(--borderListColor);
    border-top-right-radius: 0.25rem;
    border-top-left-radius: 0.25rem;
  }

  .report-item:last-child {
    border-bottom-color: var(--borderListColor);
    border-bottom-right-radius: 0.25rem;
    border-bottom-left-radius: 0.25rem;
  }

  .report-container {
    display: flex;
    flex-direction: column;
  }

  .report-header {
    display: flex;
    flex-direction: row;
  }

  .report-header-main {
    flex: 4;
    display: flex;
    flex-direction: column;
  }

  .report-header-main-top {
    flex: 1;
    display: flex;
    flex-direction: row;
  }

  .report-header-main-lines {
    flex: 1;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: flex-start;
  }

  .report-header-type {
    flex: 1;
  }

  .report-header-name {
    flex: 2;
  }

  .report-header-score {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .report-header-score-chart {
    display: inline-flex;
    height: 1.5rem;
  }

  .report-header-score-supported {
    background-color: #39b54a;
    flex-basis: auto;
    flex-shrink: 0;
    flex-grow: 0;
  }

  .report-header-score-mitigated {
    background-color: #eda745;
    flex-basis: auto;
    flex-shrink: 0;
    flex-grow: 0;
  }

  .report-header-score-unsupported {
    background-color: #c44230;
    flex-basis: auto;
    flex-shrink: 0;
    flex-grow: 0;
  }

  .report-header-score-summary-supported-value {
    background-color: #39b54a;
    display: inline-block;
    padding: 0 0.25rem;
    line-height: 1.4rem;
  }

  .report-header-score-summary-mitigated-value {
    background-color: #eda745;
    display: inline-block;
    padding: 0 0.25rem;
    line-height: 1.4rem;
  }

  .report-line-button {
    background: none;
    border: none;
    cursor: pointer;
    text-decoration: underline;
    font-size: 0.9rem;
  }

  .report-line-button:hover {
    text-decoration: none;
  }
</style>

<li class="report-item">
  <div class="report-container">
    <div class="report-header">
      <div class="report-header-main">
        <div class="report-header-main-top">
          <div class="report-header-type">{reportInfo.title}</div>
          <div class="report-header-name">
            {itemName}
            {#if itemVal.length > 0}
              {normalizeItemVal(itemVal)}
            {/if}
          </div>
        </div>
        <div class="report-header-main-lines">
          {#each report.lines as line, i}
            <button on:click|preventDefault={() => handleLineClick(line)} class="report-line-button">{line}</button>
            {#if i < report.lines.length - 1},{/if}
          {/each}
          {#if report.more_lines}
            and more...
          {/if}
        </div>
      </div>
      {#if itemStats}
        <div class="report-header-score">
          <div class="report-header-score-chart">
            {#if itemStats.supportedPercentage > 0}
              <div tabindex="0" title="{itemStats.supportedPercentage}% supported" role="group" style="width:{itemStats.supportedPercentage}%;" class="report-header-score-supported"></div>
            {/if}
            {#if itemStats.mitigatedPercentage > 0}
              <div tabindex="0" title="{itemStats.mitigatedPercentage}% partially supported" role="group" style="width:{itemStats.mitigatedPercentage}%;" class="report-header-score-mitigated"></div>
            {/if}
            {#if itemStats.unsupportedPercentage > 0}
              <div tabindex="0" title="{itemStats.unsupportedPercentage}% not supported" role="group" style="width:{itemStats.unsupportedPercentage}%;" class="report-header-score-unsupported"></div>
            {/if}
          </div>
          <div class="report-header-score-summary">
            <span class="report-header-score-summary-supported-value" title="{itemStats.supportedPercentage}% supported">{itemStats.supportedPercentage}%</span>+
            <span class="report-header-score-summary-mitigated-value" title="{itemStats.mitigatedPercentage}% partially supported">{itemStats.mitigatedPercentage}%</span> = {itemStats.fullSupportPercentage}%
          </div>
        </div>
      {/if}
    </div>
    {#if normalizedClientsList}
      {#if normalizedClientsList.unsupported.length > 0}
        <div>Unsupported clients</div>
        {#each normalizedClientsList.unsupported as client}
          <div>
            <div>{client.title}</div>
            {#if client.notes}
              {#each client.notes as noteKey}
                <div>
                  <button>{noteKey}</button>
                </div>
              {/each}
            {/if}
          </div>
        {/each}
      {/if}
      {#if normalizedClientsList.mitigated.length > 0}
        <div>Mitigated clients</div>
        {#each normalizedClientsList.mitigated as client}
          <div>
            <div>{client.title}</div>
            {#if client.notes}
              {#each client.notes as noteKey}
                <div>
                  <button>{noteKey}</button>
                </div>
              {/each}
            {/if}
          </div>
        {/each}
      {/if}
      {#if normalizedClientsList.supported.length > 0}
        <div>Supported clients</div>
        {#each normalizedClientsList.supported as client}
          <div>
            <div>{client.title}</div>
            {#if client.notes}
              {#each client.notes as noteKey}
                <div>
                  <button>{noteKey}</button>
                </div>
              {/each}
            {/if}
          </div>
        {/each}
      {/if}
    {/if}
    {#if report.rules?.notes}
      <div>Notes</div>
      {#each Object.keys(report.rules.notes).sort() as noteKey}
        <div>
          <button>{noteKey}</button>
          <div>{report.rules.notes[noteKey]}</div>
        </div>
      {/each}
    {/if}
    {#if report.rules?.url}
      <div>
        <a href="{report.rules.url}">More info</a>
      </div>
    {/if}
  </div>
</li>
