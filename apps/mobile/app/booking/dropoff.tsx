import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { BackHeader, Card, SectionLabel, Button } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

const RECENT_LOCATIONS = [
  'Casa — BF Homes, Parañaque',
  'Mitsubishi Ortigas',
  'AutoHub SLEX, Alabang',
];

export default function DropoffScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <BackHeader title="Drop-off Location" onBack={() => router.back()} />
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <SectionLabel>Pickup location</SectionLabel>
        <Card style={styles.locationCard}>
          <View style={styles.locationRow}>
            <View style={[styles.dot, { backgroundColor: colors.orange }]} />
            <View>
              <Text style={styles.locationName}>Current Location</Text>
              <Text style={styles.locationAddr}>EDSA cor. Ayala Ave, Makati City</Text>
            </View>
          </View>
        </Card>

        <SectionLabel style={{ marginTop: spacing[3] }}>Drop-off location</SectionLabel>
        <Card selected style={styles.locationCard}>
          <View style={styles.locationRow}>
            <View style={[styles.dot, { backgroundColor: colors.teal }]} />
            <View>
              <Text style={styles.locationName}>Toyota Shaw, Mandaluyong</Text>
              <Text style={styles.locationAddr}>EDSA cor. Shaw Blvd • 4.2 km</Text>
            </View>
          </View>
        </Card>

        <SectionLabel style={{ marginTop: spacing[3] }}>Recent locations</SectionLabel>
        {RECENT_LOCATIONS.map((loc) => (
          <Card key={loc} style={styles.recentCard}>
            <View style={styles.locationRow}>
              <Text style={{ fontSize: 14 }}>📍</Text>
              <Text style={styles.recentText}>{loc}</Text>
            </View>
          </Card>
        ))}
      </ScrollView>

      <View style={styles.footer}>
        <Button onPress={() => router.push('/booking/price')} fullWidth>
          Confirm Route →
        </Button>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  scrollContent: { padding: spacing[4], gap: spacing[2] },
  locationCard: { marginBottom: spacing[2] },
  locationRow: { flexDirection: 'row', alignItems: 'center', gap: 10 },
  dot: { width: 10, height: 10, borderRadius: 5 },
  locationName: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700', color: colors.navy },
  locationAddr: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
  recentCard: { marginBottom: spacing[1], padding: spacing[3] },
  recentText: { fontFamily: fontFamily.regular, fontSize: 11, color: colors.navy },
  footer: { padding: spacing[4], paddingBottom: spacing[1] },
});
