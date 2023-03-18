<svelte:options immutable="{true}" />

<script>
  import { onMount } from 'svelte'
  import { wrap } from 'comlink'
  import { memoize } from '@utils/memoize'
  import AppComponent from '@components/App.svelte'
  import ErrorComponent from '@components/Error.svelte'

  const getWebWorker = memoize(
    () => import('@utils/ww.js?worker').then(({ default: WWorker }) => {
      const webWorker = new WWorker({ name: 'Parser Worker' })
      return wrap(webWorker)
    })
  )

  let isRendered = false

  const resetState = () => {
    isRendered = false
  }

  onMount(() => {
    isRendered = true

    const eventAbortController = new AbortController()
    const { signal } = eventAbortController

    document.addEventListener('turbo:before-cache', resetState, { signal })
    return () => eventAbortController?.abort()
  })
</script>

{#if !window.WebAssembly}
  <ErrorComponent title="Your browser do not support WebAssembly" message="Your browser do not support WebAssembly" />
{:else if !window.Worker}
  <ErrorComponent title="Your browser do not support Web Workers" message="Your browser do not support Web Workers" />
{:else}
  {#await getWebWorker()}
    <div>loading...</div>
  {:then webWorkerObject}
    {#if isRendered}
      <AppComponent webWorkerObject={webWorkerObject}>
        <slot slot="githubIcon" name="githubIcon" />
      </AppComponent>
    {/if}
  {:catch error}
     <ErrorComponent title="Error to load web worker" message={error.toString()} />
  {/await}
{/if}
