import {Controller} from '@hotwired/stimulus'

export default class extends Controller {
  initialize() {
    // this.navigationMedia = window.matchMedia('(max-width: 768px)')
    // this.onNavigationMediaChange = this.onNavigationMediaChange.bind(this)
    // this.cleanupNavigationForTurboCache = this.cleanupNavigationForTurboCache.bind(this)
  }

  connect() {
    // document.addEventListener('turbo:before-cache', this.cleanupNavigationForTurboCache)
    // this.navigationMedia.addEventListener('change', this.onNavigationMediaChange)

    const go = new window.Go()
    const WASM_URL = '/parser.wasm'
    window.WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then((obj) => {
      go.run(obj.instance)
    }).then(() => window.VMail('<html></html>')).then((message) => console.log('Message', message))
  }

  disconnect() {
    // document.removeEventListener('turbo:before-cache', this.cleanupNavigationForTurboCache)
    // this.navigationMedia.removeEventListener('change', this.onNavigationMediaChange)
  }
}
