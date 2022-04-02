import {getCache, setCache} from './localCache';
import settings from '@/config/settings';

export const getInitStatus = async () => {
    const initStatus = await getCache(settings.initStatus);
    return initStatus
}
export const setInitStatus = async () => {
    await setCache(settings.initStatus, true);
}

export const getCurrSiteId = async () => {
    const currSiteId = await getCache(settings.currSiteId);
    return currSiteId
}
export const setCurrSiteId = async (val) => {
    await setCache(settings.currSiteId, val);
}

export const getCurrProductIdBySite = async (currSiteId) => {
    const mp = await getCache(settings.currProductIdBySite);
    const currProductId = mp ? mp[currSiteId] : 0
    return currProductId
}
export const setCurrProductIdBySite = async (currSiteId, currProductId) => {
    if (!currSiteId) return

    let mp = await getCache(settings.currProductIdBySite);
    if (!mp) mp = {}
    mp[currSiteId + ''] = currProductId
    await setCache(settings.currProductIdBySite, mp);
}

// script filters
export const getScriptFilters = async (siteId, productId) => {
    console.log('getScriptFilters')

    const cachedData = await getCache(settings.scriptFilters);
    const key = `${siteId}-${productId}`

    if (!cachedData || !cachedData[key]) {
        return {by: 'workspace', val: ''}
    }

    const mp = cachedData[key]
    const by = mp.by ? mp.by : 'workspace'
    const val = mp[by]

    return {by: by, val: val}
}
export const setScriptFilters = async (siteId, productId, by, val) => {
    console.log('setScriptFilters')

    let cachedData = await getCache(settings.scriptFilters);
    if (!cachedData) cachedData = {}

    const key = `${siteId}-${productId}`

    const mp = cachedData[key] ? cachedData[key] : {}
    mp.by = by
    if (val) mp[by] = val

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
