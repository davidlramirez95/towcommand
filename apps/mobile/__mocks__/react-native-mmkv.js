const store = new Map();

class MMKV {
  // eslint-disable-next-line no-unused-vars
  constructor(options) {
    // Don't auto-clear — let tests call clearAll() or reset stores explicitly
  }
  getString(key) { return store.get(key); }
  set(key, value) { store.set(key, value); }
  delete(key) { store.delete(key); }
  contains(key) { return store.has(key); }
  clearAll() { store.clear(); }
}

// Helper for explicit test cleanup
module.exports = { MMKV, __clearStore: () => store.clear() };
