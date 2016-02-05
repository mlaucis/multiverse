let localStorageMock = function () {
  let store = {}

  return {
    clear: () => {
      store = {}
    },

    getItem: (key) => {
      return store[key]
    },

    removeItem: (key) => {
      delete store[key]
    },

    setItem: (key, value) => {
      store[key] = value
    }
  }
}

Object.defineProperty(window, 'localStorage', { value: localStorageMock() })
