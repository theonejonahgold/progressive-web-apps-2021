/* eslint-env node */
process.env.SNOWPACK_PUBLIC_SALT = new Date().valueOf()

/** @type {import("snowpack").SnowpackUserConfig } */
module.exports = {
  mount: {
    src: { url: '/' },
    public: { url: '/', static: true },
  },
  plugins: ['@snowpack/plugin-typescript'],
  devOptions: {
    output: 'stream',
  },
  buildOptions: {
    out: 'dist',
    watch: process.env.NODE_ENV !== 'production',
    clean: false,
  },
  packageOptions: { source: 'remote', types: true },
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
