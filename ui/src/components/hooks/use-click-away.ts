import {onUnmounted, onMounted, Ref} from 'vue';

export default function useClickAway(eleRef: HTMLElement | Ref<HTMLElement | undefined>, handleEvent: (event: Event) => void, events = 'click'): void {
    const handler = (event: Event) => {
        let ele: HTMLElement | undefined;
        if (eleRef instanceof HTMLElement) {
            ele = eleRef;
        } else if (eleRef.value) {
            ele = eleRef.value;
        }
        if (ele && event?.target && ele?.contains && !ele.contains(event.target as Node)) {
            handleEvent(event);
        }
    };

    onMounted(() => {
        document.addEventListener(events, handler, false);
    });
    onUnmounted(() => {
        document.removeEventListener(events, handler, false);
    });
}
