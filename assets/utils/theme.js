import { APP_THEMES_LIGHT, APP_THEMES_DARK } from 'lib/constants'
import LocalStorage from 'lib/localStorage'

export const getTheme = () => {
  let theme = LocalStorage.getItem('theme')

  if (!theme) {
    if (window.matchMedia('(prefers-color-scheme: dark)')?.matches) {
      theme = APP_THEMES_DARK
    }
  }

  return theme || APP_THEMES_LIGHT
}

export const activateTheme = (theme) => {
  if (document) {
    const doc = document.querySelector(':root')
    if (doc) {
      doc.classList.remove(APP_THEMES_LIGHT, APP_THEMES_DARK)
      doc.classList.add(theme)
    }
    // update <meta name="theme-color"
    const themeMeta = document.querySelector('meta[name="theme-color"]')
    const style = getComputedStyle(document.body)
    if (themeMeta && style) {
      themeMeta.setAttribute('content', style.getPropertyValue('--bgColor'))
    }
  }
}
