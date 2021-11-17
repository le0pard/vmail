import {APP_THEMES_LIGHT, APP_THEMES_DARK} from 'lib/constants'

export const activateTheme = (theme) => {
  if (document) {
    const doc = document.querySelector(':root')
    if (doc) {
      doc.classList.remove(APP_THEMES_LIGHT, APP_THEMES_DARK)
      doc.classList.add(theme)
    }
  }
}
