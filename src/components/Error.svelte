<svelte:options immutable="{true}" />

<script>
  export let title
  export let message
</script>

<div class="error-container">
  <div class="error-msg">
    <div class="error-msg-body">
      <h3>{title}</h3>
      <div class="error-msg-text">{message}</div>
    </div>
  </div>
</div>

<style>
  .error-container {
    display: flex;
    margin: 4rem auto;
    align-items: center;
    justify-content: center;
  }

  .error-msg {
    border: 0;
    box-shadow: 0 10px 35px 0 rgb(56 71 109 / 8%);
    position: relative;
    display: flex;
    flex-direction: column;
    min-width: 0;
    background-color: var(--cardBgColor);
    border-radius: 0.4rem;
  }

  .error-msg-body {
    padding: 2rem;
  }

  .error-msg-text {
    color: var(--cardBaseColor);
    font-weight: 500;
  }
</style>
