<script>
  import {getContext} from 'svelte'
  import ClientListComponent from './ClientList'
  import {normalizeItemVal, clientsListWithStats} from 'lib/report-helpers'

  export let reportInfo
  export let itemName
  export let itemVal
  export let report
  export let handleLineClick

  let clientsWithStats = null

  const {getWebWorker} = getContext('ww')
  const webWorker = getWebWorker()

  if (webWorker?.clientsListWithStats) {
    webWorker.clientsListWithStats(report.rules).then((resData) => {
      clientsWithStats = resData
    })
  } else {
    clientsWithStats = clientsListWithStats(report.rules)
  }
</script>

<style>
  .report-item {
    border: 0;
    box-shadow: 0 10px 35px 0 rgb(56 71 109 / 8%);
    position: relative;
    display: flex;
    flex-direction: column;
    min-width: 0;
    background-color: var(--cardBgColor);
    margin-bottom: 1rem;
  }

  .report-item:last-child {
    margin-bottom: 0;
  }

  .report-item-container {
    display: flex;
    flex-direction: column;
    padding: 1rem;
  }

  .report-item-header {
    display: flex;
    flex-direction: column;
  }

  .report-item-header-info {
    flex: 1;
    display: flex;
    flex-direction: row;
  }

  .report-item-header-type {
    flex: 2;
  }

  .report-item-header-name {
    flex: 3;
  }

  .report-item-header-link {
    flex: 1;
    text-align: right;
  }

  .report-header-main-lines {
    flex: 1;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: flex-start;
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
  <div class="report-item-container">
    <div class="report-item-header">
      <div class="report-item-header-info">
        <div class="report-item-header-type">{reportInfo.title}</div>
        <div class="report-item-header-name">
          {itemName}
          {#if itemVal.length > 0}
            {normalizeItemVal(itemVal)}
          {/if}
        </div>
        {#if report.rules?.url}
          <div class="report-item-header-link">
            <a href="{report.rules.url}" target="_blank" rel="noopener noreferrer">More info</a>
          </div>
        {/if}
      </div>
      <div class="report-header-main-lines">
        {#each report.lines as line, i}
          <button on:click|preventDefault={() => handleLineClick(line)} class="report-line-button">{line}</button>
          {#if i < report.lines.length - 1},{/if}
        {/each}
        {#if report.more_lines}
          <div>and more...</div>
        {/if}
      </div>
    </div>
    {#if clientsWithStats}
      {#if clientsWithStats.unsupported.length > 0}
        <ClientListComponent
          title="Unsupported clients"
          bullet="error"
          clients={clientsWithStats.unsupported}
          percentage={clientsWithStats.unsupportedPercentage}
        />
      {/if}
      {#if clientsWithStats.mitigated.length > 0}
        <ClientListComponent
          title="Mitigated clients"
          bullet="warning"
          clients={clientsWithStats.mitigated}
          percentage={clientsWithStats.mitigatedPercentage}
        />
      {/if}
      {#if clientsWithStats.supported.length > 0}
        <ClientListComponent
          title="Supported clients"
          bullet="success"
          clients={clientsWithStats.supported}
          percentage={clientsWithStats.supportedPercentage}
        />
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
  </div>
</li>


<!-- {#if false}
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
{/if} -->
