import {Controller} from '@hotwired/stimulus'
import {wrap} from 'comlink'
import {memoize} from 'utils/memoize'
import AppComponent from 'components/App'
import ErrorComponent from 'components/Error'

const getWebWorker = memoize((url) => (
  new Promise((resolve) => {
    const webWorker = new Worker(url)
    return resolve(wrap(webWorker))
  })
))

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

    getWebWorker(this.workerUrlValue).then((webWorkerObject) => {
      this.appComponent = new AppComponent({
        target: this.appContainerTarget,
        props: {
          webWorkerObject
        }
      })
    }).catch((e) => {
      this.errorComponent = new ErrorComponent({
        target: this.appContainerTarget,
        props: {
          title: 'Error to load web worker',
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
