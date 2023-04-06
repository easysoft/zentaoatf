import 'vue-toast-notification/dist/theme-sugar.css';
import { useToast } from "vue-toastification";

const toast = useToast()
export default {
    success(options) {
        if(typeof(options) === 'string') {
            options = {message: options}
        }
        options['hideProgressBar'] = true;
        options['timeout'] = 3000;
        options['toastClassName'] = 'toast-notification toast-notification-success';
        options['bodyClassName'] = 'toast-notification-container';
        options['closeButtonClassName'] = 'toast-notification-close';
        toast.success(options.message, options);
    },
    error(options) {
        options['hideProgressBar'] = true;
        options['timeout'] = options['timeout'] ? options['timeout']: 6000;
        options['toastClassName'] = 'toast-notification toast-notification-error';
        options['bodyClassName'] = 'toast-notification-container';
        options['closeButtonClassName'] = 'toast-notification-close';
        toast.error(options.message, options);
    },
    info(options) {
        options['hideProgressBar'] = true;
        options['timeout'] = 3000;
        options['toastClassName'] = 'toast-notification toast-notification-info';
        options['bodyClassName'] = 'toast-notification-container';
        options['closeButtonClassName'] = 'toast-notification-close';
        toast.info(options.message, options);
    },
    warning(options) {
        options['hideProgressBar'] = true;
        options['timeout'] = 3000;
        options['toastClassName'] = 'toast-notification toast-notification-warning';
        options['bodyClassName'] = 'toast-notification-container';
        options['closeButtonClassName'] = 'toast-notification-close';
        toast.warning(options.message, options);
    }
}