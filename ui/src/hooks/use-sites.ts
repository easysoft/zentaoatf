import { computed, ComputedRef, onMounted } from 'vue';
import { useStore } from 'vuex';
import { StateType } from '@/views/site/store';

interface SiteInfo {
    id: number;
    createdAt: string;
    updatedAt: string;
    deleted: boolean;
    enabled: boolean;
    isDefault: boolean;
    name: string;
    username: string;
    url: string;
}

export default function useSites(): {sites: ComputedRef<SiteInfo[]>, fetchSites: () => void} {
    const store = useStore<{ Site: StateType }>();

    const sites = computed<any[]>(() => {
        const list = store.state.Site.queryResult.result;
        if (Array.isArray(list) && list.length) {
            return list.filter(x => x.url);
        }
        return [];
    });

    function fetchSites() {
        store.dispatch("Site/list", {
            keywords: '',
            enabled: '1',
            pageSize: 99999,
            page: 1,
        });
    }

    onMounted(fetchSites);

    return {sites, fetchSites};
}
