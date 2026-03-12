import { Platform } from 'react-native';
import { StateStorage } from 'zustand/middleware';

/**
 * MMKV storage with web fallback (localStorage).
 * react-native-mmkv is native-only; web uses localStorage for Playwright E2E.
 */

let storage: { getString: (k: string) => string | undefined; set: (k: string, v: string) => void; delete: (k: string) => void };

if (Platform.OS === 'web') {
  storage = {
    getString: (name: string) => {
      const v = localStorage.getItem(name);
      return v ?? undefined;
    },
    set: (name: string, value: string) => {
      localStorage.setItem(name, value);
    },
    delete: (name: string) => {
      localStorage.removeItem(name);
    },
  };
} else {
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  const { MMKV } = require('react-native-mmkv');
  storage = new MMKV({ id: 'towcommand-storage' });
}

export { storage };

/**
 * Zustand-compatible storage adapter backed by MMKV (native) or localStorage (web).
 * Synchronous reads for fast hydration on app launch.
 */
export const mmkvStorage: StateStorage = {
  getItem: (name: string) => {
    const value = storage.getString(name);
    return value ?? null;
  },
  setItem: (name: string, value: string) => {
    storage.set(name, value);
  },
  removeItem: (name: string) => {
    storage.delete(name);
  },
};
