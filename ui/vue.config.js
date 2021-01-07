var ManifestPlugin = require('webpack-manifest-plugin')

process.env.VUE_APP_VERSION = require('./updatehub-ce.json').version

module.exports = {
  publicPath: '/ui',
  configureWebpack: {
    plugins: [
      new ManifestPlugin()
    ]
  }
}
