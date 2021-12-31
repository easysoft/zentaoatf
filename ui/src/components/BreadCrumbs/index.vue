<template>
    <a-breadcrumb>
        <a-breadcrumb-item v-for="item in list" :key="item.path">
            <a-link :to="item.path">{{t(item.title)}}</a-link>
        </a-breadcrumb-item>
    </a-breadcrumb>
</template>
<script lang="ts">
import { defineComponent, PropType } from 'vue';
import { useI18n } from "vue-i18n";
import { BreadcrumbType } from '@/utils/routes';
import ALink from '@/components/ALink/index.vue';

interface BreadCrumbsSetupData {
    t: (key: string | number) => string;
}

export default defineComponent({
    name: 'BreadCrumbs',
    props: {
        list: {
            type: Array as PropType<BreadcrumbType[]>,
            default: () => {
                return [];
            }
        }
    },
    components: {
        ALink
    },
    setup(): BreadCrumbsSetupData {
        const { t } = useI18n();

        return {
            t
         }
    }
})
</script>