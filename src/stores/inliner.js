import { writable } from 'svelte/store'

const createBasicStore = (initialVal = null) => {
  const { subscribe, set } = writable(initialVal)

  return {
    subscribe,
    set,
    reset: () => set(initialVal)
  }
}

export const inlinerLoading = createBasicStore(false)
export const inlinerError = createBasicStore(null)
