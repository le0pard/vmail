<svelte:options immutable={true} />

<script>
  export let title
  export let clients
  export let percentage
  export let bullet = 'success'
  export let notesStore
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
    margin-bottom: 0.4rem;
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

  .client-list-client {
    font-size: 0.85rem;
    color: var(--clientColor);
    margin-bottom: 0.2rem;
    margin-right: 0.4rem;
    border: 1px solid var(--splitBorderColor);
    padding: 0.1rem 0.2rem;
    border-radius: 0.4rem;
  }

  .client-list-line {
    border: 0;
    box-shadow: none;
    color: var(--mutedButtonColor);
    background-color: var(--mutedButtonBgColor);
    padding: calc(0.2rem + 1px) calc(0.3rem + 1px);
    border-radius: 0.4rem;
    margin-right: 0.2rem;
    cursor: pointer;
    user-select: none;
  }

  .client-list-line:hover,
  .client-list-line:active {
    color: var(--mutedButtonHoverColor);
    background-color: var(--mutedButtonHoverBgColor);
  }

  .client-list-line-active {
    color: var(--mutedButtonHoverColor);
    background-color: var(--mutedButtonHoverBgColor);
  }

  .client-list-percentage {
    display: flex;
    align-items: center;
    flex-shrink: 0;
    font-weight: 600;
    justify-content: center;
    padding: 0.5rem;
    min-width: 4rem;
    border-radius: 0.4rem;
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
    class:client-bullet-success="{bullet === 'success'}"></span>

  <div class="client-list-body">
    <div class="client-list-title">{title}</div>
    <div class="client-list-items">
      {#each clients as client}
        <div class="client-list-client">
          <span>{client.title}</span>
          {#if client.notes}
            {#each client.notes as noteKey}
              <button
                class="client-list-line"
                class:client-list-line-active="{$notesStore.line === noteKey}"
                on:focus="{() => notesStore.setLine(noteKey)}"
                on:mouseover="{() => notesStore.setLine(noteKey)}"
                on:blur="{() => notesStore.reset()}"
                on:mouseout="{() => notesStore.reset()}"
              >
                {noteKey}
              </button>
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
