<svelte:options immutable="{true}" />

<script>
  import {onMount} from 'svelte'
  import {splitState, screenSizeMinMedia} from 'stores/split'

  const onScreenSizeMinMediaChange = (e) => {
    if (e.matches) {
      splitState.hideForceRight()
    }
  }

  const handleHideLeft = () => {
    splitState.hideLeft()
  }

  const handleHideRight = () => {
    splitState.hideRight()
  }

  onMount(() => {
    onScreenSizeMinMediaChange(screenSizeMinMedia)
    screenSizeMinMedia.addEventListener('change', onScreenSizeMinMediaChange)
    return () => screenSizeMinMedia.removeEventListener('change', onScreenSizeMinMediaChange)
  })
</script>

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

<div class="split-container">
  <div
    on:click|preventDefault="{handleHideLeft}"
    class="split-left"
    class:split-hidden="{$splitState.visible === 'right'}"
  >
    <i class="arrow-left"></i>
  </div>
  <div
    on:click|preventDefault="{handleHideRight}"
    class="split-right"
    class:split-hidden="{$splitState.visible === 'left'}"
  >
    <i class="arrow-right"></i>
  </div>
</div>
