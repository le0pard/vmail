
export const normalizeItemVal = (itemVal) => {
  if (itemVal.indexOf('||') >= 0) {
    const [itemV1, itemV2] = itemVal.split('||')
    return `${itemV1}=${itemV2}`
  }
  return itemVal
}
