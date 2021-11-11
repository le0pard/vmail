<script>
  export let reportInfo
  export let itemName
  export let report
  export let handleLineClick
</script>

<style>
  .report-item {
    background-color: var(--headColor);
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

  .report-header-score-chart>div:hover::before {
    content: '';
    display: block;
    position: absolute;
    left: -2px;
    right: -2px;
    top: -2px;
    bottom: -2px;
    border: 2px solid #1b1b1d;
    outline: 1px solid #fff;
    z-index: 1;
  }

  .report-header-score-chart>div:hover::after {
    content: attr(title);
    position: absolute;
    z-index: 2;
    top: 1.25rem;
    bottom: auto;
    right: 0.5rem;
    padding: 0.5rem;
    min-width: 6rem;
    max-width: calc(100% - 2rem);
    font-size: .75rem;
    line-height: 1.25;
    color: rgba(255,255,255,0.85);
    background: #2a2a2e;
    border: 1px solid #3c3c3d;
    border-radius: 0.25rem;
    box-shadow: -2px -2px 4px rgb(0 0 0 / 50%);
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

</style>

<li class="report-item">
  <div class="report-container">
    <div class="report-header">
      <div class="report-header-main">
        <div class="report-header-main-top">
          <div class="report-header-type">{reportInfo.title}</div>
          <div class="report-header-name">
            {itemName}
          </div>
        </div>
        <div class="report-header-main-lines">
          {#each report.lines as line}
            <button on:click|preventDefault={() => handleLineClick(line)}>{line}</button>
          {/each}
        </div>
      </div>
      <div class="report-header-score">
        <div class="report-header-score-chart">
          <div tabindex="0" title="57.58% supported" role="group" style="width:57.58%;" class="report-header-score-supported"></div>
          <div tabindex="0" title="3.03% partially supported" role="group" style="width:3.03%;" class="report-header-score-mitigated"></div>
          <div tabindex="0" title="39.39% not supported" role="group" style="width:39.39%;" class="report-header-score-unsupported"></div>
        </div>
        <div class="report-header-score-summary">
			    <span class="report-header-score-summary-supported-value" title="81.82% supported">81.82%</span>+
          <span class="report-header-score-summary-mitigated-value" title="12.12% partially supported">12.12%</span> = 93.94%
        </div>
      </div>
    </div>
    <div>{report}</div>
  </div>
</li>
