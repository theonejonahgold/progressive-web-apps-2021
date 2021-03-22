/// <reference lib="webworker" />

// const CACHE_NAME = `henkerkesh-${
//   (import.meta as Record<string, any>).env.SNOWPACK_PUBLIC_SALT
// }`
const CORE_CACHE = 'henkercore'
const PAGE_CACHE = 'henkerpage'
const ASSET_CACHE = 'henkerasset'
const CORE_CACHE_URLS = ['/offline', '/styles/main.css', '/scripts/main.js']
let timestamp = ''

// Solution for type problems from: https://github.com/Microsoft/TypeScript/issues/11781#issuecomment-785350836
const sw: ServiceWorkerGlobalScope & typeof globalThis = self as any

sw.addEventListener('install', e => {
  e.waitUntil(
    (async () => {
      const cache = await caches.open(CORE_CACHE)
      await cache.addAll(CORE_CACHE_URLS)
      await sw.skipWaiting()
    })(),
  )
})

sw.addEventListener('activate', e => {
  sw.clients.claim()
  e.waitUntil(
    (async () => {
      const keys = await caches.keys()
      await Promise.all(
        keys
          .filter(
            key =>
              (key.includes(PAGE_CACHE) && key !== PAGE_CACHE) ||
              (key.includes(CORE_CACHE) && key !== CORE_CACHE) ||
              (key.includes(ASSET_CACHE) && key !== ASSET_CACHE),
          )
          .map(caches.delete),
      )
    })(),
  )
})

sw.addEventListener('fetch', e => {
  if (e.request.method !== 'GET') return e.respondWith(fetch(e.request))
  e.respondWith(
    caches.match(e.request, { ignoreSearch: true }).then(cacheRes => {
      if (cacheRes) return cacheRes
      const url = new URL(e.request.url)
      if (url.pathname.includes('/story/') || url.pathname === '/')
        return addToCache(PAGE_CACHE, e.request)
      if (CORE_CACHE_URLS.includes(url.pathname))
        return addToCache(CORE_CACHE, e.request)
      return addToCache(ASSET_CACHE, e.request)
    }),
  )
})

sw.addEventListener('message', e => {
  if ('timestamp' in e.data) timestamp = e.data.timestamp
  if ('sync' in e.data) e.waitUntil(synchronisePages)
})

sw.addEventListener('sync', e => {
  if (e.tag === 'sync-pages') e.waitUntil(synchronisePages)
})

// @ts-expect-error: Periodic sync is an event that is supported, but not typed
sw.addEventListener('periodicsync', (e: SyncEvent) => {
  if (e.tag === 'sync-pages') e.waitUntil(synchronisePages)
})

async function synchronisePages() {
  const res = await fetch('/version')
  const newTimestamp = await res.text()
  if (timestamp !== newTimestamp) {
    timestamp = newTimestamp
    sw.postMessage({ timestamp })
    const cache = await caches.open(PAGE_CACHE)
    const reqs = await cache.keys()
    await Promise.all(reqs.map(req => addToCache(PAGE_CACHE, req)))
    await sw.skipWaiting()
  }
}

async function addToCache(name: string, req: Request): Promise<Response> {
  const cache = await caches.open(name)
  try {
    const res = await fetch(req)
    cache.put(req, res.clone())
    return res
  } catch {
    return (await caches.match('/offline/')) as Response
  }
}
