import {writable, derived} from 'svelte/store'
import {APP_THEMES_LIGHT, APP_THEMES_DARK} from 'lib/constants'
import LocalStorage from 'lib/localStorage'

const getTheme = () => {
  let theme = LocalStorage.getItem('theme')

  if (!theme) {
    if (window.matchMedia('(prefers-color-scheme: dark)')?.matches) {
      theme = APP_THEMES_DARK
    }
  }

  return theme || APP_THEMES_LIGHT
}

const createBasicStore = () => {
  const {subscribe, set} = writable(getTheme())

  return {
    subscribe,
    switchTheme: (isDarkThemeActive) => {
      const activeTheme = isDarkThemeActive ? APP_THEMES_DARK : APP_THEMES_LIGHT
      set(activeTheme)
      LocalStorage.setItem(
        'theme',
        activeTheme
      )
    }
  }
}

export const theme = createBasicStore()
export const isDarkThemeON = derived(
  theme,
  $theme => $theme === APP_THEMES_DARK
)
