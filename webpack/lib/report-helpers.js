
const round = (num, precision = 2) => (
  Math.round((num + Number.EPSILON) * Math.pow(10, precision)) / Math.pow(10, precision)
)

export const normalizeItemVal = (itemVal) => {
  if (itemVal.indexOf('||') >= 0) {
    const [itemV1, itemV2] = itemVal.split('||')
    return `${itemV1}=${itemV2}`
  }
  return itemVal
}

export const reportStats = (rules) => {
  const countValues = Object.keys(rules.stats).reduce((agg, family) => {
    Object.keys(rules.stats[family]).forEach((platform) => {
      Object.keys(rules.stats[family][platform]).forEach((version) => {
        const state = rules.stats[family][platform][version][0]
        if (state === 'y') {
          agg = {...agg, supported: agg.supported + 1}
        } else if (state === 'n') {
          agg = {...agg, unsupported: agg.unsupported + 1}
        } else {
          agg = {...agg, mitigated: agg.mitigated + 1}
        }
      })
    })

    return agg
  }, {
    supported: 0,
    mitigated: 0,
    unsupported: 0,
  })

  const countAll = countValues.supported + countValues.mitigated + countValues.unsupported

  return {
    ...countValues,
    supportedPercentage: round(countValues.supported * 100 / countAll),
    mitigatedPercentage: round(countValues.mitigated * 100 / countAll),
    unsupportedPercentage: round(countValues.unsupported * 100 / countAll),
    fullSupportPercentage: round((countValues.supported + countValues.mitigated) * 100 / countAll),
  }
}
