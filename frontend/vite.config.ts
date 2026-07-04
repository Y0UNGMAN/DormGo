import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import type { ConfigEnv, UserConfig } from 'vite' // 补充类型导入

// https://vitejs.dev/config/
export default defineConfig(({}: ConfigEnv): UserConfig => {
  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src')
      }
    },
    server: {
      proxy: {
        '/api': {
          target: 'http://127.0.0.1:4523/m1/7294675-7023101-default',
          changeOrigin: true,
          rewrite: path => path.replace(/^\/api/, '/api')
        }
      },
      port: 3000,
      open: true
    },
    build: {
      target: 'es2015',
      sourcemap: false,
      chunkSizeWarningLimit: 1000
    },
    define: {
      'import.meta.env': {
        VITE_API_URL: JSON.stringify(process.env.VITE_API_URL || '/api')
      }
    }
  }
})
