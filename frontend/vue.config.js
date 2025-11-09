module.exports = {
  devServer: {
    proxy: {
      // 匹配所有以/api开头的请求
      '/api': {
        target: 'http://localhost:8080', // 后端Go Gin服务地址（需与后端一致）
        changeOrigin: true, // 允许跨域
        pathRewrite: {
          '^/api': '/api' // 路径重写（无需修改，直接转发）
        }
      }
    }
  }
}