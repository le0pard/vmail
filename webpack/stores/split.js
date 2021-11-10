import {writable} from 'svelte/store'

const createBasicStore = () => {
  const initialVisibility = 'both'
  const initialVal = {visible: initialVisibility}
  const {subscribe, set, update} = writable(initialVal)
  const toggleFn = (state, isSkipInitial = false) => (currentVal) => ({
    visible: (isSkipInitial ? state : (currentVal.visible === initialVisibility ? state : initialVisibility))
  })

  return {
    subscribe,
    set,
    hideLeft: (isSkipInitial = false) => update(toggleFn('right', isSkipInitial)),
    hideRight: (isSkipInitial = false) => update(toggleFn('left', isSkipInitial)),
    hideForceRight: () => set({visible: 'left'}),
    reset: () => set(initialVal)
  }
}

export const splitState = createBasicStore()
