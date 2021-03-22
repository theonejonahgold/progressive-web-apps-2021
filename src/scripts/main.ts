/// <reference lib="DOM" />

if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/serviceWorker.js', { updateViaCache: 'all' })
      .then(() => navigator.serviceWorker.ready)
      .then(reg => {
        navigator.serviceWorker.addEventListener('message', e => {
          if ('timestamp' in e.data)
            localStorage.setItem('build-timestamp', e.data.timestamp)
        })
        const timestamp = localStorage.getItem('build-timestamp')
        if (timestamp) reg.active?.postMessage({ timestamp })
        if ('periodicSync' in reg)
          //@ts-expect-error: Periodic sync is supported in some browsers, but not typed
          return reg.periodicSync.getTags().then((tags: string[]) => {
            if (!tags.includes('sync-pages'))
              //@ts-expect-error: Periodic sync is supported in some browsers, but not typed
              reg.periodicSync.register('sync-pages', {
                minInterval: 24 * 60 * 60 * 1000,
              })
          })
        else if ('sync' in reg) return reg.sync.register('sync-pages')
        // @ts-expect-error: This line is a fallback in case a browser doesn't support (periodic) background syncing
        reg.active.postMessage({ sync: true })
      })
      .catch(console.error)
  })
}
