import {writable} from 'svelte/store'

export const createNotesStore = () => {
  const initialValue = {line: null}
  const {subscribe, set} = writable(initialValue)

  return {
    subscribe,
    setLine: (line) => set({line}),
    reset: () => set(initialValue)
  }
}
