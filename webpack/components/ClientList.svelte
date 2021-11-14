<script>
  export let title
  export let clients
  export let percentage
  export let bullet = 'success'
</script>

<style>
  .client-list {
    display: flex;
    flex-direction: row;
    margin: 0.3rem 0;
  }

  .client-bullet {
    display: inline-block;
    width: 0.25rem;
    border-radius: 0.4rem;
    flex-shrink: 0;
    margin-right: 0.5rem;
  }

  .client-bullet-error {
    background-color: var(--errorColor);
  }

  .client-bullet-warning {
    background-color: var(--warningColor);
  }

  .client-bullet-success {
    background-color: var(--successColor);
  }

  .client-list-title {
    font-weight: 600;
    font-size: 1.1rem;
  }

  .client-list-body {
    display: flex;
    flex-direction: column;
    flex-grow: 1;
  }

  .client-list-items {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
  }

  .client-list-percentage {
    display: flex;
    align-items: center;
    flex-shrink: 0;
    font-weight: 600;
    justify-content: center;
    padding: 0.5rem;
    min-width: 4rem;
  }

  .client-list-percentage-error {
    color: var(--errorColor);
    background-color: var(--errorBgColor);
  }

  .client-list-percentage-warning {
    color: var(--warningColor);
    background-color: var(--warningBgColor);
  }

  .client-list-percentage-success {
    color: var(--successColor);
    background-color: var(--successBgColor);
  }
</style>

<div class="client-list">
  <span
    class="client-bullet"
    class:client-bullet-error="{bullet === 'error'}"
    class:client-bullet-warning="{bullet === 'warning'}"
    class:client-bullet-success="{bullet === 'success'}"
  ></span>

  <div class="client-list-body">
    <div class="client-list-title">{title}</div>
    <div class="client-list-items">
      {#each clients as client}
        <div>
          <span>{client.title}</span>
          {#if client.notes}
            {#each client.notes as noteKey}
              <button>{noteKey}</button>
            {/each}
          {/if}
        </div>
      {/each}
    </div>
  </div>

  <div
    class="client-list-percentage"
    class:client-list-percentage-error="{bullet === 'error'}"
    class:client-list-percentage-warning="{bullet === 'warning'}"
    class:client-list-percentage-success="{bullet === 'success'}"
  >
    {percentage}&#37;
  </div>
</div>



