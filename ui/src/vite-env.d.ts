/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_APP_TITLE: string
    // more env variables...

    readonly VUE_APP_PORT: string; // 8000

    // mock 是否开启 true|false ， development环境有效
    readonly VUE_APP_MOCK: string; //  true

    // api接口域名
    readonly VUE_APP_APIHOST_MOCK: string; //  /api
    readonly VUE_APP_APISUFFIX: string; //  api/v1
    readonly VUE_APP_APIHOST: string; //  http://127.0.0.1:8085/api/v1
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}
