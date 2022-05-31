<template>
  <div class="progress-bar" :style="style">
    <div
      v-for="(progress, index) in progressList"
      :key="index"
      :class="['progress-bar-percent', 'center', 'text-inverse', Array.isArray(barClass) ? barClass[index % barClass.length] : barClass]"
      :style="getBarStyle(index, progress)"
      :title="`${progress}%`"
    >
      <span v-if="showProgress && progress">{{progress}}</span>
    </div>
    <div v-if="label" class="progress-label">
      <span :class="labelClass" :style="labelStyle">
        {{label === true ? `${progressList[0]}%` : label}}
      </span>
    </div>
    <slot />
  </div>
</template>

<script lang="ts">
function toCssSize(size: any): string {
    if (typeof size === 'number') {
        return `${size}px`;
    }
    if (typeof size !== 'string') {
        size = `${size}`;
    }
    if (size.endsWith('px') || size.endsWith('%')) {
        return size;
    }
    size = Number.parseFloat(size);
    return Number.isNaN(size) ? '' : `${size}px`;
}
</script>

<script setup lang="ts">
import { defineProps, withDefaults, computed } from 'vue';

const props = withDefaults(defineProps<{
  barClass?: string | object | (string | object)[],
  barStyle?: object | string,
  background?: false | string,
  radius?: string | number | boolean | (string | number | boolean)[],
  colors?: string | boolean | string[],
  progress: string | number | (string | number)[],
  height: number,
  width?: string | number,
  striped?: string | number | boolean,
  label?: string | boolean,
  labelClass?: string | object | (string | object)[],
  labelStyle?: string | object,
  showProgress?: boolean,
  scaling?: string | number
}>(), {
  height: 8,
  colors: true,
  striped: false,
  background: 'var(--color-primary-pale)',
  showProgress: false,
  scaling: 1
});

const progressList = computed(() => {
    const {progress} = props;
    if(typeof progress === 'string') {
        return progress.split(',').map(x => Number.parseFloat(x));
    }
    return typeof progress === 'number' ? [progress] : progress;
});

const colorsList = computed(() => {
    const {colors} = props;
    if(colors === false) {
        return null;
    }
    let colorsList = [
        'var(--color-primary)',
        'var(--color-blue)',
        'var(--color-red)',
        'var(--color-green)',
        'var(--color-purple)',
        'var(--color-yellow)',
    ];
    if(typeof colors === 'string') {
        colorsList = colors.split(',');
    } else if(Array.isArray(colors)) {
        colorsList = colors;
    }
    return colorsList;
});

const realRadius = computed(() => {
    const {radius} = props;
    if(radius) {
        const {height} = props;
        if(Array.isArray(radius)) {
            const radiusList = radius.map((r) => toCssSize(r === '50%' ? (height / 2) : r));
            if(radiusList.length > 1 && radiusList[0] !== radiusList[1]) {
                return radiusList;
            }
            return radiusList[0];
        }
        return toCssSize(radius === '50%' ? (height / 2) : radius);
    }
    return 0;
});

const style = computed(() => {
    const {
        width, height, background
    } = props;
    const style = {
        borderRadius: '',
        width: '',
        background: '',
        height: toCssSize(height),
    };
    if(width) {
        style.width = toCssSize(width);
    }
    const realRadiusValue = realRadius.value;
    if(realRadiusValue) {
        if(Array.isArray(realRadiusValue)) {
            style.borderRadius = [realRadiusValue[0], realRadiusValue[1], realRadiusValue[1], realRadiusValue[0]].join(' ');
        } else {
            style.borderRadius = realRadiusValue;
        }
    }
    if(background) {
        style.background = background;
    }
    return style;
});

function getBarStyle(index, progress) {
    const realRadiusValue = realRadius.value;
    const colorsListValue = colorsList.value;
    const progressListValue = progressList.value;
    const {barStyle, striped, scaling} = props;

    const style = {
        ...(Array.isArray(barStyle) ? barStyle[index % barStyle.length] : barStyle),
        width: `${progress * (+scaling)}%`,
        backgroundAttachment: 'fixed'
    };

    if(colorsListValue) {
        style.backgroundColor = colorsListValue[index % colorsListValue.length];
    }

    if(striped) {
        if(typeof striped === 'string') {
            style.backgroundImage = striped;
        } else {
            style.backgroundImage = 'repeating-linear-gradient(-45deg, transparent, transparent 4px, rgba(255,255,255,.2) 4px, rgba(255,255,255,.2) 5px), linear-gradient(to right, rgba(255,255,255,.2) 0%,rgba(0,0,0,0.1) 100%)';
        }
    }

    if(realRadiusValue) {
        let radius = realRadiusValue;
        if(index === 0) {
            if(Array.isArray(realRadiusValue)) {
                ([radius] = realRadiusValue);
            }
            style.borderTopLeftRadius = radius;
            style.borderBottomLeftRadius = radius;
        }
        if(index === (progressListValue.length - 1)) {
            if(Array.isArray(realRadiusValue)) {
                ([, radius] = realRadiusValue);
            }
            style.borderTopRightRadius = radius;
            style.borderBottomRightRadius = radius;
        }
    }
    return style;
}
</script>

<style>
.progress-bar {
  display: flex;
  flex-direction: row;
  overflow: hidden;
  flex-wrap: nowrap;
  justify-items: stretch;
  overflow: visible;
}
.progress-label {
  align-self: center;
  font-size: var(--font-size-sm);
  font-family: Oswald;
  padding-left: var(--space-horz-xs);
}
</style>
