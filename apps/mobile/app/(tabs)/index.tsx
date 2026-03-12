import { View, Text, StyleSheet, Pressable, ScrollView } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card } from '@/components/ui/Card';
import { colors } from '@/lib/theme/colors';
import { textStyles, fontFamily } from '@/lib/theme/typography';
import { spacing, borderRadius } from '@/lib/theme/spacing';
import { useAuth } from '@/hooks/useAuth';

const QUICK_SERVICES = [
  { icon: '🚛', label: 'Tow', color: '#FFF3EB' },
  { icon: '⛽', label: 'Fuel', color: '#E8F8F0' },
  { icon: '🔋', label: 'Jumpstart', color: '#EBF0FF' },
  { icon: '🔧', label: 'Mechanic', color: '#FFF8E1' },
];

export default function HomeScreen() {
  const { user } = useAuth();
  const firstName = user?.fullName?.split(' ')[0] ?? 'there';

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scroll}>
        {/* Header */}
        <View style={styles.header}>
          <View>
            <Text style={styles.greetingSub}>Magandang hapon 🌤️</Text>
            <Text style={styles.greeting}>{firstName}!</Text>
          </View>
          <View style={styles.headerRight}>
            <Pressable style={styles.notifBtn}>
              <Text style={{ fontSize: 16 }}>🔔</Text>
            </Pressable>
            <Pressable
              style={styles.sosButton}
              onPress={() => router.push('/sos')}
              accessibilityLabel="Emergency SOS"
              accessibilityRole="button"
            >
              <Text style={styles.sosText}>SOS</Text>
            </Pressable>
          </View>
        </View>

        {/* AI Diagnosis Card */}
        <Card elevated onPress={() => router.push('/booking/diagnose')} style={styles.aiCard}>
          <View style={styles.aiRow}>
            <View style={styles.aiIcon}>
              <Text style={{ fontSize: 26 }}>🤖</Text>
            </View>
            <View style={{ flex: 1 }}>
              <Text style={styles.aiTitle}>What's wrong with your car?</Text>
              <Text style={styles.aiSub}>AI will diagnose & find the cheapest fix</Text>
            </View>
            <Text style={{ fontSize: 20, color: colors.gold }}>→</Text>
          </View>
        </Card>

        {/* Map Placeholder */}
        <View style={styles.mapContainer}>
          <View style={styles.mapPlaceholder}>
            <Text style={{ position: 'absolute', top: 25, left: '25%', fontSize: 18 }}>🚛</Text>
            <Text style={{ position: 'absolute', top: 80, left: '70%', fontSize: 16 }}>🚛</Text>
            <Text style={{ position: 'absolute', top: 100, left: '35%', fontSize: 14, opacity: 0.6 }}>🚛</Text>
            <View style={styles.mapPinOuter}>
              <View style={styles.mapPinInner} />
            </View>
          </View>
          <View style={styles.mapBadge}>
            <View style={styles.mapBadgeDot} />
            <Text style={styles.mapBadgeText}>3 trucks nearby</Text>
          </View>
        </View>

        {/* Quick Services */}
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Quick Services</Text>
          <Text style={styles.seeAll}>See all</Text>
        </View>
        <View style={styles.quickGrid}>
          {QUICK_SERVICES.map((s) => (
            <Pressable
              key={s.label}
              onPress={() => router.push('/booking/service')}
              style={[styles.quickItem, { backgroundColor: s.color }]}
            >
              <Text style={{ fontSize: 24 }}>{s.icon}</Text>
              <Text style={styles.quickLabel}>{s.label}</Text>
            </Pressable>
          ))}
        </View>

        {/* Suki Card */}
        <Card style={styles.sukiCard}>
          <View style={styles.sukiRow}>
            <Text style={{ fontSize: 22 }}>⭐</Text>
            <View style={{ flex: 1 }}>
              <Text style={styles.sukiTitle}>Suki Silver Member</Text>
              <Text style={styles.sukiSub}>2 more bookings for Gold • 5% off all services</Text>
            </View>
          </View>
        </Card>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  scroll: { paddingHorizontal: spacing[5], paddingBottom: spacing[3] },
  // Header
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: spacing[2],
    marginBottom: spacing[3],
  },
  greetingSub: { fontFamily: fontFamily.regular, fontSize: 11, color: colors.grey },
  greeting: { ...textStyles.h2, color: colors.navy },
  headerRight: { flexDirection: 'row', gap: 8 },
  notifBtn: {
    width: 36,
    height: 36,
    borderRadius: borderRadius.md,
    backgroundColor: colors.light,
    alignItems: 'center',
    justifyContent: 'center',
  },
  sosButton: {
    width: 36,
    height: 36,
    borderRadius: borderRadius.md,
    backgroundColor: colors.coral,
    alignItems: 'center',
    justifyContent: 'center',
  },
  sosText: { fontFamily: fontFamily.bold, fontSize: 9, color: colors.white, letterSpacing: 0.5 },
  // AI Card
  aiCard: {
    backgroundColor: colors.navy,
    borderColor: colors.navy,
    marginBottom: spacing[3],
    padding: 18,
  },
  aiRow: { flexDirection: 'row', alignItems: 'center', gap: 12 },
  aiIcon: {
    width: 48,
    height: 48,
    borderRadius: 14,
    backgroundColor: 'rgba(255,107,53,0.2)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  aiTitle: { fontFamily: fontFamily.bold, fontSize: 14, fontWeight: '700', color: colors.white },
  aiSub: { fontFamily: fontFamily.regular, fontSize: 11, color: 'rgba(255,255,255,0.6)', marginTop: 2 },
  // Map
  mapContainer: {
    height: 145,
    borderRadius: 18,
    overflow: 'hidden',
    marginBottom: spacing[3],
    backgroundColor: '#E8E4DE',
    position: 'relative',
  },
  mapPlaceholder: { flex: 1, position: 'relative' },
  mapPinOuter: {
    position: 'absolute',
    top: '35%',
    left: '53%',
    width: 24,
    height: 24,
    borderRadius: 12,
    backgroundColor: colors.orange,
    borderWidth: 3,
    borderColor: colors.white,
    alignItems: 'center',
    justifyContent: 'center',
  },
  mapPinInner: { width: 6, height: 6, borderRadius: 3, backgroundColor: colors.white },
  mapBadge: {
    position: 'absolute',
    top: 10,
    right: 10,
    backgroundColor: 'rgba(11,29,51,0.85)',
    borderRadius: 10,
    paddingHorizontal: 10,
    paddingVertical: 6,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 5,
  },
  mapBadgeDot: { width: 6, height: 6, borderRadius: 3, backgroundColor: colors.green },
  mapBadgeText: { fontFamily: fontFamily.semiBold, fontSize: 10, fontWeight: '600', color: colors.white },
  // Quick Services
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 10,
  },
  sectionTitle: { fontFamily: fontFamily.bold, fontSize: 14, fontWeight: '700', color: colors.navy },
  seeAll: { fontFamily: fontFamily.semiBold, fontSize: 11, fontWeight: '600', color: colors.orange },
  quickGrid: {
    flexDirection: 'row',
    gap: 8,
    marginBottom: spacing[4],
  },
  quickItem: {
    flex: 1,
    alignItems: 'center',
    gap: 6,
    paddingVertical: 14,
    paddingHorizontal: 8,
    borderRadius: 14,
  },
  quickLabel: { fontFamily: fontFamily.semiBold, fontSize: 10, fontWeight: '600', color: colors.navy },
  // Suki
  sukiCard: {
    backgroundColor: colors.cream,
    borderWidth: 1.5,
    borderColor: 'rgba(245,166,35,0.2)',
  },
  sukiRow: { flexDirection: 'row', alignItems: 'center', gap: 10 },
  sukiTitle: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700', color: colors.navy },
  sukiSub: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
});
