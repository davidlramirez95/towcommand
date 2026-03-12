import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card, Button, Avatar, StatusPill, InfoTip } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';
import { useBookingStore } from '@/stores/booking';

export default function MatchedScreen() {
  const { matchedProvider, otp } = useBookingStore();

  const provider = matchedProvider ?? {
    id: 'prov-001',
    name: 'Juan Reyes',
    rating: 4.9,
    jobCount: 847,
    plate: 'ABC 1234',
    eta: 8,
    verified: true,
  };

  const displayOtp = otp ?? '482917';

  return (
    <SafeAreaView style={styles.container}>
      {/* Map placeholder */}
      <View style={styles.mapPlaceholder}>
        <View style={styles.mapOverlay}>
          <View style={styles.statusDot} />
          <Text style={styles.mapTitle}>Driver is on the way</Text>
          <Text style={styles.mapEta}>ETA {provider.eta} min</Text>
        </View>
        <Text style={styles.mapTruck}>🚛</Text>
        <View style={styles.mapPin}>
          <View style={styles.mapPinInner} />
        </View>
      </View>

      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Driver card */}
        <Card elevated>
          <View style={styles.driverRow}>
            <Avatar name={provider.name} size={50} />
            <View style={{ flex: 1 }}>
              <View style={styles.nameRow}>
                <Text style={styles.driverName}>{provider.name}</Text>
                {provider.verified && <StatusPill label="Verified" variant="success" />}
              </View>
              <View style={styles.statsRow}>
                <Text style={styles.rating}>★ {provider.rating}</Text>
                <Text style={styles.stats}>
                  • {provider.jobCount} jobs • {provider.plate}
                </Text>
              </View>
            </View>
          </View>
          <View style={styles.actionRow}>
            <Button variant="teal" small fullWidth onPress={() => router.push('/booking/chat')}>
              💬 Message
            </Button>
            <Button variant="secondary" small fullWidth onPress={() => router.push('/booking/tracking')}>
              📍 Track
            </Button>
          </View>
        </Card>

        {/* OTP Card */}
        <Card style={styles.otpCard}>
          <View style={styles.otpHeader}>
            <Text style={{ fontSize: 20 }}>🔐</Text>
            <View style={{ flex: 1 }}>
              <Text style={styles.otpTitle}>Digital Padala OTP</Text>
              <Text style={styles.otpSub}>Share this code when driver arrives</Text>
            </View>
          </View>
          <View style={styles.otpDigits}>
            {displayOtp.split('').map((digit, i) => (
              <View key={i} style={styles.otpBox}>
                <Text style={styles.otpDigit}>{digit}</Text>
              </View>
            ))}
          </View>
        </Card>

        <InfoTip icon="🔒">
          <Text style={styles.infoText}>
            OTP is end-to-end encrypted. Never share it before the driver arrives at your location.
          </Text>
        </InfoTip>
      </ScrollView>

      <View style={styles.footer}>
        <Button variant="danger" fullWidth small onPress={() => router.push('/sos')}>
          🆘 Emergency SOS
        </Button>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  // Map placeholder
  mapPlaceholder: {
    height: 155,
    backgroundColor: '#E8E4DE',
    position: 'relative',
    overflow: 'hidden',
  },
  mapOverlay: {
    position: 'absolute',
    top: 8,
    left: 8,
    right: 8,
    backgroundColor: 'rgba(11,29,51,0.85)',
    borderRadius: 12,
    padding: 10,
    paddingHorizontal: 14,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    zIndex: 1,
  },
  statusDot: { width: 8, height: 8, borderRadius: 4, backgroundColor: colors.green },
  mapTitle: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700', color: colors.white, flex: 1 },
  mapEta: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700', color: colors.gold },
  mapTruck: { position: 'absolute', top: '55%', left: '35%', fontSize: 20 },
  mapPin: {
    position: 'absolute',
    top: '40%',
    left: '55%',
    width: 20,
    height: 20,
    borderRadius: 10,
    backgroundColor: colors.orange,
    borderWidth: 3,
    borderColor: colors.white,
  },
  mapPinInner: {
    width: 6,
    height: 6,
    borderRadius: 3,
    backgroundColor: colors.white,
    alignSelf: 'center',
    marginTop: 2,
  },
  scrollContent: { padding: spacing[4], gap: spacing[3] },
  // Driver card
  driverRow: { flexDirection: 'row', alignItems: 'center', gap: 12 },
  nameRow: { flexDirection: 'row', alignItems: 'center', gap: 6 },
  driverName: { fontFamily: fontFamily.bold, fontSize: 14, fontWeight: '700', color: colors.navy },
  statsRow: { flexDirection: 'row', alignItems: 'center', gap: 8, marginTop: 3 },
  rating: { fontFamily: fontFamily.regular, fontSize: 11, color: colors.gold },
  stats: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
  actionRow: { flexDirection: 'row', gap: 8, marginTop: 12 },
  // OTP
  otpCard: { backgroundColor: colors.cream, borderWidth: 1.5, borderColor: 'rgba(245,166,35,0.2)' },
  otpHeader: { flexDirection: 'row', alignItems: 'center', gap: 10 },
  otpTitle: { fontFamily: fontFamily.bold, fontSize: 11, fontWeight: '700', color: colors.navy },
  otpSub: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
  otpDigits: { flexDirection: 'row', gap: 6, marginTop: 10, justifyContent: 'center' },
  otpBox: {
    width: 36,
    height: 44,
    borderRadius: 10,
    backgroundColor: colors.white,
    borderWidth: 2,
    borderColor: colors.orange,
    alignItems: 'center',
    justifyContent: 'center',
  },
  otpDigit: {
    fontFamily: fontFamily.bold,
    fontSize: 22,
    fontWeight: '800',
    color: colors.orange,
  },
  infoText: { fontFamily: fontFamily.regular, fontSize: 10, color: '#5D6D7E' },
  footer: { padding: spacing[4], paddingBottom: spacing[1] },
});
