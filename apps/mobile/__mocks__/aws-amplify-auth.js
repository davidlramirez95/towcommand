// Mock for aws-amplify/auth (used by useAuth hook and API client)
const mockSession = {
  tokens: {
    idToken: { toString: () => 'mock-jwt-token-123' },
    accessToken: { toString: () => 'mock-access-token' },
  },
  credentials: {},
};

const mockUser = {
  userId: 'user-123',
  username: 'test@example.com',
  signInDetails: { loginId: 'test@example.com' },
};

module.exports = {
  fetchAuthSession: jest.fn(() => Promise.resolve(mockSession)),
  getCurrentUser: jest.fn(() => Promise.resolve(mockUser)),
  fetchUserAttributes: jest.fn(() =>
    Promise.resolve({
      email: 'test@example.com',
      phone_number: '+639171234567',
      name: 'Juan Cruz',
      'custom:user_type': 'customer',
    }),
  ),
  signIn: jest.fn(() => Promise.resolve({ isSignedIn: true, nextStep: { signInStep: 'DONE' } })),
  signUp: jest.fn(() => Promise.resolve({ isSignUpComplete: true, nextStep: { signUpStep: 'DONE' } })),
  signOut: jest.fn(() => Promise.resolve()),
  signInWithRedirect: jest.fn(() => Promise.resolve()),
};
