/// <reference lib="webworker" />

const VERSION = 'v2'
const CORE_CACHE = `henkercore-${VERSION}`
const PAGE_CACHE = `henkerpage-${VERSION}`
const ASSET_CACHE = `henkerasset-${VERSION}`
const FAVOURITES_CACHE = `henkerfavourites`
const CORE_CACHE_URLS = [
  '/favourites',
  '/offline',
  '/styles/main.css',
  '/scripts/main.js',
]
let timestamp = ''

// Solution for type problems from: https://github.com/Microsoft/TypeScript/issues/11781#issuecomment-785350836
const sw: ServiceWorkerGlobalScope & typeof globalThis = self as any

sw.addEventListener('install', e => {
  e.waitUntil(
    (async () => {
      const cache = await caches.open(CORE_CACHE)
      await cache.addAll(CORE_CACHE_URLS)
      await sw.skipWaiting()
    })()
  )
})

sw.addEventListener('activate', e => {
  e.waitUntil(
    (async () => {
      await sw.clients.claim()
      const keys = await caches.keys()
      await Promise.all(
        keys
          .filter(
            key =>
              key !== PAGE_CACHE && key !== CORE_CACHE && key !== ASSET_CACHE
          )
          .map(key => caches.delete(key))
      )
      await sw.skipWaiting()
    })()
  )
  e.waitUntil(
    (async () => {
      let ts = await getTimestampFromDB()
      if (ts === '') {
        const res = await fetch('/version')
        ts = await res.text()
        saveTimestampToDB(ts)
      }
      timestamp = ts
    })()
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
    })
  )
})

sw.addEventListener('message', e => {
  if ('sync' in e.data) {
    e.waitUntil(synchronisePages(PAGE_CACHE))
    e.waitUntil(synchronisePages(FAVOURITES_CACHE))
  }
})

sw.addEventListener('sync', e => {
  if (e.tag === 'sync-pages') e.waitUntil(synchronisePages(PAGE_CACHE))
  if (e.tag === 'sync-favourites')
    e.waitUntil(synchronisePages(FAVOURITES_CACHE))
})

// @ts-expect-error: Periodic sync is an event that is supported, but not typed
sw.addEventListener('periodicsync', (e: SyncEvent) => {
  if (e.tag === 'sync-pages') e.waitUntil(synchronisePages(PAGE_CACHE))
  if (e.tag === 'sync-favourites')
    e.waitUntil(synchronisePages(FAVOURITES_CACHE))
})

async function synchronisePages(cacheName: string) {
  const res = await fetch('/version')
  const newTimestamp = await res.text()
  if (timestamp !== newTimestamp) {
    timestamp = newTimestamp
    saveTimestampToDB(timestamp)
    const cache = await caches.open(cacheName)
    const reqs = await cache.keys()
    await Promise.all(
      reqs.map(req => fetch(req).then(res => cache.put(req, res.clone())))
    )
  }
  await sw.skipWaiting()
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

function getTimestampFromDB(): Promise<string> {
  return new Promise((resolve, reject) => {
    const dbReq = indexedDB.open('henkernieuws', 3)
    dbReq.addEventListener('upgradeneeded', function (this: IDBOpenDBRequest) {
      const db = this.result
      if (!db.objectStoreNames.contains('timestamp'))
        db.createObjectStore('timestamp', {
          keyPath: 'timestamp',
          autoIncrement: false,
        })
    })
    dbReq.addEventListener('success', function (this: IDBOpenDBRequest) {
      const db = this.result
      const tx = db.transaction('timestamp', 'readwrite')
      const store = tx.objectStore('timestamp')
      const req = store.getAll()
      req.addEventListener('success', function () {
        if (!this.result.length) resolve('')
        resolve(this.result[0].timestamp)
      })
      req.addEventListener('error', reject)
      tx.addEventListener('complete', () => console.log('tx complete'))
    })
  })
}

function saveTimestampToDB(timestamp: string) {
  removeTimestampFromDB()
  const dbReq = indexedDB.open('henkernieuws', 3)
  dbReq.addEventListener('upgradeneeded', function (this: IDBOpenDBRequest) {
    const db = this.result
    if (!db.objectStoreNames.contains('timestamp'))
      db.createObjectStore('timestamp', {
        keyPath: 'timestamp',
        autoIncrement: false,
      })
  })
  dbReq.addEventListener('success', function (this: IDBOpenDBRequest) {
    const db = this.result
    const tx = db.transaction('timestamp', 'readwrite')
    const store = tx.objectStore('timestamp')
    const req = store.add({ timestamp })
    req.addEventListener('success', () =>
      console.log(`timestamp ${timestamp} saved`)
    )
    req.addEventListener('error', () =>
      console.log('something went wrong while saving')
    )
    tx.addEventListener('complete', () => console.log('tx complete'))
  })
}

function removeTimestampFromDB() {
  const dbReq = indexedDB.open('henkernieuws', 3)
  dbReq.addEventListener('upgradeneeded', function (this: IDBOpenDBRequest) {
    const db = this.result
    if (!db.objectStoreNames.contains('timestamp'))
      db.createObjectStore('timestamp', {
        keyPath: 'timestamp',
        autoIncrement: false,
      })
  })
  dbReq.addEventListener('success', function (this: IDBOpenDBRequest) {
    const db = this.result
    const tx = db.transaction('timestamp', 'readwrite')
    const store = tx.objectStore('timestamp')
    const req = store.getAllKeys()
    req.addEventListener('success', function (this: IDBRequest<IDBValidKey[]>) {
      const keys = this.result
      keys.forEach(key => store.delete(key))
    })
    req.addEventListener('error', () =>
      console.log('something went wrong while saving')
    )
    tx.addEventListener('complete', () => console.log('tx complete'))
  })
}
