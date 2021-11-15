<script>
  export let notes
  export let notesStore
</script>

<style>
  .notes-list-title {
    font-size: 1rem;
    font-weight: 600;
  }

  .notes-list-item {
    display: flex;
    align-items: flex-start;
    color: var(--cardBaseColor);
    cursor: pointer;
    margin-top: 0.2rem;
    width: fit-content;
    font-size: 0.8rem;
    font-weight: 500;
  }

  .notes-list-item:hover .notes-list-button {
    color: var(--mutedButtonHoverColor);
    background-color: var(--mutedButtonHoverBgColor);
  }

  .notes-list-button {
    display: inline-block;
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

  .notes-list-button:hover, .notes-list-button:active {
    color: var(--mutedButtonHoverColor);
    background-color: var(--mutedButtonHoverBgColor);
  }

  .notes-list-button-active {
    color: var(--mutedButtonHoverColor);
    background-color: var(--mutedButtonHoverBgColor);
  }
</style>

<div class="notes-list-title">Notes:</div>
{#each Object.keys(notes).sort() as noteKey}
  <div
    class="notes-list-item"
    on:focus={() => notesStore.setLine(noteKey)}
    on:mouseover={() => notesStore.setLine(noteKey)}
    on:blur={() => notesStore.reset()}
    on:mouseout={() => notesStore.reset()}
  >
    <button
      class="notes-list-button"
      class:notes-list-button-active="{$notesStore.line === noteKey}"
      on:focus={() => notesStore.setLine(noteKey)}
      on:mouseover={() => notesStore.setLine(noteKey)}
      on:blur={() => notesStore.reset()}
      on:mouseout={() => notesStore.reset()}
    >
      {noteKey}
    </button>
    <div>{notes[noteKey]}</div>
  </div>
{/each}
