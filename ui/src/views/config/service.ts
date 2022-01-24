import request from '@/utils/request';
import {Ref} from "vue";
import {Config, Interpreter} from './data.d';
import {Languages} from "@/utils/const";

const apiPath = 'config';

export function getInterpretersFromConfig(currConfig: any): any {
    const interpreters: any[] = []
    const languages: string[] = []
    const languageMap = {}

    Languages.forEach(item => {
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
