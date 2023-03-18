import { memoize } from '@utils/memoize'

export const loadSampleContent = memoize(() => {
  const fetchController = new AbortController()
  const timeoutId = setTimeout(() => fetchController.abort(), 5000) // 5 seconds timeout

  return fetch('/samples/email_sample.html', {
    signal: fetchController.signal
  }).then((r) => r.text()).finally(() => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }
  })
})
