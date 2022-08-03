/**
 * 国际化 入口
 * @author LiQingSong
 */

import { createI18n } from "vue-i18n";
import { getLocale, setLocale, importAllLocales, defaultLang } from "@/utils/i18n";

/**
 * 框架 多语言 配置
 */
export const messages = importAllLocales();
const sysLocale = getLocale();
const i18n = createI18n({
    legacy: false,
    locale: sysLocale || defaultLang,
    messages,
});

/**
 * 设置语言
 * @param locale
 */
export function setI18nLanguage(locale: string, realReload = false): void {
    setLocale(locale, function() {
        // i18n.global.locale = locale // legacy: true
        i18n.global.locale.value = locale;
    }, realReload)
}

export default i18n;
