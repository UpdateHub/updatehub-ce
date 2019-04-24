var ManifestPlugin = require('webpack-manifest-plugin')

module.exports = {
  baseUrl: '/ui',
  configureWebpack: {
    plugins: [
      new ManifestPlugin()
    ]
  }
}
