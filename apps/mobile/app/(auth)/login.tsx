import { useState } from 'react';
import { View, Text, StyleSheet, Pressable, Alert, KeyboardAvoidingView, Platform, ScrollView } from 'react-native';
import { Link, router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { colors } from '@/lib/theme/colors';
import { textStyles, fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';
import { useAuth } from '@/hooks/useAuth';

export default function LoginScreen() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const { signIn, socialSignIn } = useAuth();

  const handleLogin = async () => {
    if (!email || !password) {
      Alert.alert('Error', 'Please fill in all fields');
      return;
    }
    setLoading(true);
    try {
      await signIn({ email, password });
      router.replace('/(tabs)');
    } catch (err) {
      Alert.alert('Login Failed', err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleSocialLogin = async (provider: 'Google' | 'Facebook' | 'Apple') => {
    try {
      await socialSignIn(provider);
    } catch (err) {
      Alert.alert('Login Failed', err instanceof Error ? err.message : 'An error occurred');
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.keyboardView}
      >
        <ScrollView contentContainerStyle={styles.scroll} keyboardShouldPersistTaps="handled">
          {/* Logo */}
          <View style={styles.logoSection}>
            <Text style={styles.logoText}>TowCommand</Text>
            <Text style={styles.logoSubtext}>PILIPINAS</Text>
          </View>

          {/* Welcome Text */}
          <View style={styles.welcomeSection}>
            <Text style={styles.welcomeTitle}>Get help on the road,{'\n'}anytime, anywhere.</Text>
            <Text style={styles.welcomeSubtitle}>Sign in to book a tow truck in seconds</Text>
          </View>

          {/* Social Login Buttons */}
          <View style={styles.socialButtons}>
            <Pressable
              style={styles.socialButton}
              onPress={() => handleSocialLogin('Google')}
              accessibilityLabel="Continue with Google"
              accessibilityRole="button"
            >
              <Text style={styles.socialButtonText}>Continue with Google</Text>
            </Pressable>
            <Pressable
              style={[styles.socialButton, styles.facebookButton]}
              onPress={() => handleSocialLogin('Facebook')}
              accessibilityLabel="Continue with Facebook"
              accessibilityRole="button"
            >
              <Text style={[styles.socialButtonText, styles.whiteText]}>Continue with Facebook</Text>
            </Pressable>
            <Pressable
              style={[styles.socialButton, styles.appleButton]}
              onPress={() => handleSocialLogin('Apple')}
              accessibilityLabel="Continue with Apple"
              accessibilityRole="button"
            >
              <Text style={[styles.socialButtonText, styles.whiteText]}>Continue with Apple</Text>
            </Pressable>
          </View>

          {/* Divider */}
          <View style={styles.divider}>
            <View style={styles.dividerLine} />
            <Text style={styles.dividerText}>or sign in with email</Text>
            <View style={styles.dividerLine} />
          </View>

          {/* Email/Password */}
          <View style={styles.form}>
            <Input
              label="Email"
              value={email}
              onChangeText={setEmail}
              placeholder="juan@example.com"
              keyboardType="email-address"
              autoCapitalize="none"
              autoComplete="email"
            />
            <Input
              label="Password"
              value={password}
              onChangeText={setPassword}
              placeholder="Enter your password"
              secureTextEntry
              autoComplete="password"
            />
            <Button onPress={handleLogin} loading={loading} fullWidth>
              Sign In
            </Button>
          </View>

          {/* Sign Up Link */}
          <View style={styles.signupSection}>
            <Text style={styles.signupText}>Don't have an account? </Text>
            <Link href="/(auth)/signup" asChild>
              <Pressable accessibilityRole="link">
                <Text style={styles.signupLink}>Sign Up</Text>
              </Pressable>
            </Link>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  keyboardView: {
    flex: 1,
  },
  scroll: {
    flexGrow: 1,
    paddingHorizontal: spacing[6],
    justifyContent: 'center',
  },
  logoSection: {
    alignItems: 'center',
    marginBottom: spacing[6],
  },
  logoText: {
    fontFamily: fontFamily.bold,
    fontSize: 28,
    color: colors.orange,
  },
  logoSubtext: {
    fontFamily: fontFamily.semiBold,
    fontSize: 10,
    color: colors.teal,
    letterSpacing: 4,
    marginTop: -2,
  },
  welcomeSection: {
    alignItems: 'center',
    marginBottom: spacing[8],
  },
  welcomeTitle: {
    ...textStyles.h2,
    color: colors.navy,
    textAlign: 'center',
    lineHeight: 28,
  },
  welcomeSubtitle: {
    ...textStyles.bodySmall,
    color: colors.grey,
    marginTop: spacing[2],
  },
  socialButtons: {
    gap: spacing[2],
    marginBottom: spacing[6],
  },
  socialButton: {
    width: '100%',
    paddingVertical: spacing[3],
    borderRadius: 14,
    borderWidth: 1.5,
    borderColor: colors.greyLight,
    backgroundColor: colors.white,
    alignItems: 'center',
  },
  facebookButton: {
    backgroundColor: '#1877F2',
    borderWidth: 0,
  },
  appleButton: {
    backgroundColor: colors.navy,
    borderWidth: 0,
  },
  socialButtonText: {
    fontFamily: fontFamily.semiBold,
    fontSize: 14,
    color: colors.navy,
  },
  whiteText: {
    color: colors.white,
  },
  divider: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing[6],
    gap: spacing[3],
  },
  dividerLine: {
    flex: 1,
    height: 1,
    backgroundColor: colors.greyLight,
  },
  dividerText: {
    ...textStyles.caption,
    color: colors.grey,
  },
  form: {
    gap: spacing[4],
    marginBottom: spacing[6],
  },
  signupSection: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
  },
  signupText: {
    ...textStyles.bodySmall,
    color: colors.grey,
  },
  signupLink: {
    ...textStyles.bodySmall,
    color: colors.orange,
    fontFamily: fontFamily.bold,
  },
});
