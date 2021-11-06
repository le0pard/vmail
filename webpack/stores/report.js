import {writable} from 'svelte/store'

function createReport() {
  const {subscribe, set} = writable({})

  return {
    subscribe,
    update: (data) => set(data),
    reset: () => set({})
  };
}

export const report = createReport()
