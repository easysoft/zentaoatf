import 'vue-toast-notification/dist/theme-sugar.css';
import { useToast } from "vue-toastification";

const toast = useToast()
export default {
    success(options) {
        console.log(options)
        toast.success(options.message, options);
    },
    error(options) {
        toast.error(options.message, options);
    },
    info(options) {
        toast.info(options.message, options);
    },
    warning(options) {
        toast.warning(options.message, options);
    }
}