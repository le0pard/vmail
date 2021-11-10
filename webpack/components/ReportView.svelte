<script>
	import {report} from 'stores/report'
  import ReportListComponent from './ReportList'
</script>

<style>
  .parser-report-area {
		flex-grow: 1;
		position: relative;
		overflow: scroll;
	}
</style>

<div class="parser-report-area">
  {#if Object.keys($report).length > 0}
    <ReportListComponent />
  {:else}
    <p>Submit your HTML</p>
  {/if}
</div>
