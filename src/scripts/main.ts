/// <reference lib="DOM" />

window.addEventListener('load', async () => {
  if ('serviceWorker' in navigator) {
    await navigator.serviceWorker.register('/serviceWorker.js', {
      updateViaCache: 'all',
    })
    const reg = await navigator.serviceWorker.ready
    if ('periodicSync' in reg) periodicallySyncPages(reg)
    else if ('sync' in reg) {
      if (navigator.onLine) return syncPages(reg)
      window.addEventListener('online', () => syncPages(reg))
    }
    if (window.location.pathname.includes('favourites')) initFavouritePage()
    else initFavouriteButtons()
  }
})

async function periodicallySyncPages(reg: ServiceWorkerRegistration) {
  const status = await navigator.permissions.query({
    // @ts-expect-error: Periodic Sync is supported in Chrome but not typed
    name: 'periodic-background-sync',
  })
  if (status.state !== 'granted') return syncPages(reg)
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
  if (!tags.includes('sync-favourites')) {
    try {
      // @ts-expect-error: Periodic Sync is supported in Chrome but not typed
      await reg.periodicSync.register('sync-favourites', {
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
    await reg.sync.register('sync-favourites')
  } catch (err) {
    console.error(err)
    if (reg.active) reg.active.postMessage({ sync: true })
  }
}

async function initFavouriteButtons() {
  const buttons: NodeListOf<HTMLButtonElement> = document.querySelectorAll(
    '[data-favourite-button]'
  )
  const cache = await caches.open('henkerfavourites')
  const keys = await cache.keys()
  const urls = keys.map(key => new URL(key.url))
  buttons.forEach(button => {
    if (
      urls.some(
        url => url.pathname === `/story/${button.dataset.id as string}/`
      )
    )
      button.textContent = 'Remove from favourites'
    button.removeAttribute('disabled')
    button.addEventListener('click', toggleFavouriteArticle)
  })
}

async function toggleFavouriteArticle(this: HTMLButtonElement) {
  const { id } = this.dataset
  const url = `/story/${id}/`
  const cache = await caches.open('henkerfavourites')
  const cacheRes = await cache.match(url)
  if (!cacheRes) {
    const res = await fetch(url)
    await cache.put(url, res)
    this.textContent = 'Remove from favourites'
  } else {
    await cache.delete(url)
    this.textContent = 'Favourite'
  }
}

async function initFavouritePage() {
  const cache = await caches.open('henkerfavourites')
  const keys = await cache.keys()
  const favouritesHtml = (await Promise.all(keys.map(parseFavourite))).join('')
  const listEl = document.querySelector('[data-favourites]') as HTMLUListElement
  listEl.innerHTML = favouritesHtml
  prepareFavourites()
}

async function parseFavourite(req: Request) {
  const res = (await caches.match(req)) as Response
  const html = await res.text()
  const template = document.createElement('template')
  template.innerHTML = html
  const favourite: HTMLTemplateElement = template.content.querySelector(
    '[data-favourite]'
  ) as HTMLTemplateElement
  return favourite.innerHTML
}

async function prepareFavourites() {
  const removeButtons = document.querySelectorAll('[data-favourite-button]')
  removeButtons.forEach(button => {
    button.removeAttribute('disabled')
    button.textContent = 'Remove from favourites'
    button.addEventListener('click', removeFromFavourites)
  })
}

async function removeFromFavourites(this: HTMLButtonElement) {
  const { id } = this.dataset
  const cache = await caches.open('henkerfavourites')
  await cache.delete(`/story/${id}/`)
  if (this.parentElement && this.parentElement.parentElement)
    this.parentElement.parentElement.remove()
}
