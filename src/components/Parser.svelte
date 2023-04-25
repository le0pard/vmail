<svelte:options immutable="{true}" />

<script>
  import { onMount } from 'svelte'
  import { getWebWorker } from '@utils/worker'
  import AppComponent from '@components/App.svelte'
  import ErrorComponent from '@components/Error.svelte'

  let isPageRendered = false // trigger destroy for nested components, if turbo change page

  onMount(() => {
    isPageRendered = true

    const eventAbortController = new AbortController()
    const { signal } = eventAbortController

    document.addEventListener(
      'turbo:before-cache',
      () => {
        isPageRendered = false
        eventAbortController?.abort()
      },
      { signal, once: true }
    )
    return () => eventAbortController?.abort()
  })
</script>

{#if !window.WebAssembly}
  <ErrorComponent
    title="Your browser do not support WebAssembly"
    message="Your browser do not support WebAssembly"
  />
{:else if !window.Worker}
  <ErrorComponent
    title="Your browser do not support Web Workers"
    message="Your browser do not support Web Workers"
  />
{:else}
  {#await getWebWorker()}
    <div>loading...</div>
  {:then webWorkerObject}
    {#if isPageRendered}
      <AppComponent webWorkerObject="{webWorkerObject}">
        <slot slot="githubIcon" name="githubIcon" />
      </AppComponent>
    {/if}
  {:catch error}
    <ErrorComponent title="Error to load web worker" message="{error.toString()}" />
  {/await}
{/if}
