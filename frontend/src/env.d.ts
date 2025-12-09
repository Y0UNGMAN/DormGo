/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, unknown>
  export default component
}

// 正确声明@/utils/request的模块类型
declare module '@/utils/request' {
  import type { AxiosInstance } from 'axios'
  const request: AxiosInstance
  export default request
}

// 扩展Vite环境变量类型
interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_API_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
