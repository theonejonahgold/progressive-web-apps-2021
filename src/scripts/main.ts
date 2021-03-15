/// <reference lib="DOM" />

if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/serviceWorker.js', {
        updateViaCache: 'all',
      })
      .then(() => {
        console.log('Service worker registered')
      })
  })
}
