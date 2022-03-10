import TextMap from './text-map';
import {formatString} from './string';
import {logInfo} from "./log";

/**
 * 语言访问辅助类
 */
export default class LangHelper extends TextMap {
    /**
     * 创建一个语言访问辅助类对象
     * @param {?String} name 语言名称
     * @param {?Map<String, String>} langData 语言文本表对象
     * @memberof LangHelper
     */
    constructor(name, langData) {
        super(langData);
        this._name = name;
    }

    /**
     * 变更语言名称和语言数据
     * @param {String} name 语言名称
     * @param {Map<String, String>} langData 语言文本表对象
     * @return {void}
     * @memberof LangHelper
     */
    change(name, langData) {
        this._data = langData;
        this._name = name;
    }

    /**
     * 获取语言名称
     *
     * @readonly
     * @memberof LangHelper
     * @type {String}
     */
    get name() {
        return this._name;
    }

    /**
     * 获取错误信息对应的语言文本
     *
     * @param {string|Error} err 错误信息或错误对象本身
     * @return {string} 语言文本
     */
    error(err) {
        if (!err) {
            if (DEBUG) {
                console.collapse('LANG.error', 'redBg', '<Unknown Error>', 'redPale');
                console.error(err);
                console.groupEnd();
            }
            return '<Unknown Error>';
        }
        if (typeof err === 'string') {
            return this.string(err.startsWith('error.') ? err : `error.${err}`, err);
        }
        if (Array.isArray(err)) {
            return err.map(this.error).join(';');
        }
        let message = '';
        if (err.code) {
            message += this.string(`error.${err.code}`, `${err.message || ''}[${err.code}]`);
        } else if (err.message) {
            message = this.string(`error.${err.message}`, err.message);
        }
        if (message) {
            let formatParams = err.formats || err.extras;
            if (formatParams) {
                if (typeof formatParams === 'object' && !Array.isArray(formatParams)) {
                    message = formatString(message, formatParams);
                } else {
                    if (!Array.isArray(formatParams)) {
                        formatParams = [formatParams];
                    }
                    message = formatString(message, ...formatParams);
                }
            }
        }
        return message;
    }
}
