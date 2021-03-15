if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/serviceWorker.js', {
        scope: '/',
      })
      .then(registration => {
        console.log(registration)
      })
  })
}
