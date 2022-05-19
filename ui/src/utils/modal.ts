import { $vfm } from "vue-final-modal";

// Modal.confirm({title:'dddd', content:'确定删除吗？', showOkBtn:false})
export default {
    confirm(options) {
        options.isConfirm = true;
        $vfm.show({ component: 'ZModal', bind: options }) //options {title:'title'}
    },
}