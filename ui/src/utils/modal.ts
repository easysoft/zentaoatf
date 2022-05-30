import { $vfm } from "vue-final-modal";

// Modal.confirm({title:'dddd', content:'确定删除吗？', showOkBtn:false})
export default {
    confirm(options, on?: { [key: string]: Function | Function[] } ) {
        options.isConfirm = true;
        options.clickToClose = false;
        $vfm.show({ component: 'ZModal', bind: options, on: on}) //options {title:'title'}
    },
}