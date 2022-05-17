<template>
    <vue-final-modal v-model="showModal" classes="modal-container" content-class="modal-content">
        <Button class="modal__close" @click="onCancel" icon="close" size="sm" />
        <span class="modal__title">{{title}}</span>
        <div class="modal__content">
            <slot />
        </div>
        <div class="modal__action">
            <Button @click="onOk" class="btn-modal" :label="t('confirm')" />
            <Button @click="onCancel" class="btn-modal" :label="t('cancel')" />
        </div>
    </vue-final-modal>
</template>

<script lang="ts">
export default {
    name: 'ZModal',
}
</script>

<script setup lang="ts">
import { defineProps, defineEmits, withDefaults, defineExpose, ref, computed } from 'vue';
import Icon from './Icon.vue';
import Button from './Button.vue';
import { $vfm, VueFinalModal } from 'vue-final-modal'
import { useI18n } from "vue-i18n";

export interface ZModalProps {
    showModal: boolean,
    title: string
}
const { t } = useI18n();
const props = withDefaults(defineProps<ZModalProps>(), {
    showModal: true,
    title: '标题',
});
console.log(props)

const showModal = computed(() => {
    return props.showModal;
  });

const emit = defineEmits<{
    (type: 'onCancel', event: {}): void,
    (type: 'onOk', event: {}): void,
}>();

const confirm = (params) => {
    console.log(params)
}
const onCancel = () => {
    emit('onCancel', {});
}
const onOk = () => {
    emit("onOk", {});
}
defineExpose({
    confirm
})
</script>

<style scoped>
:deep(.modal-container) {
    display: flex;
    justify-content: center;
    align-items: center;
}

:deep(.modal-content) {
    position: relative;
    display: flex;
    flex-direction: column;
    margin: 0 1rem;
    padding: 1rem;
    border: 1px solid #e2e8f0;
    border-radius: 0.25rem;
    background: #fff;
    min-width: 300px;
    min-height: 200px;
    justify-content: space-between;
}

.modal__title {
    margin: 0 2rem 0 0;
    font-size: 1.2rem;
    font-weight: 700;
}

.modal__close {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    border: 0 solid #e2e8f0;
    margin-left: 10px;
    background-color: transparent;
    font-size: 1rem;
    cursor: pointer;
    padding: .25rem .5rem;
}

.modal__action {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-shrink: 0;
    padding: 1rem 0 0;
}

.dark-mode div::v-deep .modal-content {
    border-color: #2d3748;
    background-color: #1a202c;
}

.btn-modal {
    border: 0 solid #e2e8f0;
    margin-left: 10px;
    text-transform: none;
    background-color: transparent;
    font-size: 100%;
    cursor: pointer;
    padding: .25rem .5rem;
    border-width: 1px;
    border-radius: .25rem;
}
</style>
