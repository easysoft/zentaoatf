import moment from "moment";
import {AutoTestTools, ScriptLanguages, UnitTestFrameworks, UnitTestTools} from "@/utils/const";
import {Ref} from "vue";
import {Config, Interpreter} from "@/views/config/data";

export function getInterpretersFromConfig(currConfig: any): any {
    const interpreters: any[] = []
    const languages: string[] = []
    const languageMap = {}

    ScriptLanguages.forEach(item => {
        const lang = item.toLowerCase()
        languageMap[lang] = item

        if (currConfig && currConfig[lang] && currConfig[lang].trim() != '') {
            interpreters.push({ lang: lang, val: currConfig[lang] })
        } else {
            languages.push(lang)
        }
    })
    return {interpreters: interpreters, languages: languages, languageMap: languageMap}
}

export function setInterpreter(config: Ref<Config>, interpreters: Ref<Interpreter[]>): Ref<Config> {
    interpreters.value.forEach((item, i) => {
        config[item.lang] = item.val
    })
    return config
}

export function getUnitTestFrameworks(): any {
    const list = new Array<string>()
    const map = {}
    UnitTestFrameworks.forEach((item) => {
        const lowerCase = item.toLowerCase()
        list.push(lowerCase)
        map[lowerCase] = item
    })

    return {list: list, map: map}
}
export function getUnitTestTools(): any {
    const data = {}
    const map = {}

    Object.keys(UnitTestTools).forEach((key) => {
        if (! (key in data)) data[key] = []

        UnitTestTools[key].forEach(item => {
            const lowerCase = item.toLowerCase()
            data[key].push(lowerCase)
            map[lowerCase] = item
        })
    })

    return {data: data, map: map}
}

export function getAutoTestTools(): any {
    const list = new Array<string>()
    const map = {}
    AutoTestTools.forEach((item) => {
        const lowerCase = item.toLowerCase()
        list.push(lowerCase)
        map[lowerCase] = item
    })

    return {list: list, map: map}
}

const execByMap = {
    case: '选用例',
    module: '按模块',
    suite: '按套件',
    task: '按任务',
}
export function execByDef(code) {
    return execByMap[code]
}
export function momentTimeDef(tm) {
    return moment.unix(tm).format("YYYY-MM-DD HH:mm:ss")
}
export function percentDef(numb, total) {
    if (total == 0) return '0%'
    return Number(numb / total * 100).toFixed(2) + '%'
}

const osMap = {
    windows: 'Windows',
    linux: 'Linux',
    mac: 'Mac',
}
export function testEnvDef(code) {
    return osMap[code]
}
const testTypeMap = {
    func: 'Functional Testing',
    unit: 'Unit Testing',
    auto: 'Automated Testing',
}
export function testTypeDef(code) {
    return testTypeMap[code]
}
export function resultStatusDef(code) {
    if (code === true || code === 'pass')
    return '通过'

    return '失败'
}
