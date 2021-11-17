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
  static values = {
    workerUrl: String
  }
  static targets = ['appContainer']

  connect() {
    if (!window.WebAssembly) {
      this.errorComponent = new ErrorComponent({
        target: this.appContainerTarget,
        props: {
          title: 'Your browser do not support WebAssembly',
          message: 'Your browser do not support WebAssembly'
        }
      })
      return
    }

    loadWasmParser().then(() => {
      if (window.VMail) {
        this.appComponent = new AppComponent({
          target: this.appContainerTarget,
          props: {
            workerURL: this.workerUrlValue,
            parserFunction: window.VMail
          }
        })
      } else {
        this.errorComponent = new ErrorComponent({
          target: this.appContainerTarget,
          props: {
            title: 'Error to load WebAssembly module',
            message: 'Error to load wasm module'
          }
        })
      }
    }).catch((e) => {
      this.errorComponent = new ErrorComponent({
        target: this.appContainerTarget,
        props: {
          title: 'Error to load WebAssembly module',
          message: e.toString()
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
