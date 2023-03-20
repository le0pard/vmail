import { clientsClaim } from 'workbox-core'
import { precacheAndRoute } from 'workbox-precaching'
import { registerRoute } from 'workbox-routing'
import { NetworkFirst } from 'workbox-strategies'
import { BackgroundSyncPlugin } from 'workbox-background-sync'
import { ExpirationPlugin } from 'workbox-expiration'
import { cleanupOutdatedCaches } from 'workbox-precaching/cleanupOutdatedCaches'

self.addEventListener('message', (event) => {
  if (event.data?.type === 'SKIP_WAITING') {
    self.skipWaiting()
  }
})

clientsClaim()
cleanupOutdatedCaches()

// wasm cache router
registerRoute(
  new RegExp('\\.wasm$', 'i'),
  new NetworkFirst({
    cacheName: 'wasm-modules',
    networkTimeoutSeconds: 5,
    matchOptions: {
      ignoreSearch: true
    },
    plugins: [
      new ExpirationPlugin({
        maxEntries: 10,
        maxAgeSeconds: 90 * 24 * 60 * 60, // 90 days
        matchOptions: {
          ignoreSearch: true
        },
        purgeOnQuotaError: true
      }),
      new BackgroundSyncPlugin('wasmQueue', {
        maxRetentionTime: 24 * 60 // Retry for max of 24 Hours (specified in minutes)
      })
    ]
  })
)

const cachedAssets = self.__WB_MANIFEST

precacheAndRoute(cachedAssets)
