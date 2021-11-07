import {writable, derived} from 'svelte/store'
import {
  MULTI_LEVEL_REPORT_KEYS,
  SINGLE_LEVEL_REPORT_KEYS,
  REPORT_CSS_VARIABLES
} from 'lib/constants'

const selectLinesAndSelectors = (report) => {
  let lineToSelector = {}
  MULTI_LEVEL_REPORT_KEYS.forEach((reportInfo) => {
    if (report[reportInfo.key]) {
      Object.keys(report[reportInfo.key]).forEach((name) => {
        Object.keys(report[reportInfo.key][name]).forEach((value) => {
          report[reportInfo.key][name][value].lines.forEach((line) => {
            lineToSelector[line] ||= []
            lineToSelector[line] = [...lineToSelector[line], [reportInfo.key, name, value]]
          })
        })
      })
    }
  })
  SINGLE_LEVEL_REPORT_KEYS.forEach((reportInfo) => {
    if (report[reportInfo.key]) {
      Object.keys(report[reportInfo.key]).forEach((name) => {
        report[reportInfo.key][name].lines.forEach((line) => {
          lineToSelector[line] ||= []
          lineToSelector[line] = [...lineToSelector[line], [reportInfo.key, name]]
        })
      })
    }
  })
  if (report[REPORT_CSS_VARIABLES.key]) {
    Object.keys(report[REPORT_CSS_VARIABLES.key]).forEach((name) => {
      report[REPORT_CSS_VARIABLES.key][name].lines.forEach((line) => {
        lineToSelector[line] ||= []
        lineToSelector[line] = [...lineToSelector[line], REPORT_CSS_VARIABLES.key]
      })
    })
  }
  return lineToSelector
}

const createBasicStore = (initialVal = null) => {
  const {subscribe, set} = writable(initialVal)

  return {
    subscribe,
    set,
    reset: () => set({})
  }
}

export const reportLoading = createBasicStore(false)
export const reportError = createBasicStore(null)
export const report = createBasicStore({})
export const linesAndSelectors = derived(
  report,
  $report => selectLinesAndSelectors($report)
)
