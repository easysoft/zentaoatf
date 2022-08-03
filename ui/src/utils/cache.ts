import {getCache, setCache} from './localCache';
import settings from '@/config/settings';
import {proxyArrToVal} from "@/utils/comm";

export const getExecBy = async () => {
    let execBy = await getCache(settings.execBy);
    if (!execBy) execBy = 'opened'
    return execBy
}
export const setExecBy = async (execBy) => {
    await setCache(settings.execBy, execBy);
}

export const getInitStatus = async () => {
    return await getCache(settings.initStatus);
}
export const setInitStatus = async () => {
    await setCache(settings.initStatus, true);
}

export const getCurrSiteId = async () => {
    return await getCache(settings.currSiteId);
}
export const setCurrSiteId = async (val) => {
    await setCache(settings.currSiteId, val);
}

export const getCurrProductIdBySite = async (currSiteId) => {
    const mp = await getCache(settings.currProductIdBySite);
    return mp ? mp[currSiteId] : 0
}
export const setCurrProductIdBySite = async (currSiteId, currProductId) => {
    if (!currSiteId) return

    let mp = await getCache(settings.currProductIdBySite);
    if (!mp) mp = {}
    mp[currSiteId + ''] = currProductId
    await setCache(settings.currProductIdBySite, mp);
}

export const getScriptDisplayBy = async (siteId, productId) => {
    console.log('getScriptDisplayBy')

    const cachedData = await getCache(settings.displayBy);
    const key = `${siteId}-${productId}`

    let displayBy = 'workspace'
    if (cachedData && cachedData[key]) {
        displayBy = cachedData[key]
    }

    return displayBy
}
export const setScriptDisplayBy = async (displayBy, siteId, productId) => {
    console.log('setScriptDisplayBy')

    let cachedData = await getCache(settings.displayBy);
    if (!cachedData) cachedData = {}

    const key = `${siteId}-${productId}`

    cachedData[key] = displayBy
    await setCache(settings.displayBy, cachedData);
}

// script filters
export const getScriptFilters = async (displayBy, siteId, productId, byDefault = '') => {
    console.log('getScriptFilters')

    const cachedData = await getCache(settings.scriptFilters);
    const key = `${displayBy}-${siteId}-${productId}`

    if (!cachedData || !cachedData[key]) {
        return {by: '', val: ''}
    }

    const mp = cachedData[key]
    const by = byDefault ? byDefault : (mp.by ? mp.by : '')
    const val = mp[by]

    return {by: by, val: val}
}
export const setScriptFilters = async (displayBy, siteId, productId, by, val) => {
    console.log('setScriptFilters')

    let cachedData = await getCache(settings.scriptFilters);
    if (!cachedData) cachedData = {}

    const key = `${displayBy}-${siteId}-${productId}`

    const mp = cachedData[key] ? cachedData[key] : {}
    mp.by = by
    mp[by] = val

    cachedData[key] = mp
    await setCache(settings.scriptFilters, cachedData);
}

export const getExpandedKeys = async (siteId, productId) => {
    console.log('getExpandedKeys')
    const key = `${siteId}-${productId}`

    const cachedData = await getCache(settings.expandedKeys);
    if (!cachedData || !cachedData[key]) {
        return []
    }

    const keys = cachedData[key] ? cachedData[key] : []

    return [...keys]
}
export const setExpandedKeys = async (siteId, productId, keys) => {
    console.log('setExpandedKeys')
    const key = `${siteId}-${productId}`

    let cachedData = await getCache(settings.expandedKeys);
    if (!cachedData) cachedData = {}

    const items = []as any[]
    keys.forEach((item) => {
        items.push(item)
    })
    cachedData[key] = items
    await setCache(settings.expandedKeys, cachedData);
}

export const getCmdHistories = async (workspaceId) => {
    const mp = await getCache(settings.cmdHistories);
    return mp && mp[workspaceId] ? mp[workspaceId] : []
}
export const setCmdHistories = async (workspaceId, items) => {
    console.log('setCmdHistories', workspaceId)

    if (!workspaceId) return

    let mp = await getCache(settings.cmdHistories);

    if (!mp) mp = {}
    mp[workspaceId] = proxyArrToVal(items)
    await setCache(settings.cmdHistories, mp);
}

export const getOpenedScripts = async () => {
    const openedScripts : string[] = await getCache(settings.openedScripts);
    return openedScripts
}
export const openScript = async (script : string) => {
    let openedScripts : string[] = await getCache(settings.openedScripts);
    if (!openedScripts) {
        openedScripts = [script]
    } else {
        openedScripts.push(script)
    }

    await setCache(settings.openedScripts, openedScripts);
}
export const closeScript = async (script : string) => {
    let newValue = [] as string[]

    const openedScripts : string[] = await getCache(settings.openedScripts);
    if (openedScripts) {
        newValue = openedScripts.filter((item) => {
            return item !== script;
        });
    }

    await setCache(settings.openedScripts, newValue);
}