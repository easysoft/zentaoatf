import { computed, ComputedRef } from 'vue';
import { useStore } from 'vuex';
import {ZentaoData} from '@/store/zentao';

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

export default function useCurrentSite(): ComputedRef<SiteInfo> {
    const store = useStore<{ Zentao: ZentaoData }>();
    return computed<any>(() => {
        const site = store.state.Zentao.currSite;
        return site.username ? site : null;
    });
}
