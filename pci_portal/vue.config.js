module.exports = {
  baseUrl: 'pci',
  productionSourceMap: false,
  devServer: {
    proxy: {
      "/rpc": {
        target: "http://127.0.0.1:8079",
        changeOrigin: true,
      }
    }
  },
}