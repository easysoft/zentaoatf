import {createApp} from 'vue';
import VueToast from 'vue-toast-notification';
import 'vue-toast-notification/dist/theme-sugar.css';

export function notification() {
    const app = createApp({});
    app.use(VueToast);
    app.mount('#app');
    return app;
}