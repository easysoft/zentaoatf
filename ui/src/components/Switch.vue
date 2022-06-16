<template>
    <div class="z-switch" :class="{ 'is-checked': checked }">
        <input
            class="z-switch-input"
            ref="input"
            type="checkbox"
            :checked="checked"
            @change="handleInput"
            :true-value="trueValue"
            :false-value="falseValue"
        />
        <span class="z-switch-action"></span>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, nextTick, defineProps, defineEmits } from 'vue'
const props = defineProps({
    modelValue: {
        type: [Number, String, Boolean],
    },
    trueValue: { 
        type: [Number, String, Boolean],
        default: true,
    },
    falseValue: {
        type: [Number, String, Boolean],
        default: true,
    },
    activeColor: {
        type: [String],
        default: '#1890ff',
    },
    width:{
        type: [String],
        default: '40px',
    }
})
const emits = defineEmits(['update:modelValue', 'change'])

const input = ref(null)
const checked = computed(() => {
    return props.modelValue === props.trueValue
})
const handleInput = () => {
    nextTick(() => {
        const val = input.value.checked
        emits("update:modelValue", val);
        emits("change", val);
    })
};

</script>

<style lang='less' scoped>
.z-switch {
    position: relative;
    height: 18px;
    transition: background 0.2s;
    width: v-bind(width);
    background-image: linear-gradient(to right,rgba(0,0,0,.25),rgba(0,0,0,.25)),linear-gradient(to right,#fff,#fff);
    border-radius: 10px;
    display: inline-flex;
    align-items: center;
    vertical-align: middle;
    .z-switch-input {
        position: relative;
        z-index: 1;
        margin: 0;
        width: 100%;
        height: 100%;
        opacity: 0;
    }
    .z-switch-action {
        position: absolute;
        transition: 0.2s;
        left: 2px;
        top: 2px;
        z-index: 0;
        height: 14px;
        width: 14px;
        background: #fff;
        border-radius: 50%;
    }
    &.is-checked {
        background: v-bind(activeColor);
        .z-switch-action {
            left: 100%;
            background: #fff;
            margin-left: -18px;
        }
    }
}
</style>