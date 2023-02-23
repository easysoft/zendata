import localforage from 'localforage';

export const getCache = async (key) => {
  return await localforage.getItem(key);
};

export const setCache = async (key, val) => {
  try {
    await localforage.setItem(key, val);
    return true;
  } catch (error) {
    return false;
  }
};
