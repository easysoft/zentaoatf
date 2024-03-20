/**
 * 国际化 utils
 * @author LiQingSong
 */
import { LocaleMessages } from '@intlify/core-base';
import { VueMessageType } from "vue-i18n";

// window.localStorage 存储key
export const localeKey = 'locale';

// 默认语言
export const defaultLang = 'zh-CN';

/**
 * 验证语言命名规则 zh-CN
 * @returns boolen
 * @author LiQingSong
 */
export const localeNameExp = (lang: string): boolean => {
    const localeExp = /^([a-z]{2})-?([A-Z]{2})?$/
    return localeExp.test(lang);
}

/**
 * 设置 html lang 属性值
 * @param lang 语言的 key
 * @author LiQingSong
 */
export const setHtmlLang = (lang: string): void => {
    /**
     * axios.defaults.headers.common['Accept-Language'] = locale
     */
    document.querySelector('html')?.setAttribute('lang', lang);
}

/**
 * 获取当前选择的语言
 * 获取的浏览器语言默认项目中有可能不支持，所以在config/i18n.ts中要加以判断
 * @returns string
 * @author LiQingSong
 */
export const getLocale = (): string => {
    const lang = typeof window.localStorage !== 'undefined' ? window.localStorage.getItem(localeKey) : '';
    const isNavigatorLanguageValid = typeof navigator !== 'undefined' && typeof navigator.language === 'string';
    const browserLang = isNavigatorLanguageValid ? navigator.language.split('-').join('-') : '';
    return lang || browserLang || defaultLang;
};

/**
 * 切换语言
 * @param lang 语言的 key
 * @param realReload 是否刷新页面，默认刷新
 * @author LiQingSong
 */
export const setLocale = (lang: string, callback: () => void, realReload = true): void => {

  if (lang !== undefined && !localeNameExp(lang)) {
    // for reset when lang === undefined
    throw new Error('setLocale lang format error');
  }
  if (getLocale() !== lang) {
    if (typeof window.localStorage !== 'undefined') {
      window.localStorage.setItem(localeKey, lang || '');
    }

    if (realReload) {
        window.location.reload();
    } else {
        setHtmlLang(lang);

        if(typeof callback === 'function') {
            callback();
        }
    }

  }
};

/**
 * 自动导入 框架自定义语言
 * @author LiQingSong
 */
export function importAllLocales(): LocaleMessages<VueMessageType> {
    const modules: any = {};
    try {
        const viewsRequireContext = import.meta.glob('../views/**/locales/*.ts', {eager: true});
        for (const path in viewsRequireContext) {
            if (viewsRequireContext.hasOwnProperty(path)) {
                const modulesConent = viewsRequireContext[path] as any;
                if (modulesConent.default) {
                    const modulesName = path.match(/\/([^/]+)\/locales/)?.[1];
                    if (modulesName) {
                        if (modules[modulesName]) {
                            modules[modulesName] = {
                                ...modules[modulesName],
                                ...modulesConent.default
                            };
                        } else {
                            modules[modulesName] = modulesConent.default;
                        }
                    }
                }
            }
        }

        const layoutsRequireContext = import.meta.glob('../layouts/**/locales/*.ts', {eager: true});
        for (const path in layoutsRequireContext) {
            if (layoutsRequireContext.hasOwnProperty(path)) {
                const modulesConent = layoutsRequireContext[path] as any;
                if (modulesConent.default) {
                    const modulesName = path.match(/\/([^/]+)\/locales/)?.[1];
                    if (modulesName) {
                        if (modules[modulesName]) {
                            modules[modulesName] = {
                                ...modules[modulesName],
                                ...modulesConent.default
                            };
                        } else {
                            modules[modulesName] = modulesConent.default;
                        }
                    }
                }
            }
        }

        const componentsRequireContext = import.meta.glob('../components/**/locales/*.ts', {eager: true});
        for (const path in componentsRequireContext) {
            if (componentsRequireContext.hasOwnProperty(path)) {
                const modulesConent = componentsRequireContext[path] as any;
                if (modulesConent.default) {
                    const modulesName = path.match(/\/([^/]+)\/locales/)?.[1];
                    if (modulesName) {
                        if (modules[modulesName]) {
                            modules[modulesName] = {
                                ...modules[modulesName],
                                ...modulesConent.default
                            };
                        } else {
                            modules[modulesName] = modulesConent.default;
                        }
                    }
                }
            }
        }

        const localesRequireContext = import.meta.glob('../locales/*.ts', {eager: true});
        for (const path in localesRequireContext) {
            if (localesRequireContext.hasOwnProperty(path)) {
                const modulesConent = localesRequireContext[path] as any;
                if (modulesConent.default) {
                    const modulesName = path.match(/\/([^/]+)\.ts$/)?.[1];
                    if (modulesName) {
                        if (modules[modulesName]) {
                            modules[modulesName] = {
                                ...modules[modulesName],
                                ...modulesConent.default
                            };
                        } else {
                            modules[modulesName] = modulesConent.default;
                        }
                    }
                }
            }
        }
    } catch (error) {
        console.log(error);
    }

    return modules;
}
