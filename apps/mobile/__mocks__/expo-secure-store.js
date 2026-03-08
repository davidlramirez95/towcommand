const store = new Map();

module.exports = {
  setItemAsync: jest.fn((key, value) => {
    store.set(key, value);
    return Promise.resolve();
  }),
  getItemAsync: jest.fn((key) => {
    return Promise.resolve(store.get(key) || null);
  }),
  deleteItemAsync: jest.fn((key) => {
    store.delete(key);
    return Promise.resolve();
  }),
  __clearStore: () => store.clear(),
};
