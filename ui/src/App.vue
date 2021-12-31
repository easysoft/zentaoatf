<template>
    <a-config-provider :locale="antdLocales">
      <router-view></router-view>
    </a-config-provider>
</template>
<script lang="ts">
import { defineComponent, computed, onMounted } from "vue";
import { antdMessages } from "@/config/i18n";
import { setHtmlLang } from "@/utils/i18n";
import { useI18n } from "vue-i18n";
export default defineComponent({
  name: 'App',
  setup() {
    const { locale } = useI18n();
    const antdLocales = computed(()=> antdMessages[locale.value]);

    onMounted(() => {
      setHtmlLang(locale.value);
    })

    return {
      antdLocales
    }
  }
})
</script>