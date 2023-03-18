<svelte:options immutable="{true}" />

<script>
  import { onMount, onDestroy } from 'svelte'
  import { EVENT_SUBMIT_EXAMPLE, EVENT_INLINE_CSS } from '@lib/constants'
  import { inlinerLoading, inlinerError } from '@stores/inliner'
  import IconComponent from '@components/Icon.svelte'

  let inlinerButtonText = 'Inline CSS in HTML'
  let sampleButtonText = 'Sample HTML/CSS'

  const screenSizeMediumMedia = window.matchMedia('(max-width: 1000px)')

  const genAndSubmitSample = () => {
    window.dispatchEvent(new window.CustomEvent(EVENT_SUBMIT_EXAMPLE, { detail: {} }))
  }

  const inlineCssInHTML = () => {
    window.dispatchEvent(new window.CustomEvent(EVENT_INLINE_CSS, { detail: {} }))
  }

  const onScreenSizeMediumMediaChange = (e) => {
    if (e.matches) {
      sampleButtonText = 'Sample'
      inlinerButtonText = 'Inline CSS'
    } else {
      sampleButtonText = 'Sample HTML/CSS'
      inlinerButtonText = 'Inline CSS in HTML'
    }
  }

  const unsubscribeInlinerError = inlinerError.subscribe((errorValue) => {
    if (errorValue) {
      setTimeout(() => inlinerError.set(null), 3000)
    }
  })

  onMount(() => {
    // init
    onScreenSizeMediumMediaChange(screenSizeMediumMedia)
    // listeners
    screenSizeMediumMedia.addEventListener('change', onScreenSizeMediumMediaChange)
    return () => screenSizeMediumMedia.removeEventListener('change', onScreenSizeMediumMediaChange)
  })

  onDestroy(unsubscribeInlinerError)
</script>

<style>
  .editor-header {
    display: flex;
    align-items: center;
    background-color: var(--headBgColor);
    flex-shrink: 0;
    overflow: hidden;
    white-space: nowrap;
    box-shadow: 0 10px 30px 0 rgb(82 63 105 / 8%);
    border-bottom: 1px solid var(--headBorderColor);
  }

  .editor-header-item {
    display: flex;
    align-items: stretch;
    padding-top: 0;
    padding-bottom: 0;
  }

  .editor-header-item-full {
    flex-grow: 1;
  }

  .editor-header-link {
    color: var(--headColor);
    background-color: transparent;
    border-bottom: 2px solid transparent;
    transition: color 0.2s ease, background-color 0.2s ease;
    padding: 0.5rem 0;
    margin: 0 1rem;
    text-decoration: none;
  }

  .editor-header-link-active {
    border-bottom: 2px solid var(--headHoverColor);
    color: var(--headHoverColor);
  }

  .editor-header-link:hover,
  .editor-header-link:active {
    border-bottom: 2px solid var(--headHoverColor);
  }

  .editor-header-logo-container {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .editor-header-logo-link {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2rem;
    height: 2rem;
    margin: 0 0.5rem;
  }

  .editor-header-inline-button {
    color: var(--buttonColor);
    display: inline-block;
    background-color: var(--buttonBgColor);
    line-height: 1.25;
    text-align: center;
    white-space: nowrap;
    vertical-align: middle;
    cursor: pointer;
    user-select: none;
    border: 1px solid transparent;
    padding: 0.3rem;
    font-size: 1rem;
    width: 100%;
  }

  .editor-header-inline-button-hidden {
    display: none;
  }

  .editor-header-inline-button-error {
    color: var(--errorColor);
    background-color: var(--errorBgColor);
    cursor: default;
    pointer-events: none;
  }

  .editor-header-inline-button:hover {
    background-color: var(--buttonBgHoverColor);
  }

  .editor-header-inline-button:active {
    background-color: var(--buttonBgActiveColor);
  }

  .editor-header-inline-button-error:hover,
  .editor-header-inline-button-error:hover {
    background-color: var(--errorBgColor);
  }

  @keyframes inline-loader {
    0%,
    80%,
    100% {
      box-shadow: 0 0;
      height: 4em;
    }
    40% {
      box-shadow: 0 -2em;
      height: 5em;
    }
  }

  .editor-header-inline-loader {
    display: none;
  }

  .editor-header-inline-loader,
  .editor-header-inline-loader:before,
  .editor-header-inline-loader:after {
    background: var(--cardBaseColor);
    -webkit-animation: inline-loader 1s infinite ease-in-out;
    animation: inline-loader 1s infinite ease-in-out;
    width: 1em;
    height: 4em;
  }

  .editor-header-inline-loader-show {
    display: block;
  }

  .editor-header-inline-loader {
    color: var(--cardBaseColor);
    text-indent: -9999em;
    margin: 0 2rem;
    position: relative;
    font-size: 0.25rem;
    transform: translateZ(0);
    animation-delay: -0.16s;
  }
  .editor-header-inline-loader:before,
  .editor-header-inline-loader:after {
    position: absolute;
    top: 0;
    content: '';
  }
  .editor-header-inline-loader:before {
    left: -0.5rem;
    animation-delay: -0.32s;
  }
  .editor-header-inline-loader:after {
    left: 0.5rem;
  }

  .editor-header-sample-button {
    color: var(--mutedButtonColor);
    border: 1px solid var(--buttonBgColor);
    background-color: transparent;
    border-radius: 0.4rem;
    padding: 0.2rem;
    margin: 0 1rem;
    cursor: pointer;
    user-select: none;
    font-size: 0.9rem;
  }

  .editor-header-sample-button:hover,
  .editor-header-sample-button:active {
    color: var(--mutedButtonHoverColor);
    border-color: var(--mutedButtonHoverBgColor);
    background-color: var(--mutedButtonHoverBgColor);
  }

  .editor-header-sample-button:active {
    box-shadow: inset 0 -10rem 0 rgb(158 158 158 / 15%);
  }
</style>

<div class="editor-header">
  <div class="editor-header-logo-container">
    <a class="editor-header-logo-link" href="/">
      <IconComponent />
    </a>
  </div>
  <div class="editor-header-item">
    <a class="editor-header-link editor-header-link-active" href="/"> VMail </a>
  </div>
  <div class="editor-header-item">
    <a class="editor-header-link" href="/faq.html">FAQ</a>
  </div>
  <div class="editor-header-item-full"></div>
  <div class="editor-header-item">
    <button
      class="editor-header-inline-button"
      class:editor-header-inline-button-hidden="{$inlinerLoading === true}"
      class:editor-header-inline-button-error="{$inlinerError !== null}"
      on:click|preventDefault="{inlineCssInHTML}"
    >
      {$inlinerError ? 'Inlining error' : inlinerButtonText}
    </button>
    <div
      class="editor-header-inline-loader"
      class:editor-header-inline-loader-show="{$inlinerLoading === true}"
    >
      Loading...
    </div>
  </div>
  <div class="editor-header-item">
    <button class="editor-header-sample-button" on:click|preventDefault="{genAndSubmitSample}">
      {sampleButtonText}
    </button>
  </div>
</div>
