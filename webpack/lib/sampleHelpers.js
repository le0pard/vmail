import {memoize} from 'utils/memoize'

export const loadSampleContent = memoize(() => (
  fetch('/samples/email_sample.html').then((response) => response.text())
))
