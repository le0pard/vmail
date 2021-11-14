<script>
	import {report} from 'stores/report'
  import ReportListComponent from './ReportList'
  import ReportEmptyComponent from './ReportEmpty'
</script>

<style>
  .parser-report-area {
		flex-grow: 1;
		position: relative;
		overflow: scroll;
    padding: 1rem 1rem;
	}
</style>

<div class="parser-report-area">
  {#if Object.keys($report).length > 0}
    <ReportListComponent />
  {:else}
    <ReportEmptyComponent />
  {/if}
</div>
