<template>
  <vue-final-modal
    :name="'ZModal'"
    :content-style="contentStyle"
    :overlay-style="overlayStyle"
    :styles="styles"
    :click-to-close="false"
    v-model="showModalRef"
    classes="modal-container"
    content-class="modal-content"
  >
    <Button class="modal-close" @click="onCancel" icon="close" size="sm" />
    <span class="modal-title">{{
      title == undefined ? t("title") : title
    }}</span>
    <div class="modal-content">
        <span v-html="content"></span>
      <slot />
    </div>
    <div class="modal-action">
      <Button
        v-if="showOkBtn"
        @click="onOk"
        class="btn-modal primary rounded"
        :label="okTitle == undefined ? t('confirm') : okTitle"
      />
      <Button
        v-if="showCancelBtn"
        @click="onCancel"
        class="btn-modal border rounded"
        :label="cancelTitle == undefined ? t('cancel') : cancelTitle"
      />
    </div>
  </vue-final-modal>
</template>

<script lang="ts">
export default {
  name: "ZModal",
};
</script>

<script setup lang="ts">
import {
  defineProps,
  defineEmits,
  withDefaults,
  defineExpose,
  ref,
  watch,
} from "vue";
import { useI18n } from "vue-i18n";
import Button from "./Button.vue";
import { $vfm } from "vue-final-modal";

export interface ZModalProps {
  showModal?: boolean;
  title: string;
  onCancel?: Function;
  onOk?: Function;
  okTitle?: string;
  cancelTitle?: string;
  content?: string;
  isConfirm?: boolean;
  showOkBtn?: boolean;
  showCancelBtn?: boolean;
  styles?: object; // Style that will be applied to the modal container.
  contentStyle?: object; // Style that will be applied to the modal content.
  overlayStyle?: object; // Style that will be applied to the modal overlay.
}
const { t } = useI18n();
const props = withDefaults(defineProps<ZModalProps>(), {
  showModal: false,
  isConfirm: false,
  showOkBtn: true,
  showCancelBtn: true,
});
const showModalRef = ref(props.showModal);
watch(props, () => {
  showModalRef.value = props.showModal;
});

const emit = defineEmits<{
  (type: "onCancel", event: {}): void;
  (type: "onOk", event: {}): void;
}>();

const confirm = (params) => {
  console.log(params);
};
const onCancel = () => {
  if (props.isConfirm) {
    $vfm.hide("ZModal");
    if (props.onCancel && typeof props.onCancel === "function") {
      props.onCancel();
    }
  }
  emit("onCancel", {});
};
const onOk = () => {
  if (props.isConfirm) {
    $vfm.hide("ZModal");
    if (props.onOk && typeof props.onOk === "function") {
      props.onOk();
    }
  }
  emit("onOk", {});
};
defineExpose({
  confirm,
});
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
  padding: 1rem;
  border-radius: 0.25rem;
  background: #fff;
  justify-content: space-between;
  max-height: 90%;
  overflow: auto;
}

.modal-title {
  margin: 0 2rem 0 0;
  justify-content: space-between;
  font-size: 1.1rem;
  font-weight: 700;
}

.modal-close {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  border-width: 0;
  background-color: transparent;
  font-size: 1rem;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
}

.modal-action {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-shrink: 0;
  padding: var(--space-base) 0;
}

.btn-modal {
  margin-left: 10px;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  min-width: 70px;
}
</style>
