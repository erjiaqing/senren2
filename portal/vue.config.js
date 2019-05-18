module.exports = {
  baseUrl: 'oj',
  productionSourceMap: false,
  devServer: {
    proxy: {
      "/rpc": {
        target: "https://acm.whu.edu.cn:8080",
        changeOrigin: true,
      }
    }
  },
}