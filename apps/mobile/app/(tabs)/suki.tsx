import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card, SectionLabel } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily, textStyles } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

const BENEFITS = [
  { label: '5% off all services', status: 'Active now', active: true },
  { label: 'Priority matching', status: 'Silver+', active: true },
  { label: 'Free condition report', status: 'Silver+', active: true },
  { label: '10% off + VIP support', status: 'Gold tier', active: false },
];

const POINTS_HISTORY = [
  { date: 'Feb 20', amount: '+1 point', reason: 'Flatbed Tow' },
  { date: 'Feb 14', amount: '+1 point', reason: 'Jumpstart' },
  { date: 'Jan 28', amount: '+1 point', reason: 'Fuel Delivery' },
  { date: 'Jan 15', amount: '+1 point', reason: 'Tire Change' },
];

export default function SukiScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.headerSection}>
        <Text style={styles.pageTitle}>Suki Rewards</Text>
      </View>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Tier Card */}
        <Card elevated style={styles.tierCard}>
          <View style={styles.tierRow}>
            <View style={styles.tierIcon}>
              <Text style={{ fontSize: 26 }}>⭐</Text>
            </View>
            <View>
              <Text style={styles.tierName}>Silver Member</Text>
              <Text style={styles.tierProgress}>4 of 6 bookings to Gold</Text>
            </View>
          </View>
          <View style={styles.progressOuter}>
            <View style={[styles.progressInner, { width: '66%' }]} />
          </View>
          <View style={styles.progressLabels}>
            <Text style={styles.progressLabelLeft}>Silver</Text>
            <Text style={styles.progressLabelRight}>Gold →</Text>
          </View>
        </Card>

        <SectionLabel>Your benefits</SectionLabel>
        {BENEFITS.map((b) => (
          <Card key={b.label} style={styles.benefitCard}>
            <View style={styles.benefitRow}>
              <Text style={styles.benefitLabel}>{b.label}</Text>
              <View style={[styles.statusBadge, b.active ? styles.statusActive : styles.statusLocked]}>
                <Text style={[styles.statusText, b.active ? styles.statusTextActive : styles.statusTextLocked]}>
                  {b.status}
                </Text>
              </View>
            </View>
          </Card>
        ))}

        <SectionLabel style={{ marginTop: spacing[3] }}>Points history</SectionLabel>
        {POINTS_HISTORY.map((p) => (
          <Card key={p.date + p.reason} style={styles.historyCard}>
            <View style={styles.historyRow}>
              <View>
                <Text style={styles.historyReason}>{p.reason}</Text>
                <Text style={styles.historyDate}>{p.date}</Text>
              </View>
              <Text style={styles.historyAmount}>{p.amount}</Text>
            </View>
          </Card>
        ))}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  headerSection: { paddingHorizontal: spacing[5], paddingTop: spacing[1], paddingBottom: spacing[2] },
  pageTitle: { ...textStyles.h2, color: colors.navy },
  scrollContent: { padding: spacing[4], gap: spacing[2] },
  // Tier
  tierCard: { backgroundColor: colors.navy, borderWidth: 0, padding: spacing[5], marginBottom: spacing[3] },
  tierRow: { flexDirection: 'row', alignItems: 'center', gap: 12 },
  tierIcon: {
    width: 50,
    height: 50,
    borderRadius: 16,
    backgroundColor: 'rgba(245,166,35,0.2)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  tierName: { fontFamily: fontFamily.bold, fontSize: 16, fontWeight: '800', color: colors.white },
  tierProgress: { fontFamily: fontFamily.regular, fontSize: 11, color: 'rgba(255,255,255,0.5)' },
  progressOuter: {
    height: 8,
    borderRadius: 4,
    backgroundColor: 'rgba(255,255,255,0.1)',
    marginTop: spacing[3],
  },
  progressInner: {
    height: '100%',
    borderRadius: 4,
    backgroundColor: colors.orange,
  },
  progressLabels: { flexDirection: 'row', justifyContent: 'space-between', marginTop: 6 },
  progressLabelLeft: { fontFamily: fontFamily.regular, fontSize: 9, color: 'rgba(255,255,255,0.4)' },
  progressLabelRight: { fontFamily: fontFamily.regular, fontSize: 9, color: colors.gold },
  // Benefits
  benefitCard: { marginBottom: spacing[1], padding: spacing[3] },
  benefitRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between' },
  benefitLabel: { fontFamily: fontFamily.semiBold, fontSize: 12, fontWeight: '600', color: colors.navy },
  statusBadge: { paddingHorizontal: 8, paddingVertical: 2, borderRadius: 4 },
  statusActive: { backgroundColor: 'rgba(0,196,140,0.1)' },
  statusLocked: { backgroundColor: 'rgba(142,155,174,0.1)' },
  statusText: { fontFamily: fontFamily.bold, fontSize: 9, fontWeight: '700' },
  statusTextActive: { color: colors.green },
  statusTextLocked: { color: colors.grey },
  // History
  historyCard: { marginBottom: spacing[1], padding: spacing[3] },
  historyRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  historyReason: { fontFamily: fontFamily.semiBold, fontSize: 11, fontWeight: '600', color: colors.navy },
  historyDate: { fontFamily: fontFamily.regular, fontSize: 9, color: colors.grey },
  historyAmount: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700', color: colors.green },
});
