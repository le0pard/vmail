import {Controller} from '@hotwired/stimulus'
import {memoize} from 'utils/memoize'
import AppComponent from 'components/App'
import ErrorComponent from 'components/Error'

const loadWasmParser = memoize(async () => {
  const go = new window.Go()
  const fetchPromise = window.fetch('/parser.wasm')
  const {instance} = await window.WebAssembly.instantiateStreaming(fetchPromise, go.importObject)
  go.run(instance) // do not wait for this promise
  return instance
})

export default class extends Controller {
  static targets = ['appContainer']

  connect() {
    if (!window.WebAssembly) {
      this.errorComponent = new ErrorComponent({
        target: this.appContainerTarget,
        props: {
          message: 'Your browser do not support WebAssembly'
        }
      })
      return
    }

    loadWasmParser().then(() => {
      this.appComponent = new AppComponent({
        target: this.appContainerTarget,
        props: {}
      })
    }).catch(() => {
      this.errorComponent = new ErrorComponent({
        target: this.appContainerTarget,
        props: {
          message: 'Error to load wasm module'
        }
      })
    })
  }

  disconnect() {
    if (this.appComponent) {
      this.appComponent.$destroy()
      this.appComponent = null
    }
    if (this.errorComponent) {
      this.errorComponent.$destroy()
      this.errorComponent = null
    }
  }
}
