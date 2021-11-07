import {writable, derived} from 'svelte/store'
import {
  MULTI_LEVEL_REPORT_KEYS,
  SINGLE_LEVEL_REPORT_KEYS,
  REPORT_CSS_VARIABLES,
  multiSelectorName,
  singleSelectorName,
  cssVarsSelectorName
} from 'lib/constants'

const selectLinesAndSelectors = (report) => {
  let lineToSelector = {}
  const lines = new Set()
  MULTI_LEVEL_REPORT_KEYS.forEach((reportInfo) => {
    if (report[reportInfo.key]) {
      Object.keys(report[reportInfo.key]).forEach((name) => {
        Object.keys(report[reportInfo.key][name]).forEach((value) => {
          report[reportInfo.key][name][value].lines.forEach((line) => {
            lineToSelector[line] ||= multiSelectorName(reportInfo.key, name, value)
            lines.add(line)
          })
        })
      })
    }
  })
  SINGLE_LEVEL_REPORT_KEYS.forEach((reportInfo) => {
    if (report[reportInfo.key]) {
      Object.keys(report[reportInfo.key]).forEach((name) => {
        report[reportInfo.key][name].lines.forEach((line) => {
          lineToSelector[line] ||= singleSelectorName(reportInfo.key, name)
          lines.add(line)
        })
      })
    }
  })
  if (report[REPORT_CSS_VARIABLES.key]) {
    Object.keys(report[REPORT_CSS_VARIABLES.key]).forEach((name) => {
      report[REPORT_CSS_VARIABLES.key][name].lines.forEach((line) => {
        lineToSelector[line] ||= cssVarsSelectorName()
        lines.add(line)
      })
    })
  }
  return {
    lines: [...lines],
    lineToSelector
  }
}

const createReport = () => {
  const {subscribe, set} = writable({})

  return {
    subscribe,
    update: (data) => set(data),
    reset: () => set({})
  }
}

export const report = createReport()
export const linesAndSelectors = derived(
  report,
  $report => selectLinesAndSelectors($report)
)
