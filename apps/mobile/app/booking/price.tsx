import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { BackHeader, Card, SectionLabel, InfoTip, Button } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';
import { useBookingStore } from '@/stores/booking';

const PRICE_LINES = [
  { label: 'Base fare', value: '₱800' },
  { label: 'Distance (4.2 km × ₱100)', value: '₱420' },
  { label: 'Weight surcharge (SUV)', value: '₱250' },
  { label: 'Time surcharge (peak)', value: '₱180' },
  { label: 'Platform fee', value: '₱200' },
];

export default function PriceScreen() {
  const { setPaymentMethod, setMatchingState } = useBookingStore();

  const handleBook = () => {
    setPaymentMethod({ type: 'gcash', label: 'GCash', last4: '8847' });
    setMatchingState('searching');
    router.push('/booking/matching');
  };

  return (
    <SafeAreaView style={styles.container}>
      <BackHeader title="Price Estimate" onBack={() => router.back()} />
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Card elevated style={styles.priceCard}>
          <View style={styles.totalSection}>
            <Text style={styles.estimateLabel}>Estimated Total</Text>
            <Text style={styles.totalPrice}>₱1,850</Text>
            <Text style={styles.regulation}>MMDA Reg. 24-004 compliant pricing</Text>
          </View>
          <View style={styles.divider} />
          {PRICE_LINES.map(({ label, value }) => (
            <View key={label} style={styles.lineItem}>
              <Text style={styles.lineLabel}>{label}</Text>
              <Text style={styles.lineValue}>{value}</Text>
            </View>
          ))}
        </Card>

        <SectionLabel>Payment method</SectionLabel>
        <Card selected style={styles.paymentCard}>
          <View style={styles.paymentRow}>
            <View style={styles.gcashIcon}>
              <Text style={styles.gcashText}>G</Text>
            </View>
            <View style={{ flex: 1 }}>
              <Text style={styles.paymentName}>GCash</Text>
              <Text style={styles.paymentLast4}>**** 8847</Text>
            </View>
            <View style={styles.checkCircle}>
              <Text style={styles.checkMark}>✓</Text>
            </View>
          </View>
        </Card>

        <InfoTip icon="💡">
          <Text style={styles.infoText}>
            A <Text style={{ fontWeight: '700' }}>₱200 hold</Text> will be placed on your GCash.
            Final amount charged after job completion.
          </Text>
        </InfoTip>
      </ScrollView>

      <View style={styles.footer}>
        <Button onPress={handleBook} fullWidth>
          🚛 Book Now — ₱1,850
        </Button>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  scrollContent: { padding: spacing[4], gap: spacing[3] },
  priceCard: { marginBottom: spacing[2] },
  totalSection: { alignItems: 'center', marginBottom: spacing[3] },
  estimateLabel: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
  totalPrice: {
    fontFamily: fontFamily.bold,
    fontSize: 36,
    fontWeight: '800',
    color: colors.orange,
    marginTop: spacing[1],
  },
  regulation: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey, marginTop: 2 },
  divider: { height: 1, backgroundColor: colors.greyLight, marginVertical: 10 },
  lineItem: { flexDirection: 'row', justifyContent: 'space-between', paddingVertical: 6 },
  lineLabel: { fontFamily: fontFamily.regular, fontSize: 11, color: colors.grey },
  lineValue: { fontFamily: fontFamily.semiBold, fontSize: 11, fontWeight: '600', color: colors.navy },
  paymentCard: { marginBottom: spacing[2] },
  paymentRow: { flexDirection: 'row', alignItems: 'center', gap: 10 },
  gcashIcon: {
    width: 36,
    height: 36,
    borderRadius: 10,
    backgroundColor: '#00B4D8',
    alignItems: 'center',
    justifyContent: 'center',
  },
  gcashText: {
    fontFamily: fontFamily.bold,
    fontSize: 12,
    fontWeight: '800',
    color: colors.white,
  },
  paymentName: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700', color: colors.navy },
  paymentLast4: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
  checkCircle: {
    width: 18,
    height: 18,
    borderRadius: 9,
    backgroundColor: colors.orange,
    alignItems: 'center',
    justifyContent: 'center',
  },
  checkMark: { color: colors.white, fontSize: 10 },
  infoText: { fontFamily: fontFamily.regular, fontSize: 10, color: '#5D6D7E' },
  footer: { padding: spacing[4], paddingBottom: spacing[1] },
});
