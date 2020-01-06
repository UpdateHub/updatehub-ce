var ManifestPlugin = require('webpack-manifest-plugin')

module.exports = {
  publicPath: '/ui',
  configureWebpack: {
    plugins: [
      new ManifestPlugin()
    ]
  }
}
