<svelte:options immutable="{true}" />

<script>
  import {onMount, getContext} from 'svelte'
  import {createNotesStore} from 'stores/notes'
  import ClientListComponent from './ClientList'
  import NotesListComponent from './NotesList'
  import {normalizeItemName, normalizeItemVal, clientsListWithStats} from 'lib/reportHelpers'

  export let elementID
  export let reportInfo
  export let itemName
  export let itemVal
  export let report
  export let handleLineClick

  let notesStore = createNotesStore()
  let clientsWithStats = null

  const {getWebWorker} = getContext('ww')

  onMount(() => {
    const webWorker = getWebWorker()
    if (webWorker?.clientsListWithStats) {
      webWorker.clientsListWithStats(report.rules).then((resData) => {
        clientsWithStats = resData
      })
    } else {
      clientsWithStats = clientsListWithStats(report.rules)
    }
  })
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
    border-radius: 0.4rem;
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
    color: var(--baseColor);
    border-bottom: 1px dashed var(--splitBorderColor);
    padding-bottom: 0.3rem;
  }

  .report-item-header-type {
    flex: 2;
    font-weight: 600;
    font-size: 1.2rem;
  }

  .report-item-header-name {
    flex: 3;
    font-size: 1.2rem;
    font-weight: 600;
    padding: 0 1rem;
  }

  .report-item-header-link {
    flex: 1;
    text-align: right;
  }

  @media only screen and (max-width: 480px) {
    .report-item-header-type {
      flex: 1;
    }

    .report-item-header-name {
      flex: 1;
    }

    .report-item-header-link {
      display: none;
    }
  }

  .report-item-header-more-link {
    color: var(--linkColor);
    text-decoration: none;
    font-size: 0.8rem;
  }

  .report-item-header-more-link:hover,
  .report-item-header-more-link:active {
    color: var(--linkHoverColor);
  }

  .report-item-header-description {
    color: var(--cardBaseColor);
    margin: 0.2rem 0;
    width: fit-content;
    font-size: 0.8rem;
    font-weight: 500;
  }

  .report-header-main-lines {
    flex: 1;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: flex-start;
    color: var(--cardBaseColor);
    margin: 0.3rem 0;
    border-bottom: 1px dashed var(--splitBorderColor);
    padding-bottom: 0.3rem;
  }

  .report-header-main-lines-title {
    color: var(--baseColor);
    font-weight: 500;
    font-size: 1rem;
  }

  .report-line-button {
    border: 0;
    box-shadow: none;
    color: var(--headColor);
    background-color: var(--mutedButtonBgColor);
    padding: calc(0.2rem + 1px) calc(0.3rem + 1px);
    border-radius: 0.4rem;
    margin-right: 0.2rem;
    margin-bottom: 0.2rem;
    cursor: pointer;
    user-select: none;
  }

  .report-line-button:hover,
  .report-line-button:active {
    color: var(--mutedButtonHoverColor);
    background-color: var(--mutedButtonHoverBgColor);
  }
</style>

<li id="{elementID}" class="report-item">
  <div class="report-item-container">
    <div class="report-item-header">
      <div class="report-item-header-info">
        <div class="report-item-header-type">
          {reportInfo.title}
        </div>
        {#if itemName}
          <div class="report-item-header-name">
            {normalizeItemName(reportInfo.key, itemName)}
            {#if itemVal.length > 0}
              {normalizeItemVal(itemVal)}
            {/if}
          </div>
        {/if}
        {#if report.rules?.url}
          <div class="report-item-header-link">
            <a
              class="report-item-header-more-link"
              href="{report.rules.url}"
              target="_blank"
              rel="noopener noreferrer">More info</a
            >
          </div>
        {/if}
      </div>
      {#if report.rules?.description}
        <div class="report-item-header-description">
          {report.rules?.description}
        </div>
      {/if}
      <div class="report-header-main-lines">
        <div class="report-header-main-lines-title">Found on lines:</div>
        {#each report.lines as line, i}
          <button
            on:click|preventDefault="{() => handleLineClick(line)}"
            class="report-line-button"
          >
            {line}
          </button>
        {/each}
        {#if report.more_lines}<div>and more...</div>{/if}
      </div>
    </div>
    {#if clientsWithStats}
      {#if clientsWithStats.unknown.length > 0}
        <ClientListComponent
          title="Support unknown"
          bullet="unknown"
          clients="{clientsWithStats.unknown}"
          count="{clientsWithStats.unknownCount}"
          percentage="{clientsWithStats.unknownPercentage}"
          notesStore="{notesStore}"
        />
      {/if}
      {#if clientsWithStats.unsupported.length > 0}
        <ClientListComponent
          title="Unsupported clients"
          bullet="error"
          clients="{clientsWithStats.unsupported}"
          count="{clientsWithStats.unsupportedCount}"
          percentage="{clientsWithStats.unsupportedPercentage}"
          notesStore="{notesStore}"
        />
      {/if}
      {#if clientsWithStats.mitigated.length > 0}
        <ClientListComponent
          title="Partially supported clients"
          bullet="warning"
          clients="{clientsWithStats.mitigated}"
          count="{clientsWithStats.mitigatedCount}"
          percentage="{clientsWithStats.mitigatedPercentage}"
          notesStore="{notesStore}"
        />
      {/if}
      {#if clientsWithStats.supported.length > 0}
        <ClientListComponent
          title="Supported clients"
          bullet="success"
          clients="{clientsWithStats.supported}"
          count="{clientsWithStats.supportedCount}"
          percentage="{clientsWithStats.supportedPercentage}"
          notesStore="{notesStore}"
        />
      {/if}
    {/if}
    {#if report.rules?.notes && Object.keys(report.rules.notes || {}).length > 0}
      <NotesListComponent notes="{report.rules.notes}" notesStore="{notesStore}" />
    {/if}
  </div>
</li>
