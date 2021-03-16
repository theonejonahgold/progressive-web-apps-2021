/// <reference lib="webworker" />

// const CACHE_NAME = `henkerkesh-${
//   (import.meta as Record<string, any>).env.SNOWPACK_PUBLIC_SALT
// }`
const CACHE_NAME = 'henkerkesh'
const CACHE_URLS = ['/', '/offline/', '/styles/main.css', '/scripts/main.js']
let cacheName = CACHE_NAME

// Solution for type problems from: https://github.com/Microsoft/TypeScript/issues/11781#issuecomment-785350836
const sw: ServiceWorkerGlobalScope & typeof globalThis = self as any

sw.addEventListener('install', e => {
  e.waitUntil(async () => {
    const res = await fetch('/version')
    const version = await res.text()
    cacheName = `${CACHE_NAME}-${version}`
    const cache = await caches.open(cacheName)
    await cache.addAll(CACHE_URLS)
  })
})

sw.addEventListener('activate', e => {
  sw.clients.claim()
  e.waitUntil(
    (async () => {
      const res = await fetch('/version')
      const version = await res.text()
      cacheName = `${CACHE_NAME}-${version}`
      const keys = await caches.keys()
      await Promise.all(
        keys.filter(key => key !== cacheName).map(key => caches.delete(key)),
      )
    })(),
  )
})

sw.addEventListener('fetch', e => {
  if (e.request.method !== 'GET') {
    return
  }
  e.respondWith(
    (async () => {
      const cacheRes = await caches.match(e.request, {
        ignoreSearch: true,
      })
      if (cacheRes) return cacheRes
      const cache = await caches.open(cacheName)
      try {
        const res = await fetch(e.request)
        cache.put(e.request, res.clone())
        return res
      } catch {
        return (await cache.match('/offline/')) as Response
      }
    })(),
  )
})
