// general css, css-only components, third-party libraries
import './css/app.sass'

// general polifils
import 'focus-visible'

// wasm
import 'vendors/wasm_exec'

if (!window.WebAssembly.instantiateStreaming) { // polyfill
  window.WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer()
    return await WebAssembly.instantiate(source, importObject)
  }
}
