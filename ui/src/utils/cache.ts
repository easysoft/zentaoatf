
import {getCache, setCache} from './localCache';
import settings from '@/config/settings';

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
  let mp = await getCache(settings.currProductIdBySite);
  if (!mp) mp = {}
  mp[currSiteId + ''] = currProductId
  await setCache(settings.currProductIdBySite, mp);
}

// script filters
export const getScriptFilters = async () => {
  const mp = await getCache(settings.scriptFilters);

  if (!mp) {
    return {by: 'workspace', val: ''}
  }

  const by = mp.by
  const val = mp[by]

  return {by: by, val: val}
}
export const setScriptFilters = async (by, val) => {
  let mp = await getCache(settings.scriptFilters);

  if (!mp) mp = {}
  mp.by = by
  mp[by] = val

  await setCache(settings.scriptFilters, mp);
}
