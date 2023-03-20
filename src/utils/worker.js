import { wrap } from 'comlink'
import { memoize } from '@utils/memoize'

export const getWebWorker = memoize(
  () => import('@utils/ww.js?worker').then(({ default: WWorker }) => {
    const webWorker = new WWorker({ name: 'Parser Worker' })
    return wrap(webWorker)
  })
)
