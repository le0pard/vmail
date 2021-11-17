import {CSS_SELECTORS_MAP} from './constants'

const FAMILY_MAP = {
  'gmail': 'Gmail',
  'outlook': 'Outlook',
  'yahoo': 'Yahoo! Mail',
  'apple-mail': 'Apple Mail',
  'aol': 'AOL',
  'thunderbird': 'Mozilla Thunderbird',
  'microsoft': 'Microsoft',
  'samsung-email': 'Samsung Email',
  'sfr': 'SFR',
  'orange': 'Orange',
  'protonmail': 'ProtonMail',
  'hey': 'HEY',
  'mail-ru': 'Mail.ru',
  'fastmail': 'Fastmail',
  'laposte': 'LaPoste.net',
  't-online-de': 'T-online.de',
  'free-fr': 'Free.fr'
}

const PLATFORM_MAP = {
  'desktop-app': 'Desktop',
  'desktop-webmail': 'Desktop Webmail',
  'mobile-webmail': 'Mobile Webmail',
  'webmail': 'Webmail',
  'ios': 'iOS',
  'android': 'Android',
  'windows': 'Windows',
  'macos': 'macOS',
  'windows-mail': 'Windows Mail',
  'outlook-com': 'Outlook.com'
}

const getFamily = (family) => FAMILY_MAP[family] ?? family
const getPlatform = (platform) => PLATFORM_MAP[platform] ?? platform

const round = (num, precision = 2) => (
  Math.round((num + Number.EPSILON) * Math.pow(10, precision)) / Math.pow(10, precision)
)

export const camelize = (str) => (
  str.replace(/([-_][a-z])/ig, ($1) => (
    $1.toUpperCase()
      .replace('-', '')
      .replace('_', '')
  ))
)

export const normalizeItemName = (reportKey, itemName) => {
  switch (reportKey) {
    case 'css_selector_types':
      return CSS_SELECTORS_MAP[itemName]?.title ?? ''
    default:
      return itemName
  }
}

export const normalizeItemVal = (itemVal) => {
  if (itemVal.indexOf('||') >= 0) {
    const [itemV1, itemV2] = itemVal.split('||')
    return `[${itemV1}=${itemV2}]`
  }
  return itemVal
}

export const getTooltipText = (matches) => (
  matches.map(([reportInfo, itemName, itemVal]) => (
    [
      reportInfo.title.replace(/\s/g, '\u00a0'), // "\u00a0" is non break space
      itemName && `:\u00a0${normalizeItemName(reportInfo.key, itemName)}`, // "\u00a0" is non break space
      itemVal && `(${normalizeItemVal(itemVal)})`
    ].filter(Boolean).join('')
  )).join(', ')
)

export const clientsListWithStats = (rules) => {
  const reducedData = Object.keys(rules.stats).reduce((agg, family) => {
    Object.keys(rules.stats[family]).forEach((platform) => {
      const versionsCount = Object.keys(rules.stats[family][platform]).length
      Object.keys(rules.stats[family][platform]).forEach((version) => {
        const [state, ...notes] = rules.stats[family][platform][version]
        const clientData = {
          title: `${getFamily(family)} ${getPlatform(platform)}${versionsCount > 1 ? `(${version})` : ''}`,
          notes
        }
        if (state === 'y') {
          agg = {
            ...agg,
            supported: [
              ...agg.supported,
              clientData
            ],
            supportedCount: agg.supportedCount + 1
          }
        } else if (state === 'n') {
          agg = {
            ...agg,
            unsupported: [
              ...agg.unsupported,
              clientData
            ],
            unsupportedCount: agg.unsupportedCount + 1
          }
        } else {
          agg = {
            ...agg,
            mitigated: [
              ...agg.mitigated,
              clientData
            ],
            mitigatedCount: agg.mitigatedCount + 1
          }
        }
      })
    })

    return agg
  }, {
    supported: [],
    supportedCount: 0,
    mitigated: [],
    mitigatedCount: 0,
    unsupported: [],
    unsupportedCount: 0
  })

  const countAll = reducedData.supportedCount + reducedData.mitigatedCount + reducedData.unsupportedCount

  return {
    ...reducedData,
    supportedPercentage: round(reducedData.supportedCount * 100 / countAll),
    mitigatedPercentage: round(reducedData.mitigatedCount * 100 / countAll),
    unsupportedPercentage: round(reducedData.unsupportedCount * 100 / countAll)
  }
}
