const store = new Map();

class MMKV {
  constructor() {
    store.clear();
  }
  getString(key) { return store.get(key); }
  set(key, value) { store.set(key, value); }
  delete(key) { store.delete(key); }
  contains(key) { return store.has(key); }
  clearAll() { store.clear(); }
}

module.exports = { MMKV };
