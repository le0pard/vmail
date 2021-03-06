export const APP_THEMES_LIGHT = 'light'
export const APP_THEMES_DARK = 'dark'

export const MULTI_LEVEL_REPORT_KEYS = [
  {
    key: 'css_properties',
    title: 'CSS Property'
  },
  {
    key: 'at_rule_css_statements',
    title: 'At-rules'
  },
  {
    key: 'html_tags',
    title: 'HTML Tag'
  },
  {
    key: 'html_attributes',
    title: 'HTML Attribute'
  }
]

export const SINGLE_LEVEL_REPORT_KEYS = [
  {
    key: 'css_selector_types',
    title: 'CSS Selector'
  },
  {
    key: 'css_dimentions',
    title: 'CSS Dimention'
  },
  {
    key: 'css_functions',
    title: 'CSS Function'
  },
  {
    key: 'css_pseudo_selectors',
    title: 'CSS Pseudo-class'
  },
  {
    key: 'link_types',
    title: 'Link'
  },
  {
    key: 'img_formats',
    title: 'Image format'
  }
]

export const REPORT_CSS_VARIABLES = {
  key: 'css_variables',
  title: 'CSS Variable'
}

export const REPORT_CSS_IMPORTANT = {
  key: 'css_important',
  title: 'CSS !important keyword'
}

export const REPORT_HTML5_DOCTYPE = {
  key: 'html5_doctype',
  title: 'HTML5 doctype'
}

export const CSS_SELECTORS_MAP = {
  0: {
    title: 'Adjacent sibling combinator'
  },
  1: {
    title: 'Attribute selector'
  },
  2: {
    title: 'Chaining selectors'
  },
  3: {
    title: 'Child combinator'
  },
  4: {
    title: 'Class selector'
  },
  5: {
    title: 'Descendant combinator'
  },
  6: {
    title: 'General sibling combinator'
  },
  7: {
    title: 'Grouping selectors'
  },
  8: {
    title: 'ID selector'
  },
  9: {
    title: 'Type selector'
  },
  10: {
    title: 'Universal selector *'
  }
}

export const EVENT_LINE_TO_EDITOR = 'vmail:line-to-editor'
export const EVENT_LINE_TO_REPORT = 'vmail:line-to-report'
export const EVENT_INLINE_CSS = 'vmail:inline-css-in-html'
export const EVENT_SUBMIT_EXAMPLE = 'vmail:submit-example-for-editor'
