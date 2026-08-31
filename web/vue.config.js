const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8080,
    proxy: {
      '/adb': {
        target: 'http://localhost:8088',
        changeOrigin: true
      },
      '/devices': {
        target: 'http://localhost:8088',
        changeOrigin: true
      },
      '/ip': {
        target: 'http://localhost:8088',
        changeOrigin: true
      }
    }
  }
})
