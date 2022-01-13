
import localforage from 'localforage';
import settings from '@/config/settings';

export const getCache = async (key: string): Promise<string | null> => {
  return await localforage.getItem(key);
};

export const setCache = async (key: string, val: string): Promise<boolean> => {
  try {
    await localforage.setItem(key, val);
    return true;
  } catch (error) {
    return false;
  }
};

export const removeCache = async (key: string): Promise<boolean> => {
  try {
    await localforage.removeItem(settings.siteTokenKey);
    return true;
  } catch (error) {
    return false;
  }
};