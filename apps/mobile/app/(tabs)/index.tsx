import { View, Text, StyleSheet, Pressable } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { colors } from '@/lib/theme/colors';
import { textStyles, fontFamily } from '@/lib/theme/typography';
import { spacing, borderRadius } from '@/lib/theme/spacing';
import { useAuth } from '@/hooks/useAuth';

export default function HomeScreen() {
  const { user } = useAuth();
  const firstName = user?.fullName?.split(' ')[0] ?? 'there';

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <View>
          <Text style={styles.greeting}>Mabuhay, {firstName}!</Text>
          <Text style={styles.subtitle}>Need help on the road?</Text>
        </View>
        <Pressable
          style={styles.sosButton}
          onPress={() => router.push('/sos')}
          accessibilityLabel="Emergency SOS"
          accessibilityRole="button"
        >
          <Text style={styles.sosText}>SOS</Text>
        </Pressable>
      </View>

      {/* Map Placeholder */}
      <View style={styles.mapContainer}>
        <View style={styles.mapPlaceholder}>
          <Text style={styles.mapPlaceholderText}>Map loads here</Text>
          <Text style={styles.mapPlaceholderSubtext}>Mapbox GL integration</Text>
        </View>
      </View>

      {/* Quick Actions */}
      <View style={styles.actions}>
        <Card elevated onPress={() => router.push('/booking/service')} style={styles.mainAction}>
          <Text style={styles.actionTitle}>Request a Tow</Text>
          <Text style={styles.actionSubtitle}>Get help in minutes</Text>
        </Card>

        <View style={styles.secondaryActions}>
          <Card
            onPress={() => router.push('/booking/diagnose')}
            style={styles.secondaryAction}
          >
            <View style={styles.aiBadge}>
              <Text style={styles.aiBadgeText}>AI</Text>
            </View>
            <Text style={styles.secondaryTitle}>Diagnose</Text>
            <Text style={styles.secondarySubtitle}>What's wrong?</Text>
          </Card>
          <Card
            onPress={() => router.push('/(tabs)/history')}
            style={styles.secondaryAction}
          >
            <Text style={styles.secondaryTitle}>History</Text>
            <Text style={styles.secondarySubtitle}>Past trips</Text>
          </Card>
        </View>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing[5],
    paddingVertical: spacing[3],
  },
  greeting: {
    ...textStyles.h2,
    color: colors.navy,
  },
  subtitle: {
    ...textStyles.bodySmall,
    color: colors.grey,
    marginTop: 2,
  },
  sosButton: {
    width: 48,
    height: 48,
    borderRadius: borderRadius.full,
    backgroundColor: colors.coral,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: colors.coral,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.35,
    shadowRadius: 8,
    elevation: 4,
  },
  sosText: {
    fontFamily: fontFamily.bold,
    fontSize: 12,
    color: colors.white,
    letterSpacing: 1,
  },
  mapContainer: {
    flex: 1,
    marginHorizontal: spacing[4],
    marginVertical: spacing[2],
    borderRadius: borderRadius.xl,
    overflow: 'hidden',
  },
  mapPlaceholder: {
    flex: 1,
    backgroundColor: colors.navy,
    alignItems: 'center',
    justifyContent: 'center',
  },
  mapPlaceholderText: {
    ...textStyles.h3,
    color: colors.white,
    opacity: 0.5,
  },
  mapPlaceholderSubtext: {
    ...textStyles.caption,
    color: colors.teal,
    marginTop: spacing[1],
  },
  actions: {
    paddingHorizontal: spacing[4],
    paddingBottom: spacing[3],
    gap: spacing[3],
  },
  mainAction: {
    backgroundColor: colors.orange,
    borderColor: colors.orange,
  },
  actionTitle: {
    ...textStyles.h3,
    color: colors.white,
  },
  actionSubtitle: {
    ...textStyles.bodySmall,
    color: 'rgba(255,255,255,0.8)',
    marginTop: 2,
  },
  secondaryActions: {
    flexDirection: 'row',
    gap: spacing[3],
  },
  secondaryAction: {
    flex: 1,
  },
  aiBadge: {
    backgroundColor: '#667eea',
    borderRadius: 6,
    paddingHorizontal: 8,
    paddingVertical: 2,
    alignSelf: 'flex-start',
    marginBottom: spacing[2],
  },
  aiBadgeText: {
    fontFamily: fontFamily.bold,
    fontSize: 8,
    color: colors.white,
    letterSpacing: 0.5,
  },
  secondaryTitle: {
    ...textStyles.label,
    color: colors.navy,
  },
  secondarySubtitle: {
    ...textStyles.caption,
    color: colors.grey,
    marginTop: 2,
  },
});
