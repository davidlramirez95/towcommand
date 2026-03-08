/** @type {import('jest').Config} */
module.exports = {
  preset: 'jest-expo',
  setupFilesAfterEnv: ['./jest.setup.js'],
  // In a pnpm monorepo, dependencies are hoisted to root node_modules/.pnpm/
  // We need to transform React Native and Expo packages regardless of their location
  transformIgnorePatterns: [
    '<rootDir>/../../node_modules/(?!(.pnpm/[^/]+/node_modules/)?(react-native|@react-native|expo|@expo|@tanstack|zustand|react-native-mmkv|react-native-reanimated|react-native-gesture-handler|react-native-screens|react-native-safe-area-context|react-native-svg|@react-native-async-storage))',
  ],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/$1',
    '^aws-amplify/auth$': '<rootDir>/__mocks__/aws-amplify-auth.js',
    '^aws-amplify$': '<rootDir>/__mocks__/aws-amplify.js',
    '^react-native-mmkv$': '<rootDir>/__mocks__/react-native-mmkv.js',
    '^expo-secure-store$': '<rootDir>/__mocks__/expo-secure-store.js',
    '^expo-haptics$': '<rootDir>/__mocks__/expo-haptics.js',
    '^expo-router$': '<rootDir>/__mocks__/expo-router.js',
  },
  testMatch: ['**/__tests__/**/*.test.{ts,tsx}'],
  collectCoverageFrom: [
    'lib/**/*.{ts,tsx}',
    'stores/**/*.{ts,tsx}',
    'components/**/*.{ts,tsx}',
    'hooks/**/*.{ts,tsx}',
    '!**/*.d.ts',
    '!**/index.ts',
  ],
};
