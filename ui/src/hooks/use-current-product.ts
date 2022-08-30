import { computed, ComputedRef } from 'vue';
import { useStore } from 'vuex';
import {ZentaoData} from '@/store/zentao';

interface ProductInfo {
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

export default function useCurrentProduct(): ComputedRef<ProductInfo> {
    const store = useStore<{ Zentao: ZentaoData }>();
    return computed<any>(() => {
        return store.state.Zentao.currProduct;
    });
}
