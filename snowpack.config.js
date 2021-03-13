/* eslint-env node */


/** @type {import("snowpack").SnowpackUserConfig } */
module.exports = {
  mount: {
    src: { url: '/' },
  },
  plugins: ['@snowpack/plugin-typescript'],
  devOptions: {
    output: 'stream',
  },
  buildOptions: {
    out: 'dist/static',
    watch: process.env.NODE_ENV !== 'production',
    clean: false
  },
  packageOptions: { source: 'remote', types: true, },
  // optimize: {
  //   bundle: true,
  //   splitting: true,
  //   treeshake: true,
  //   target: 'es2020',
  //   entrypoints: 'auto',
  //   manifest: true,
  //   minify: true,
  // },
}
