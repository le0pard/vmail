<svelte:options immutable="{true}" />

<script>
  import { isReportReady, reportLoading, reportError } from '@stores/report'
  import ReportListComponent from '@components/ReportList.svelte'
  import ReportLoadingComponent from '@components/ReportLoading.svelte'
  import ReportEmptyComponent from '@components/ReportEmpty.svelte'
  import ReportErrorComponent from '@components/ReportError.svelte'
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
