<script>
  import { onMount } from 'svelte'
  import { splitState, screenSizeMinMedia } from '@stores/split'

  const onScreenSizeMinMediaChange = (e) => {
    if (e.matches) {
      splitState.hideForceRight()
    }
  }

  const handleHideLeftKey = (e) => {
    e.preventDefault()
    if (e.key === ' ' || e.code === 'Space' || e.keyCode === 32) {
      splitState.hideLeft()
    }
  }

  const handleHideLeft = (e) => {
    e.preventDefault()
    splitState.hideLeft()
  }

  const handleHideRightKey = (e) => {
    e.preventDefault()
    if (e.key === ' ' || e.code === 'Space' || e.keyCode === 32) {
      splitState.hideRight()
    }
  }

  const handleHideRight = (e) => {
    e.preventDefault()
    splitState.hideRight()
  }

  onMount(() => {
    const eventAbortController = new AbortController()
    const { signal } = eventAbortController
    // init
    onScreenSizeMinMediaChange(screenSizeMinMedia())
    // subscribe
    screenSizeMinMedia().addEventListener('change', onScreenSizeMinMediaChange, { signal })
    return () => eventAbortController?.abort()
  })
</script>

<div class="split-container">
  <div
    onclick={handleHideLeft}
    onkeypress={handleHideLeftKey}
    tabindex="0"
    role="button"
    class="split-left"
    class:split-hidden={$splitState.visible === 'right'}
  >
    <i class="arrow-left"></i>
  </div>
  <div
    onclick={handleHideRight}
    onkeypress={handleHideRightKey}
    tabindex="0"
    role="button"
    class="split-right"
    class:split-hidden={$splitState.visible === 'left'}
  >
    <i class="arrow-right"></i>
  </div>
</div>

<style>
  .split-container {
    display: flex;
    width: 30px;
  }

  .split-left,
  .split-right {
    align-items: center;
    background-color: var(--splitBgColor);
    color: var(--splitColor);
    cursor: pointer;
    display: flex;
    flex: 1;
    justify-content: center;
    transition: background-color 0.3s ease-in-out;
  }

  .split-left {
    border-left: 1px solid var(--splitBorderColor);
    border-right: 1px solid var(--splitBorderColor);
  }

  .split-right {
    border-right: 1px solid var(--splitBorderColor);
  }

  .split-left:hover,
  .split-right:hover {
    background-color: var(--splitBgHoverColor);
  }

  .split-left:active,
  .split-right:active {
    background-color: var(--splitBgHoverColor);
  }

  .arrow-left,
  .arrow-right {
    border-style: solid;
    border-color: var(--splitColor);
    border-width: 0 2px 2px 0;
    display: inline-block;
    padding: 2px;
  }

  .split-left:hover .arrow-left,
  .split-right:hover .arrow-right {
    border-color: var(--splitHoverColor);
  }

  .arrow-left {
    transform: rotate(135deg);
  }

  .arrow-right {
    transform: rotate(-45deg);
  }

  .split-hidden {
    display: none;
  }
</style>
