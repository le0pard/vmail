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
<style>
  :root{--bs-blue:#0d6efd;--bs-indigo:#6610f2;--bs-purple:#6f42c1;--bs-pink:#d63384;}

  @charset 'utf-8';

  .button {
    margin: 0;
    text-transform: lowercase;
  }

  .button:hover {
    margin: 0;
    text-transform: lowercase;
  }

  .button:hover > a {
    margin: 0;
    text-transform: lowercase;
  }

  .button:hover, .btn:hover {
    margin: 0;
    text-transform: lowercase;
  }

  a:active {
    color: red;
  }

  @media (min-width: 760px) {
    .button-super {
      margin: 0 100px;
      text-transform: lowercase;
    }

    #someUniqEl {
      display: flex;
    }
  }

  .g-4,.gx-4{--bs-gutter-x:1.5rem}.g-4,.gy-4{--bs-gutter-y:1.5rem}.g-5,.gx-5{--bs-gutter-x:3rem}.g-5,.gy-5{--bs-gutter-y:3rem}@media (min-width:576px){.col-sm{flex:1 0 0%}
  .modal-backdrop{position:fixed;top:0;left:0;z-index:1050;width:100vw;height:100vh;background-color:#000}.modal-backdrop.fade{opacity:0}
  .row>*{flex-shrink:0;width:100%;max-width:100%;padding-right:calc(var(--bs-gutter-x) * .5);padding-left:calc(var(--bs-gutter-x) * .5);margin-top:var(--bs-gutter-y)}.col{flex:1 0 0%}

  @page {
    margin: 1cm;
    padding: 2px;
  }

  @page :first {
    margin: 2cm;
    padding: 50%;
  }

  @keyframes slidein {
    from {
      transform: translateX(0%);
    }

    to {
      transform: translateX(100%);
    }
  }
</style>
<audio /><audio />
<audio />
<button type="submit">Submit</button>
<button type="reset">Reset</button>
<button type="reset11">Reset111</button>
<button type="submit">Submit 2</button>
<div style="margin: 0; text-transform: lowercase; font-size: 12px">Text</div>
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
    let startTime, endTime

    loadWasmParser().then(() => {
      startTime = performance.now()
      return window.VMail(testHTML)
    }).then((message) => {
      endTime = performance.now()
      console.log(`Call to VMail took ${endTime - startTime} milliseconds`)
      console.log(message)
    })
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
