<svelte:options immutable="{true}" />

<script>
  import { onMount, onDestroy } from 'svelte'
  import { theme } from '@stores/theme'
  import { activateTheme } from '@utils/theme'

  onMount(() => {
    activateTheme($theme)
  })

  const unsubscribe = theme.subscribe(v => activateTheme(v))

	onDestroy(unsubscribe)
</script>

