import {Controller} from '@hotwired/stimulus'
import {memoize} from 'utils/memoize'

const loadWasmParser = memoize(async () => {
  const go = new window.Go()
  const fetchPromise = window.fetch('/parser.wasm')
  const {instance} = await window.WebAssembly.instantiateStreaming(fetchPromise, go.importObject)
  go.run(instance)
  return instance
})

const testHTML = `<html>
<body>
<audio /><audio />
<audio />
<button type="submit">Submit</button>
<button type="reset">Reset</button>
<button type="reset11">Reset111</button>
<button type="submit">Submit 2</button>
</body>
</html>
`

export default class extends Controller {
  initialize() {
    // this.navigationMedia = window.matchMedia('(max-width: 768px)')
    // this.onNavigationMediaChange = this.onNavigationMediaChange.bind(this)
    // this.cleanupNavigationForTurboCache = this.cleanupNavigationForTurboCache.bind(this)
  }

  connect() {
    loadWasmParser().then(() => window.VMail(testHTML)).then((message) => console.log('Message', JSON.parse(message)))
    // document.addEventListener('turbo:before-cache', this.cleanupNavigationForTurboCache)
    // this.navigationMedia.addEventListener('change', this.onNavigationMediaChange)
    // window.VMail('<html></html>').then((message) => console.log('Message', message)).catch((e) => console.log('Error', e))
    // const go2 = new window.Go()
    // window.WebAssembly.instantiateStreaming(fetch(WASM_URL), go2.importObject).then((obj) => {
    //   go2.run(obj.instance)
    // }).then(() => window.VMail('<html></html>')).then((message) => console.log('Message', message))
  }

  disconnect() {
    // document.removeEventListener('turbo:before-cache', this.cleanupNavigationForTurboCache)
    // this.navigationMedia.removeEventListener('change', this.onNavigationMediaChange)
  }
}
