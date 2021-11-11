import {writable} from 'svelte/store'

export const screenSizeMinMedia = window.matchMedia('(max-width: 800px)')

const createBasicStore = () => {
  const initialVisibility = 'both'
  const initialValue = {
    visible: initialVisibility
  }
  const {subscribe, set, update} = writable(initialValue)
  const toggleFn = (state) => (currentVal) => ({
    visible: (
      screenSizeMinMedia.matches ? state : (
        currentVal.visible === initialVisibility ? state : initialVisibility
      )
    )
  })

  return {
    subscribe,
    set,
    hideLeft: () => update(toggleFn('right')),
    hideRight: () => update(toggleFn('left')),
    hideForceRight: () => set({visible: 'left'}),
    reset: () => set(initialValue)
  }
}

export const splitState = createBasicStore()
