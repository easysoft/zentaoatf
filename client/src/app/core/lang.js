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

    logInfo(`load language res from ${pth}`)

    const buf = fs.readFileSync(pth)
    const obj = JSON.parse(buf.toString())
    return obj
}

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

if (DEBUG) {
    global.$lang = langHelper;
}

export default langHelper;
