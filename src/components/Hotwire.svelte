<svelte:options immutable="{true}" />

<script>
  import { onMount } from 'svelte'
  onMount(() => {
    const eventAbortController = new AbortController()
    const { signal } = eventAbortController
    // Before every page navigation, remove any previously added component hydration scripts
    document.addEventListener(
      'turbo:before-render',
      () => {
        const scripts = document.querySelectorAll('script[data-astro-component-hydration]')
        for (const script of scripts) {
          script.remove()
        }
      },
      { signal }
    )
    // After every page navigation, move the bundled styles into the body
    // document.addEventListener('turbo:render', () => {
    //   const styles = document.querySelectorAll('link[href^="/assets/asset"][href$=".css"]')
    //   for (const style of styles) {
    //     document.body.append(style)
    //   }
    // }, { signal })
    import('@hotwired/turbo')
    return () => eventAbortController?.abort()
  })
</script>
