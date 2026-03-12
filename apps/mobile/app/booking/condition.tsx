import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { BackHeader, Card, InfoTip, Button } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

const PHOTO_ANGLES = ['Front', 'Rear', 'Left', 'Right', 'FL Tire', 'FR Tire', 'RL Tire', 'RR Tire'];
const COMPLETED_COUNT = 3;

export default function ConditionScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <BackHeader title="Pre-Tow Condition Report" onBack={() => router.back()} />
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Card elevated style={styles.heroCard}>
          <View style={styles.heroRow}>
            <Text style={{ fontSize: 20 }}>📸</Text>
            <View>
              <Text style={styles.heroTitle}>8 photos required before towing</Text>
              <Text style={styles.heroSub}>All photos are timestamped & GPS-tagged</Text>
            </View>
          </View>
        </Card>

        <View style={styles.photoGrid}>
          {PHOTO_ANGLES.map((label, i) => {
            const done = i < COMPLETED_COUNT;
            return (
              <View
                key={label}
                style={[styles.photoSlot, done ? styles.photoSlotDone : styles.photoSlotPending]}
              >
                <Text style={{ fontSize: done ? 14 : 18 }}>{done ? '✅' : '📷'}</Text>
                <Text style={[styles.photoLabel, done && styles.photoLabelDone]}>{label}</Text>
              </View>
            );
          })}
        </View>

        <Card style={styles.videoCard}>
          <Text style={{ fontSize: 26, textAlign: 'center' }}>🎥</Text>
          <Text style={styles.videoTitle}>Record 360° Walk-Around</Text>
          <Text style={styles.videoSub}>Min 30 seconds • GPS & timestamp embedded</Text>
        </Card>

        <InfoTip icon="🔒">
          <Text style={styles.infoText}>
            All evidence is <Text style={{ fontWeight: '700' }}>tamper-proof</Text> with SHA-256
            hashing and stored for 1 year.
          </Text>
        </InfoTip>
      </ScrollView>

      <View style={styles.footer}>
        <Button onPress={() => router.push('/booking/complete')} fullWidth>
          {`Submit Report (${COMPLETED_COUNT}/${PHOTO_ANGLES.length} photos)`}
        </Button>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  scrollContent: { padding: spacing[4], gap: spacing[3] },
  heroCard: { backgroundColor: colors.navy, borderWidth: 0 },
  heroRow: { flexDirection: 'row', alignItems: 'center', gap: 10 },
  heroTitle: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700', color: colors.white },
  heroSub: { fontFamily: fontFamily.regular, fontSize: 10, color: 'rgba(255,255,255,0.6)' },
  photoGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 6,
  },
  photoSlot: {
    width: '23.5%',
    aspectRatio: 1,
    borderRadius: 12,
    alignItems: 'center',
    justifyContent: 'center',
    gap: 2,
  },
  photoSlotDone: {
    backgroundColor: 'rgba(0,196,140,0.12)',
    borderWidth: 2,
    borderColor: colors.green,
  },
  photoSlotPending: {
    backgroundColor: colors.light,
    borderWidth: 2,
    borderStyle: 'dashed',
    borderColor: colors.greyLight,
  },
  photoLabel: { fontFamily: fontFamily.semiBold, fontSize: 8, fontWeight: '600', color: colors.grey },
  photoLabelDone: { color: colors.green },
  videoCard: {
    alignItems: 'center',
    padding: spacing[5],
    borderWidth: 2,
    borderStyle: 'dashed',
    borderColor: colors.greyLight,
  },
  videoTitle: {
    fontFamily: fontFamily.semiBold,
    fontSize: 12,
    fontWeight: '600',
    color: colors.navy,
    marginTop: spacing[1],
  },
  videoSub: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
  infoText: { fontFamily: fontFamily.regular, fontSize: 10, color: '#5D6D7E' },
  footer: { padding: spacing[4], paddingBottom: spacing[1] },
});
