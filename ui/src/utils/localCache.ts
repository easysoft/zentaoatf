
import localforage from 'localforage';
import settings from '@/config/settings';

export const getCache = async (key: string): Promise<any | null> => {
  return localforage.getItem(key);
};

export const setCache = async (key: string, val: any): Promise<boolean> => {
  try {
    await localforage.setItem(key, val);
    return true;
  } catch (error) {
    return false;
  }
};
