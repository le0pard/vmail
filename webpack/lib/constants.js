export const MULTI_LEVEL_REPORT_KEYS = [
  {
    key: 'html_tags'
  },
  {
    key: 'html_attributes'
  },
  {
    key: 'css_properties'
  },
  {
    key: 'at_rule_css_statements'
  }
]

export const SINGLE_LEVEL_REPORT_KEYS = [
  {
    key: 'css_selector_types'
  },
  {
    key: 'css_dimentions'
  },
  {
    key: 'css_functions'
  },
  {
    key: 'css_pseudo_selectors'
  },
  {
    key: 'img_formats'
  }
]

export const REPORT_CSS_VARIABLES = {
  key: 'css_variables'
}

export const CSS_SELECTORS_MAP = {
  0: {
    title: 'Adjacent sibling combinator',
    description: `The adjacent sibling combinator (h1 + p) allows to target an element that is directly after another.`
  },
  1: {
    title: 'Attribute selector',
    description: `The attribute selector ([attr]) targets elements with this specific attribute.`
  },
  2: {
    title: 'Chaining selectors',
    description: `Chaining selectors (.foo.bar) allows to apply styles to elements matching all the corresponding selectors.`
  },
  3: {
    title: 'Child combinator',
    description: `The child combinator is represented by a superior sign (>) between two selectors and matches the second selector if it is a direct child of the first selector.`
  },
  4: {
    title: 'Class selector',
    description: `The class selector (.className) allows to apply styles to a group of elements with the corresponding class attribute.`
  },
  5: {
    title: 'Descendant combinator',
    description: `The descendant combinator is represented by a space ( ) between two selectors and matches the second selector if it has ancestor matching the first selector.`
  },
  6: {
    title: 'General sibling combinator',
    description: `The general sibling combinator (img ~ p) allows to target any element that after another (directly or indirectly).`
  },
  7: {
    title: 'Grouping selectors',
    description: `Grouping selectors (.foo, .bar) allows to apply the same styles to the different corresponding elements.`
  },
  8: {
    title: 'ID selector',
    description: `The ID selector (#id) allows to apply styles to an element with the corresponding id attribute.`
  },
  9: {
    title: 'Type selector',
    description: `Type selector or element selectors allow to apply styles by HTML element names.`
  },
  10: {
    title: 'Universal selector *',
    description: `The universal selector (*) allows to apply styles to every elements.`
  },
}

export const EVENT_LINE_TO_EDITOR = 'line-to-editor'
export const EVENT_LINE_TO_REPORT = 'line-to-report'
