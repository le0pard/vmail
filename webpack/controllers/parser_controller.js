import {Controller} from '@hotwired/stimulus'

import AppComponent from 'components/App'
import ErrorComponent from 'components/Error'

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

    this.appComponent = new AppComponent({
      target: this.appContainerTarget,
      props: {
        workerURL: this.workerUrlValue
      }
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
