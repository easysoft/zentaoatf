import { createApp } from 'vue';

// 全局样式
import '@/assets/css/global.less';

// 引入 Antd 
import Antd from 'ant-design-vue';

// 导入 svg
import { importAllSvg } from "@/components/IconSvg/index";
importAllSvg();

import App from '@/App.vue';
import router from '@/config/routes';
import store from '@/config/store';
import i18n from '@/config/i18n';

import _ from "lodash";
import mitt, {Emitter} from "@/utils/mitt";

const app = createApp(App)
app.use(store);
app.use(router)
app.use(Antd);
app.use(i18n);
app.mount('#app');

const _emitter: Emitter = mitt();

// 全局发布
app.config.globalProperties.$pub = (...args) => {
    _emitter.emit(_.head(args), args.slice(1));
};
// 全局订阅
app.config.globalProperties.$sub = function (_event, _callback) {
    // eslint-disable-next-line prefer-rest-params
    Reflect.apply(_emitter.on, _emitter, _.toArray(arguments));
};
// 取消订阅
app.config.globalProperties.$unsub = function (_event, _callback) {
    // eslint-disable-next-line prefer-rest-params
    Reflect.apply(_emitter.off, _emitter, _.toArray(arguments));
};
