import { computed, ComputedRef } from 'vue';
import { useElementSize, MaybeElementRef } from '@vueuse/core';

/**
 * Use width pixels value and convert to percent
 * @param ref Element ref
 * @param pixels Pixels values
 * @returns
 */
export function useWidthToPercent(ref: MaybeElementRef, ...pixels: number[]): ComputedRef<number>[] {
    const { width } = useElementSize(ref);
    return pixels.map(x => computed(() => (100 * x) / width.value));
}

/**
 * Use height pixels value and convert to percent
 * @param ref Element ref
 * @param pixels Pixels values
 * @returns
 */
export function useHeightToPercent(ref: MaybeElementRef, ...pixels: number[]): ComputedRef<number>[] {
    const { height } = useElementSize(ref);
    return pixels.map(x => computed(() => (100 * x) / height.value));
}
