import { Controller } from '@hotwired/stimulus'
import { Workbox, messageSW } from 'workbox-window'

const hiddenNotificationClassName = 'sw-notification__hidden'

export default class extends Controller {
  static targets = ['toast']

  initialize() {
    this.showUpdateNotification = this.showUpdateNotification.bind(this)

    this.wb = null
    this.wbRegistration = null
  }

  connect() {
    this.initServiceWorker()
  }

  disconnect() {
    this.cleanupServiceWorker()
  }

  initServiceWorker() {
    if ('serviceWorker' in navigator) {
      this.wb = new Workbox('/sw.js')

      this.wb.addEventListener('waiting', this.showUpdateNotification)
      this.wb.addEventListener('externalwaiting', this.showUpdateNotification)

      // Register the service worker after event listeners have been added.
      this.wb.register().then((r) => {
        this.wbRegistration = r
      })
    }
  }

  cleanupServiceWorker() {
    if ('serviceWorker' in navigator && this.wb) {
      this.wb.removeEventListener('waiting', this.showUpdateNotification)
      this.wb.removeEventListener('externalwaiting', this.showUpdateNotification)
    }
  }

  showUpdateNotification() {
    this.toastTarget.classList.remove(hiddenNotificationClassName)
  }

  reloadPage(e) {
    e.preventDefault()

    this.wb.addEventListener('controlling', () => window.location.reload())

    if (this.wbRegistration?.waiting) {
      // Send a message to the waiting service worker,
      // instructing it to activate.
      messageSW(this.wbRegistration.waiting, { type: 'SKIP_WAITING' })
    }
  }

  dismiss(e) {
    e.preventDefault()
    this.toastTarget.classList.add(hiddenNotificationClassName)
  }
}
