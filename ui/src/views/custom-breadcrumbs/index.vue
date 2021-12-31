<template>
    <div className="indexlayout-main-conent">
        <a-card :bordered="false">
            <ArrowUpOutlined :style="{fontSize: '35px', color: '#FF0000'}" />
            {{t('page.custom-breadcrumbs.msg')}}
       </a-card>
    </div>
</template>
<script lang="ts">
import { defineComponent } from "vue";
import { useI18n } from "vue-i18n";
import { ArrowUpOutlined } from '@ant-design/icons-vue';
export default defineComponent({
    components: {
        ArrowUpOutlined
    },
    setup() {
        const { t } = useI18n();

        return {
            t
        }
    }
})
</script>