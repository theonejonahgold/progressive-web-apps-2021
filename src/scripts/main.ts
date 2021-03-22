/// <reference lib="DOM" />

window.addEventListener('load', async () => {
  if ('serviceWorker' in navigator) {
    await navigator.serviceWorker.register('/serviceWorker.js', {
      updateViaCache: 'all',
    })
    const reg = await navigator.serviceWorker.ready
    navigator.serviceWorker.addEventListener('message', e => {
      if ('timestamp' in e.data)
        localStorage.setItem('build-timestamp', e.data.timestamp)
    })
    const timestamp = localStorage.getItem('build-timestamp')
    if (timestamp) reg.active?.postMessage({ timestamp })
    if ('periodicSync' in reg) periodicallySyncPages(reg)
    else if ('sync' in reg) {
      syncPages(reg)
      if ('sync' in reg && !navigator.onLine)
        window.addEventListener('online', () => syncPages(reg))
    }
  }
  prepareDB()
  initFavouriteButtons()
  if (window.location.pathname.includes('favourites')) initFavouritePage()
})

async function periodicallySyncPages(reg: ServiceWorkerRegistration) {
  const status = await navigator.permissions.query({
    // @ts-expect-error: Periodic Sync is supported in Chrome but not typed
    name: 'periodic-background-sync',
  })
  if (status.state !== 'granted') return syncPages(reg)
  // @ts-expect-error: Periodic Sync is supported in Chrome but not typed
  await reg.periodicSync.unregister('sync-pages')
  // @ts-expect-error: Periodic Sync is supported in Chrome but not typed
  const tags = await reg.periodicSync.getTags()
  if (!tags.includes('sync-pages')) {
    try {
      // @ts-expect-error: Periodic Sync is supported in Chrome but not typed
      await reg.periodicSync.register('sync-pages', {
        minInterval: 60 * 60 * 1000,
      })
    } catch (err) {
      console.error(err)
    }
  }
}

async function syncPages(reg: ServiceWorkerRegistration) {
  try {
    await reg.sync.register('sync-pages')
  } catch (err) {
    console.error(err)
    reg.active?.postMessage({ sync: true })
  }
}

function prepareDB() {
  const dbReq = indexedDB.open('henkernieuws', 2)
  dbReq.addEventListener('upgradeneeded', function (this: IDBOpenDBRequest) {
    const db = this.result
    if (!db.objectStoreNames.contains('favourites'))
      db.createObjectStore('favourites', {
        keyPath: 'id',
        autoIncrement: false,
      })
  })
  dbReq.addEventListener('success', () =>
    console.log('DB successfully prepared')
  )
}

function initFavouriteButtons() {
  const buttons: NodeListOf<HTMLButtonElement> = document.querySelectorAll(
    '[data-favourite]'
  )
  const dbReq = indexedDB.open('henkernieuws', 2)
  dbReq.addEventListener('success', function (this: IDBOpenDBRequest) {
    const db = this.result
    const tx = db.transaction('favourites', 'readonly')
    const store = tx.objectStore('favourites')
    const req = store.getAll()
    req.addEventListener('success', function (this: IDBRequest<any[]>) {
      buttons.forEach(button => {
        if (this.result.some(val => val.id === button.dataset.id))
          button.textContent = 'Remove from favourites'
        button.removeAttribute('disabled')
        button.addEventListener('click', toggleFavouriteArticle)
      })
    })
  })
}

function toggleFavouriteArticle(this: HTMLButtonElement) {
  const id = this.dataset.id
  const title = this.dataset.title
  const url = this.dataset.url
  const author = this.dataset.author
  const dbReq = indexedDB.open('henkernieuws', 2)
  dbReq.addEventListener('success', function (this: IDBOpenDBRequest) {
    const db = this.result
    const tx = db.transaction('favourites', 'readonly')
    const store = tx.objectStore('favourites')
    const req = store.get(id as string)
    req.addEventListener('success', function (this: IDBRequest) {
      const tx = db.transaction('favourites', 'readwrite')
      const store = tx.objectStore('favourites')
      let req: IDBRequest
      if (this.result != undefined) req = store.delete(id as string)
      else req = store.add({ id, title, url, author })
      tx.addEventListener('complete', () => console.log('completed tx'))
      req.addEventListener('success', () => console.log('success'))
      req.addEventListener('error', () => console.log('error'))
    })
  })
  if (this.textContent == 'Favourite')
    this.textContent = 'Remove from favourites'
  else this.textContent = 'Favourite'
}

function initFavouritePage() {
  const dbReq = indexedDB.open('henkernieuws', 2)
  dbReq.addEventListener('success', function (this: IDBOpenDBRequest) {
    const db = this.result
    const tx = db.transaction('favourites', 'readonly')
    const store = tx.objectStore('favourites')
    const req = store.getAll()
    req.addEventListener('success', function (this: IDBRequest<any[]>) {
      const listEl: HTMLUListElement = document.querySelector(
        '[data-favourites]'
      ) as HTMLUListElement
      const content = this.result
        .map(
          fav =>
            `<li><a href="${fav.url}"><h4>${fav.title}</h4></a><p>bij ${fav.author}</p><button data-remove="${fav.id}">Remove from favourites</button></li>`
        )
        .join('')
      listEl.innerHTML = content
      const removeButtons: NodeListOf<HTMLButtonElement> = document.querySelectorAll(
        '[data-remove]'
      )
      removeButtons.forEach(button => {
        button.addEventListener('click', function () {
          const id = this.dataset.remove
          const tx = db.transaction('favourites', 'readwrite')
          const store = tx.objectStore('favourites')
          const req = store.delete(id as string)
          req.addEventListener('success', () => button.parentElement?.remove())
        })
      })
    })
  })
}
