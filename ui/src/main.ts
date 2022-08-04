import { createApp } from 'vue';

import App from '@/App.vue';
import router from '@/config/router';
import store from '@/config/store';
import i18n from '@/config/i18n';

import _ from "lodash";
import mitt, {Emitter} from "@/utils/mitt";
import Toast, { PluginOptions } from "vue-toastification";
import "vue-toastification/dist/index.css";
import vfmPlugin from "vue-final-modal";
import ZModal from "@/components/Modal.vue";

const app = createApp(App)
app.use(store);
app.use(router)
app.use(i18n);

app.component("ZModal", ZModal);
app.use(vfmPlugin)

const options: PluginOptions = {
    // You can set your default options here
};
app.use(Toast, options);

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
