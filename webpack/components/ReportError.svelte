<svelte:options immutable={true} />

<script>
  export let error
</script>

<style>
  .error-report {
    border: 0;
    box-shadow: 0 10px 35px 0 rgb(56 71 109 / 8%);
    position: relative;
    display: flex;
    flex-direction: column;
    min-width: 0;
    background-color: var(--cardBgColor);
    border-radius: 0.4rem;
  }

  .error-report-body {
    flex: 1 1 auto;
    padding: 2rem;
  }

  .error-report-title {
    color: var(--errorColor);
  }

  .error-report-text {
    color: var(--cardBaseColor);
  }
</style>

<div class="error-report">
  <div class="error-report-body">
    <h3 class="error-report-title">Error to parse HTML</h3>
    <div class="error-report-text">
      {error.toString()}
    </div>
  </div>
</div>
