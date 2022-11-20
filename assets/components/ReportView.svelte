<svelte:options immutable="{true}" />

<script>
  import { isReportReady, reportLoading, reportError } from 'stores/report'
  import ReportListComponent from './ReportList'
  import ReportLoadingComponent from './ReportLoading'
  import ReportEmptyComponent from './ReportEmpty'
  import ReportErrorComponent from './ReportError'
</script>

<style>
  .parser-report-area {
    flex-grow: 1;
    position: relative;
    overflow: scroll;
    padding: 1rem;
  }
</style>

<div class="parser-report-area">
  {#if $reportError}
    <ReportErrorComponent error="{$reportError}" />
  {:else if $reportLoading}
    <ReportLoadingComponent />
  {:else if $isReportReady}
    <ReportListComponent />
  {:else}
    <ReportEmptyComponent />
  {/if}
</div>
