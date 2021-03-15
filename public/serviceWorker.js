const CACHE_NAME = 'henkerkesh'
const CACHE_URLS = ['/', '/offline/', '/styles/main.css', '/scripts/main.js']

self.addEventListener('install', e => {
  e.waitUntil(async () => {
    const cache = await caches.open(CACHE_NAME)
    await cache.addAll(CACHE_URLS)
  })
})

self.addEventListener('activate', e => {
  self.clients.claim()
  e.waitUntil(
    (async () => {
      const keys = await caches.keys()
      await Promise.all(
        keys.map(key => {
          key !== CACHE_NAME && caches.delete(key)
        }),
      )
    })(),
  )
})

self.addEventListener('fetch', e => {
  if (e.request.method !== 'GET') {
    return
  }
  e.respondWith(
    (async () => {
      const cacheRes = await caches.match(e.request, {
        ignoreSearch: true,
      })
      if (cacheRes) return cacheRes
      const cache = await caches.open(CACHE_NAME)
      try {
        const res = await fetch(e.request)
        cache.put(e.request, res.clone())
        return res
      } catch {
        return await cache.match('/offline/')
      }
    })(),
  )
})
