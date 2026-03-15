import { View, Text, StyleSheet } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card, Button, InfoTip } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

const JOB_DETAILS = [
  ['Service', 'Flatbed Tow'],
  ['Distance', '12.4 km'],
  ['Duration', '45 min'],
];

export default function CompleteScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.content}>
        <View style={styles.successIcon}>
          <Text style={{ fontSize: 38 }}>✅</Text>
        </View>
        <Text style={styles.title}>Tow Complete!</Text>
        <Text style={styles.subtitle}>Your vehicle has been delivered safely</Text>

        <Card elevated style={styles.summaryCard}>
          <View style={styles.summaryRow}>
            <Text style={styles.summaryLabel}>Job ID</Text>
            <Text style={styles.summaryJobId}>TC-2026-00847</Text>
          </View>
          {JOB_DETAILS.map(([label, value]) => (
            <View key={label} style={styles.detailRow}>
              <Text style={styles.detailLabel}>{label}</Text>
              <Text style={styles.detailValue}>{value}</Text>
            </View>
          ))}
          <View style={styles.divider} />
          <View style={styles.totalRow}>
            <Text style={styles.totalLabel}>Total Paid</Text>
            <View style={{ alignItems: 'flex-end' }}>
              <Text style={styles.totalAmount}>₱1,850</Text>
              <Text style={styles.totalMethod}>via GCash 💚</Text>
            </View>
          </View>
        </Card>

        <InfoTip icon="⭐">
          <View>
            <Text style={styles.sukiTitle}>+1 Suki Point earned!</Text>
            <Text style={styles.sukiSub}>1 more booking to reach Gold tier</Text>
          </View>
        </InfoTip>

        <View style={styles.actions}>
          <Button onPress={() => router.push('/booking/rate')} fullWidth>
            ⭐ Rate Your Experience
          </Button>
          <Button variant="secondary" onPress={() => router.replace('/(tabs)')} fullWidth>
            Back to Home
          </Button>
        </View>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  content: { flex: 1, alignItems: 'center', justifyContent: 'center', padding: spacing[5] },
  successIcon: {
    width: 80,
    height: 80,
    borderRadius: 24,
    backgroundColor: colors.teal,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: spacing[4],
  },
  title: {
    fontFamily: fontFamily.bold,
    fontSize: 22,
    fontWeight: '800',
    color: colors.navy,
    marginBottom: spacing[1],
  },
  subtitle: {
    fontFamily: fontFamily.regular,
    fontSize: 12,
    color: colors.grey,
    marginBottom: spacing[6],
  },
  summaryCard: { width: '100%', marginBottom: spacing[3] },
  summaryRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 6,
  },
  summaryLabel: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
  summaryJobId: { fontFamily: fontFamily.semiBold, fontSize: 10, fontWeight: '600', color: colors.navy },
  detailRow: { flexDirection: 'row', justifyContent: 'space-between', paddingVertical: 4 },
  detailLabel: { fontFamily: fontFamily.regular, fontSize: 11, color: colors.grey },
  detailValue: { fontFamily: fontFamily.semiBold, fontSize: 11, fontWeight: '600', color: colors.navy },
  divider: { height: 1, backgroundColor: colors.greyLight, marginVertical: 8 },
  totalRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  totalLabel: { fontFamily: fontFamily.bold, fontSize: 13, fontWeight: '700', color: colors.navy },
  totalAmount: { fontFamily: fontFamily.bold, fontSize: 18, fontWeight: '800', color: colors.orange },
  totalMethod: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
  sukiTitle: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700', color: '#F57F17' },
  sukiSub: { fontFamily: fontFamily.regular, fontSize: 10, color: '#795548' },
  actions: { width: '100%', gap: spacing[2], marginTop: spacing[5] },
});
