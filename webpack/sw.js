import {clientsClaim} from 'workbox-core'
import {precacheAndRoute} from 'workbox-precaching'
import {cleanupOutdatedCaches} from 'workbox-precaching/cleanupOutdatedCaches'

self.addEventListener('message', event => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting()
  }
})

clientsClaim()
cleanupOutdatedCaches()

const cachedFiles = self.__WB_MANIFEST

precacheAndRoute([
  ...cachedFiles
], {
  ignoreURLParametersMatching: [/.*/],
  cleanUrls: false
})
