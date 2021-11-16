export const activateTheme = (theme) => {
  if (document) {
    const doc = document.querySelector(':root')
    if (doc) {
      doc.className = theme
    }
  }
}
