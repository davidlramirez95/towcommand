import { useEffect, useCallback } from 'react';
import {
  signIn,
  signOut,
  signUp,
  getCurrentUser,
  fetchAuthSession,
  signInWithRedirect,
} from 'aws-amplify/auth';
import { useAuthStore } from '@/stores/auth';
import { api } from '@/lib/api';

interface SignUpInput {
  email: string;
  password: string;
  phone: string;
  fullName: string;
}

interface SignInInput {
  email: string;
  password: string;
}

export function useAuth() {
  const { user, isAuthenticated, isLoading, setUser, setLoading, reset } = useAuthStore();

  // Check auth state on mount
  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = useCallback(async () => {
    setLoading(true);
    try {
      const cognitoUser = await getCurrentUser();
      const session = await fetchAuthSession();
      const claims = session.tokens?.idToken?.payload;

      setUser({
        id: cognitoUser.userId,
        email: (claims?.email as string) ?? '',
        phone: (claims?.phone_number as string) ?? '',
        fullName: (claims?.name as string) ?? '',
        userType: ((claims?.['custom:user_type'] as string) ?? 'customer') as 'customer' | 'provider' | 'admin',
      });
    } catch {
      setUser(null);
    }
  }, [setUser, setLoading]);

  const handleSignIn = useCallback(async (input: SignInInput) => {
    const result = await signIn({
      username: input.email,
      password: input.password,
    });
    if (result.isSignedIn) {
      await checkAuth();
    }
    return result;
  }, [checkAuth]);

  const handleSignUp = useCallback(async (input: SignUpInput) => {
    const result = await signUp({
      username: input.email,
      password: input.password,
      options: {
        userAttributes: {
          email: input.email,
          phone_number: input.phone,
          name: input.fullName,
        },
      },
    });
    return result;
  }, []);

  const handleSocialSignIn = useCallback(async (provider: 'Google' | 'Facebook' | 'Apple') => {
    await signInWithRedirect({ provider });
  }, []);

  const handleSignOut = useCallback(async () => {
    await signOut();
    reset();
  }, [reset]);

  return {
    user,
    isAuthenticated,
    isLoading,
    signIn: handleSignIn,
    signUp: handleSignUp,
    signOut: handleSignOut,
    socialSignIn: handleSocialSignIn,
    checkAuth,
  };
}
