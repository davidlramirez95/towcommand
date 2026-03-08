import { View, Text, StyleSheet } from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors } from '@/lib/theme/colors';
import { textStyles, fontFamily } from '@/lib/theme/typography';
import { spacing, borderRadius } from '@/lib/theme/spacing';
import { useBookingStore } from '@/stores/booking';

/**
 * Live tracking screen for an active booking.
 * Shows Mapbox map with provider location (via WebSocket),
 * status badge, ETA, and action buttons.
 *
 * TODO (Week 4): Wire up Mapbox map and WebSocket location updates.
 */
export default function BookingTrackingScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const activeBooking = useBookingStore((s) => s.activeBooking);

  return (
    <SafeAreaView style={styles.container}>
      {/* Map Area */}
      <View style={styles.mapArea}>
        <View style={styles.mapPlaceholder}>
          <Text style={styles.mapText}>Live Tracking Map</Text>
          <Text style={styles.mapSubtext}>Booking #{id}</Text>
        </View>
      </View>

      {/* Status Card */}
      <View style={styles.statusCard}>
        <View style={styles.statusHeader}>
          <View style={styles.statusBadge}>
            <Text style={styles.statusBadgeText}>
              {activeBooking?.status ?? 'LOADING'}
            </Text>
          </View>
          {activeBooking?.eta && (
            <Text style={styles.etaText}>ETA: {activeBooking.eta} min</Text>
          )}
        </View>

        {activeBooking?.providerName && (
          <View style={styles.providerInfo}>
            <View style={styles.providerAvatar}>
              <Text style={styles.providerAvatarText}>
                {activeBooking.providerName.charAt(0)}
              </Text>
            </View>
            <View>
              <Text style={styles.providerName}>{activeBooking.providerName}</Text>
              <Text style={styles.providerPhone}>{activeBooking.providerPhone}</Text>
            </View>
          </View>
        )}
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  mapArea: {
    flex: 1,
  },
  mapPlaceholder: {
    flex: 1,
    backgroundColor: colors.navy,
    alignItems: 'center',
    justifyContent: 'center',
  },
  mapText: {
    ...textStyles.h3,
    color: colors.white,
    opacity: 0.5,
  },
  mapSubtext: {
    ...textStyles.bodySmall,
    color: colors.teal,
    marginTop: spacing[1],
  },
  statusCard: {
    backgroundColor: colors.white,
    borderTopLeftRadius: borderRadius.xl,
    borderTopRightRadius: borderRadius.xl,
    padding: spacing[5],
    shadowColor: colors.navy,
    shadowOffset: { width: 0, height: -4 },
    shadowOpacity: 0.08,
    shadowRadius: 12,
    elevation: 4,
  },
  statusHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing[4],
  },
  statusBadge: {
    backgroundColor: colors.teal,
    borderRadius: borderRadius.pill,
    paddingVertical: spacing[1],
    paddingHorizontal: spacing[3],
  },
  statusBadgeText: {
    fontFamily: fontFamily.bold,
    fontSize: 10,
    color: colors.white,
    letterSpacing: 0.5,
  },
  etaText: {
    ...textStyles.label,
    color: colors.orange,
  },
  providerInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[3],
  },
  providerAvatar: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: colors.teal,
    alignItems: 'center',
    justifyContent: 'center',
  },
  providerAvatarText: {
    fontFamily: fontFamily.bold,
    fontSize: 18,
    color: colors.white,
  },
  providerName: {
    ...textStyles.label,
    color: colors.navy,
  },
  providerPhone: {
    ...textStyles.caption,
    color: colors.grey,
  },
});
