import {writable, derived} from 'svelte/store'
import {APP_THEMES_LIGHT, APP_THEMES_DARK} from 'lib/constants'
import {getTheme} from 'utils/theme'
import LocalStorage from 'lib/localStorage'

const createBasicStore = () => {
  const {subscribe, set} = writable(getTheme())

  return {
    subscribe,
    switchTheme: (isDarkThemeActive) => {
      const activeTheme = isDarkThemeActive ? APP_THEMES_DARK : APP_THEMES_LIGHT
      set(activeTheme)
      LocalStorage.setItem('theme', activeTheme)
    }
  }
}

export const theme = createBasicStore()
export const isDarkThemeON = derived(theme, ($theme) => $theme === APP_THEMES_DARK)
