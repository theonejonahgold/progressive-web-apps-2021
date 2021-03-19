/// <reference lib="webworker" />

// const CACHE_NAME = `henkerkesh-${
//   (import.meta as Record<string, any>).env.SNOWPACK_PUBLIC_SALT
// }`
const CORE_CACHE = 'henkercore'
const PAGE_CACHE = 'henkerpage'
const ASSET_CACHE = 'henkerasset'
const CORE_CACHE_URLS = ['/offline', '/styles/main.css', '/scripts/main.js']
let coreCacheName = CORE_CACHE
let pageCacheName = PAGE_CACHE
let assetCacheName = ASSET_CACHE

// Solution for type problems from: https://github.com/Microsoft/TypeScript/issues/11781#issuecomment-785350836
const sw: ServiceWorkerGlobalScope & typeof globalThis = self as any

sw.addEventListener('install', e => {
  e.waitUntil(
    (async () => {
      await updateCacheNames()
      const cache = await caches.open(coreCacheName)
      await cache.addAll(CORE_CACHE_URLS)
      await sw.skipWaiting()
    })(),
  )
})

sw.addEventListener('activate', e => {
  sw.clients.claim()
  e.waitUntil(
    (async () => {
      await updateCacheNames()
      const keys = await caches.keys()
      await Promise.all(
        keys
          .filter(
            key =>
              (key.includes(PAGE_CACHE) && key !== pageCacheName) ||
              (key.includes(CORE_CACHE) && key !== coreCacheName) ||
              (key.includes(ASSET_CACHE) && key !== assetCacheName),
          )
          .map(caches.delete),
      )
    })(),
  )
})

sw.addEventListener('fetch', async e => {
  if (e.request.method !== 'GET') e.respondWith(fetch(e.request))
  const cacheRes = await caches.match(e.request, { ignoreSearch: true })
  if (cacheRes) e.respondWith(cacheRes)
  const url = new URL(e.request.url)
  if (url.pathname.includes('/story/') || url.pathname === '/')
    e.respondWith(addToCache(pageCacheName, e.request))
  if (CORE_CACHE_URLS.includes(url.pathname))
    e.respondWith(addToCache(coreCacheName, e.request))
  e.respondWith(addToCache(assetCacheName, e.request))
})

async function updateCacheNames() {
  const res = await fetch('/version')
  const version = await res.text()
  coreCacheName = `${CORE_CACHE}-${version}`
  pageCacheName = `${PAGE_CACHE}-${version}`
  assetCacheName = `${ASSET_CACHE}-${version}`
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
