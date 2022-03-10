import {app} from "electron";

const fs = require('fs');

import {DEBUG, WORK_DIR} from '../utils/consts';

import LangHelper from '../utils/lang-helper';
import {logInfo} from "../utils/log";
import path from "path";

/**
 * 语言访问辅助对象
 * @type {LangHelper}
 */
const langHelper = new LangHelper();

/**
 * 额外的语言表对象
 * @private
 * @type {Map<String, String>}
 */
let extraLangData = null;

/**
 * 获取平台预设的语言数据对象
 * 该语言数据默认会从 lang/ 目录下加载对应的语言文件
 * @param {String} langName 语言名称
 * @returns {Promise} 使用 Promise 异步返回处理结果
 */
const loadLangData = (langName) => {
    let pth = `lang/${langName}.json`
    if (!DEBUG) {
        pth = path.join(process.resourcesPath, pth)
    }

    logInfo(`===load language res ${pth}===`)

    const buf = fs.readFileSync(pth)
    const obj = JSON.parse(buf.toString())
    return obj
}

/**
 * 更改界面语言
 * @param {String} langName 界面语言名称
 * @param {boolean} [notifyPlatform=true] 是否通知平台变更语言
 * @returns {Promise} 使用 Promise 异步返回处理结果
 */
export const loadLanguage = async (langName) => {
    if (!langName) {
        return
    }

    if (langName !== langHelper.name) {
        const langData = loadLangData(langName)
        langHelper.change(langName, langData);
    }
};

/**
 * 初始化界面语言文本访问功能
 * @param {Map<String, String>} extraData 额外的语言表数据对象
 * @return {void}
 */
export const initLang = () => {
    let langName = app.getLocale()
    logInfo(`langName=${langName}`)

    langName = langName.toLowerCase()
    if (langName !== 'zh-cn' && langName.startsWith('zh-')) {
        langName = 'zh-cn';
    }

    loadLanguage(langName)
};

/**
 * 从多语言文本定义对象获取符合当前语言设置的文本
 *
 * @param {String|Object} obj 多语言文本定义对象
 * @param {String} defaultString 默认语言文本
 * @return {String} 语言文本
 * @example
 * const langTextObj = {
 *     'zh-cn': '你好',
 *     'en': 'Hello',
 *     'default': '你好'
 * };
 *
 * const langText = getStringFromObject(langTextObj);
 * // 当在中文环境下 `langText` 值为 `'你好'`
 * // 当在英文环境下 `langText` 值为 `'Hello'`
 * // 在其他语言环境下 `langText` 值为 `'你好'`
 */
export const getStringFromObject = (obj, defaultString) => {
    if (typeof obj === 'object') {
        let str = obj[langHelper.name];
        if (str === undefined) {
            str = obj.default;
            if (str === undefined) {
                str = obj[obj.defaultLang];
            }
            if (str === undefined) {
                str = obj[Object.keys(obj)[0]];
            }
        }
        return str === undefined ? defaultString : str;
    }
    return obj === undefined ? defaultString : obj;
};

Object.assign(langHelper, {
    getStringFromObject,
});

if (DEBUG) {
    global.$lang = langHelper;
}

export default langHelper;
