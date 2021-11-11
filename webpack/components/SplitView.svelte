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

  .split-left, .split-right {
    align-items: center;
    background-color: var(--splitColor);
    color: #fff;
    cursor: pointer;
    display: flex;
    flex: 1;
    justify-content: center;
    transition: background-color .3s ease-in-out;
  }

  .split-left {
    border-right: 1px solid var(--splitBorderColor);
  }

  .split-left:hover, .split-right:hover {
    background-color: var(--splitColorHover);
  }

  .arrow-left, .arrow-right {
    border: solid black;
    border-width: 0 2px 2px 0;
    display: inline-block;
    padding: 2px;
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
  <div on:click|preventDefault={handleHideLeft} class="split-left" class:split-hidden="{$splitState.visible === 'right'}">
    <i class="arrow-left"></i>
  </div>
  <div on:click|preventDefault={handleHideRight} class="split-right" class:split-hidden="{$splitState.visible === 'left'}">
    <i class="arrow-right"></i>
  </div>
</div>
