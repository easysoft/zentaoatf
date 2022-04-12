<template>
    <div></div>
</template>

<script lang="ts">
import {defineComponent, onMounted, onBeforeUnmount} from "vue";
import { useI18n } from "vue-i18n";

import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {notification} from "ant-design-vue";
import {ResultErr} from "@/utils/request";

export default defineComponent({
    name: 'TopNotify',

    setup() {
      const { t } = useI18n();

      const notify = (result: any) => {
        const msg = result.httpCode === 200 ? t('biz_'+result.resultCode) : t('http_'+result.httpCode)
        const desc = result.resultMsg ? result.resultMsg : ''

        notification.error({
          message: msg,
          description: desc,
        });
      }

      onMounted(() => {
        console.log('onMounted')
        bus.on(settings.eventNotify, notify);
      })
      onBeforeUnmount( () => {
        bus.off(settings.eventNotify, notify);
      })

      return {
        t,
      }
    }
})
</script>