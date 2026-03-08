import { View, Text, StyleSheet, Pressable, Alert } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { colors } from '@/lib/theme/colors';
import { textStyles, fontFamily } from '@/lib/theme/typography';
import { spacing, borderRadius } from '@/lib/theme/spacing';
import { useAuth } from '@/hooks/useAuth';

export default function ProfileScreen() {
  const { user, signOut } = useAuth();

  const handleSignOut = async () => {
    Alert.alert('Sign Out', 'Are you sure you want to sign out?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Sign Out',
        style: 'destructive',
        onPress: async () => {
          await signOut();
          router.replace('/(auth)/login');
        },
      },
    ]);
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Account</Text>
      </View>

      {/* Profile Card */}
      <Card elevated style={styles.profileCard}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>
            {user?.fullName?.charAt(0)?.toUpperCase() ?? 'U'}
          </Text>
        </View>
        <Text style={styles.name}>{user?.fullName ?? 'User'}</Text>
        <Text style={styles.email}>{user?.email ?? ''}</Text>
        <Text style={styles.phone}>{user?.phone ?? ''}</Text>
      </Card>

      {/* Menu Items */}
      <View style={styles.menu}>
        <MenuItem label="My Vehicles" onPress={() => {}} />
        <MenuItem label="Payment Methods" onPress={() => {}} />
        <MenuItem label="Notifications" onPress={() => {}} />
        <MenuItem label="Help & Support" onPress={() => {}} />
        <MenuItem label="About TowCommand" onPress={() => {}} />
      </View>

      <View style={styles.signOutSection}>
        <Button variant="ghost" onPress={handleSignOut} fullWidth>
          Sign Out
        </Button>
      </View>
    </SafeAreaView>
  );
}

function MenuItem({ label, onPress }: { label: string; onPress: () => void }) {
  return (
    <Pressable
      style={styles.menuItem}
      onPress={onPress}
      accessibilityRole="button"
      accessibilityLabel={label}
    >
      <Text style={styles.menuLabel}>{label}</Text>
      <Text style={styles.menuArrow}>›</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  header: {
    paddingHorizontal: spacing[5],
    paddingTop: spacing[3],
    paddingBottom: spacing[4],
  },
  title: {
    ...textStyles.h2,
    color: colors.navy,
  },
  profileCard: {
    marginHorizontal: spacing[4],
    alignItems: 'center',
    paddingVertical: spacing[6],
    marginBottom: spacing[4],
  },
  avatar: {
    width: 64,
    height: 64,
    borderRadius: borderRadius.full,
    backgroundColor: colors.teal,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: spacing[3],
  },
  avatarText: {
    fontFamily: fontFamily.bold,
    fontSize: 24,
    color: colors.white,
  },
  name: {
    ...textStyles.h3,
    color: colors.navy,
  },
  email: {
    ...textStyles.bodySmall,
    color: colors.grey,
    marginTop: 2,
  },
  phone: {
    ...textStyles.bodySmall,
    color: colors.grey,
    marginTop: 2,
  },
  menu: {
    marginHorizontal: spacing[4],
    backgroundColor: colors.white,
    borderRadius: borderRadius.lg,
    borderWidth: 1,
    borderColor: colors.greyLight,
    overflow: 'hidden',
  },
  menuItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: spacing[4],
    paddingHorizontal: spacing[4],
    borderBottomWidth: 1,
    borderBottomColor: colors.greyLight,
  },
  menuLabel: {
    ...textStyles.body,
    color: colors.navy,
  },
  menuArrow: {
    fontSize: 20,
    color: colors.grey,
  },
  signOutSection: {
    marginTop: 'auto',
    paddingHorizontal: spacing[4],
    paddingBottom: spacing[4],
  },
});
