import { $vfm } from "vue-final-modal";

// Modal.confirm({title:'dddd', content:'确定删除吗？', showOkBtn:false})
export default {
    confirm(options :{ [key: string]: any }, on?: { [key: string]: Function | Function[] } ) :void {
        options.isConfirm = true;
        options.clickToClose = false;
        $vfm.show({ component: 'ZModal', bind: options, on: on == undefined ? {} : on}) //options {title:'title'}
    },
}