/* eslint-env node */
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
    sourcemap: true,
  },
  packageOptions: { source: 'remote', types: true },
  optimize: {
    bundle: true,
    splitting: true,
    treeshake: true,
    target: 'es2020',
    entrypoints: ({ files }) =>
      files.filter(file => file.includes('index.html')),
    manifest: true,
    minify: true,
  },
}
