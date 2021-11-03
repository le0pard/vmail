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
    background-image: url("paper.webp");
  }

  a:active {
    color: red;
    background: lightblue url("img_tree.jpg") no-repeat fixed center;
  }

  .descendant .combinator {
    margin: 0;
    text-transform: lowercase;
    background: lightblue url(img_tree.jpg) no-repeat fixed center;
  }

  h1 + p {
    margin: 0;
    text-transform: lowercase;
  }

  * {
    color: green;
  }

  a {
    color: red;
  }

  a[title] {
    color: purple;
  }

  a[href="https://example.org"] {
    color: green;
  }

  .foo.bar {
    color: green;
  }

  h1 ~ p {
    margin: 0;
    text-transform: lowercase;
  }

   h2 p {
    margin: 0;
    text-transform: lowercase;
  }

  @media (min-width: 760px) {
    .button-super {
      margin: 0 100vh;
      text-transform: lowercase;
    }

    #someUniqEl {
      display: flex;
      font-size: Initial;
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
    margin-top:var(--bs-gutter-y);
  }

  @keyframes slidein {
    from {
      transform: translateX(0%);
    }

    50%  { margin-top: 150px !important; }

    to {
      transform: translateX(100%);
    }
  }

  @media (orientation: landscape) {
    body {
      flex-direction: row;
    }
  }

  @media (orientation: portrait) {
    body {
      flex-direction: column;
    }
  }

  @font-face {
    font-family: "Open Sans";
    src: url("/fonts/OpenSans-Regular-webfont.woff2") format("woff2"),
        url("/fonts/OpenSans-Regular-webfont.woff") format("woff");
  }

  a::after {
    content: "→";
  }
</style>
<audio /><audio />
<audio />
<button type="submit">Submit</button>
<button type="reset">Reset</button>
<button type="reset11">Reset111</button>
<button type="submit">Submit 2</button>
<div style="margin: 0; text-transform: lowercase; font-size: 12px">Text</div>
<img src="img.webp" alt="img" />
<div role="application" aria-labelledby="calendar" aria-describedby="info">
    <h1 id="calendar">Calendar</h1>
    <p id="info">
        This calendar shows the game schedule for the Boston Red Sox.
    </p>
    <div role="grid">
        ...
    </div>
</div>
<img src="data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUA
    AAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO
        9TXL0Y4OHwAAAABJRU5ErkJggg==" alt="Red dot" />
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
