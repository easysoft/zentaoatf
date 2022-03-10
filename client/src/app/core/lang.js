import {app} from "electron";

const fs = require('fs');

import {DEBUG, WORK_DIR} from '../utils/consts';

import LangHelper from '../utils/lang-helper';
import {logInfo} from "../utils/log";
import path from "path";

const langHelper = new LangHelper();

export const initLang = () => {
    let langName = app.getLocale()
    logInfo(`langName=${langName}`)

    langName = langName.toLowerCase()
    if (langName !== 'zh-cn' && langName.startsWith('zh-')) {
        langName = 'zh-cn';
    }

    loadLanguage(langName)
};

export const loadLanguage = async (langName) => {
    if (!langName) {
        return
    }

    if (langName !== langHelper.name) {
        const langData = loadLangData(langName)
        langHelper.change(langName, langData);
    }
};

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

if (DEBUG) {
    global.$lang = langHelper;
}

export default langHelper;
