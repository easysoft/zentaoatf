<template>
  <div>
    <a-dropdown>
        <span class="dropDown">
          {{languageLabels[locale]}}
        </span>
        <template #overlay>
            <a-menu class="menu">
                <template v-for="item in locales" :key="item">
                  <template v-if="item !== locale">
                    <a-menu-item  @click="changeLang(item)">
                      <span role="img" :aria-label="languageLabels[item]">
                          {{languageIcons[item]}}
                      </span>
                      {{languageLabels[item]}}
                    </a-menu-item>
                  </template>
                </template>
            </a-menu>
        </template>
    </a-dropdown>
  </div>
</template>
<script lang="ts">
import { defineComponent, WritableComputedRef } from "vue";
import { setI18nLanguage } from "@/config/i18n";
import { useI18n } from "vue-i18n";
import IconSvg from "@/components/IconSvg";
interface SelectLangSetupData {
    locales: string[];
    languageLabels: {[key: string]: string};
    languageIcons: {[key: string]: string};
    changeLang: (key) => void;
    locale: WritableComputedRef<string>;
}

export default defineComponent({
    name: 'SelectLang',
    components: {
    },
    setup(): SelectLangSetupData {

        const { locale }  = useI18n();
       
        const locales: string[] = ['zh-CN', 'en-US'];
        const languageLabels: {[key: string]: string} = {
            'zh-CN': '简体中文',
            'en-US': 'English',
        };
        const languageIcons: {[key: string]: string} = {
            'zh-CN': '🇨🇳',
            'en-US': '🇺🇸',
        };

        // 切换语言
        const changeLang = (key): void => {
          console.log(key)
          setI18nLanguage(key);
        }

        return {
            locales,
            languageLabels,
            languageIcons,
            changeLang,
            locale
        }
    }
})
</script>
<style lang="less" scoped>
.menu {
  .anticon {
    margin-right: 8px;
  }
  .ant-dropdown-menu-item {
    min-width: 160px;
  }
}
.dropDown {
  display: inline-block;
  width: 90px;
  cursor: pointer;
  font-size: 14px;
}
</style>