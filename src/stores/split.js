import { writable } from 'svelte/store'

export const screenSizeMinMedia = () => window.matchMedia('(max-width: 800px)')

const createBasicStore = () => {
  const initialVisibility = 'both'
  const initialValue = {
    visible: initialVisibility
  }
  const { subscribe, set, update } = writable(initialValue)
  const toggleFn = (state) => () =>
    update((currentVal) => ({
      visible: screenSizeMinMedia().matches
        ? state
        : currentVal.visible === initialVisibility
        ? state
        : initialVisibility
    }))
  const toggleForMobileFn = (state) => () =>
    set({
      visible: screenSizeMinMedia().matches ? state : initialVisibility
    })

  return {
    subscribe,
    set,
    hideLeft: toggleFn('right'),
    hideRight: toggleFn('left'),
    hideForceRight: () => set({ visible: 'left' }),
    switchToRightOnMobile: toggleForMobileFn('right'),
    switchToLeftOnMobile: toggleForMobileFn('left'),
    reset: () => set(initialValue)
  }
}

export const splitState = createBasicStore()
