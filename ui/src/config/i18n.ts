/**
 * 国际化 入口
 * @author LiQingSong
 */

import { createI18n } from "vue-i18n";
import { getLocale, setLocale, importAllLocales, defaultLang } from "@/utils/i18n";

/**
 * antd 多语言 配置
 */
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import zhTW from 'ant-design-vue/es/locale/zh_TW';
import enUS from 'ant-design-vue/es/locale/en_US';
export const antdMessages: { [key: string]: any} = {
    'zh-CN': zhCN,
    'zh-TW': zhTW,
    'en-US': enUS,
}


/**
 * 框架 多语言 配置
 */
export const messages = importAllLocales();
const sysLocale = getLocale();
const i18n = createI18n({
    legacy: false,
    locale: antdMessages[sysLocale] ? sysLocale : defaultLang,
    messages,
});


/**
 * 设置语言
 * @param locale 
 */
export function setI18nLanguage(locale: string, realReload = false): void {  
    setLocale(locale,realReload, function() {
        // i18n.global.locale = locale // legacy: true
        i18n.global.locale.value = locale;        
    })
}

export default i18n;
