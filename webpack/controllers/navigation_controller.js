import {Controller} from '@hotwired/stimulus'

export default class extends Controller {
  static targets = ['itemsList']

  initialize() {
    this.navigationMedia = window.matchMedia('(max-width: 768px)')
    this.onNavigationMediaChange = this.onNavigationMediaChange.bind(this)
    this.cleanupNavigationForTurboCache = this.cleanupNavigationForTurboCache.bind(this)
  }

  connect() {
    document.addEventListener('turbo:before-cache', this.cleanupNavigationForTurboCache)
    this.navigationMedia.addEventListener('change', this.onNavigationMediaChange)
  }

  disconnect() {
    document.removeEventListener('turbo:before-cache', this.cleanupNavigationForTurboCache)
    this.navigationMedia.removeEventListener('change', this.onNavigationMediaChange)
  }

  toggle(e) {
    e.preventDefault()
    this.changeVisibilityForNavigation(!this.isNavigationVisible())
  }

  isNavigationVisible() {
    const firstElement = this.itemsListTargets[0]
    if (firstElement) {
      return firstElement.style.display === 'block'
    }
    return false
  }

  changeVisibilityForNavigation(isShow = true) {
    const displayValue = isShow ? 'block' : 'none'
    this.itemsListTargets.forEach((el) => {
      el.style.display = displayValue
    })
  }

  onNavigationMediaChange(e) {
    this.changeVisibilityForNavigation(!e.matches)
  }

  cleanupNavigationForTurboCache() {
    if (this.navigationMedia.matches) {
      this.changeVisibilityForNavigation(false)
    }
  }
}
