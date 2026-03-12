import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card, Button } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily, textStyles } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

const ALERTS = [
  {
    icon: '🌊',
    title: 'Flood Level: Knee-Deep',
    subtitle: 'EDSA-Guadalupe area • Updated 5 min ago',
    style: { backgroundColor: 'rgba(255,255,255,0.06)', borderColor: 'rgba(255,255,255,0.1)' },
    titleColor: colors.white,
  },
  {
    icon: '🚛',
    title: '2 Trucks Available',
    subtitle: 'Heavy-duty flatbeds only during typhoon',
    style: { backgroundColor: 'rgba(255,255,255,0.06)', borderColor: 'rgba(255,255,255,0.1)' },
    titleColor: colors.white,
  },
  {
    icon: '💰',
    title: 'Surge: 1.5× Base Rate',
    subtitle: 'MMDA-regulated typhoon rate cap applies',
    style: { backgroundColor: 'rgba(255,71,87,0.1)', borderColor: 'rgba(255,71,87,0.2)' },
    titleColor: colors.coral,
  },
  {
    icon: '🛡️',
    title: 'Safety Guaranteed',
    subtitle: 'Full insurance coverage during calamity',
    style: { backgroundColor: 'rgba(0,196,140,0.1)', borderColor: 'rgba(0,196,140,0.2)' },
    titleColor: colors.green,
  },
];

export default function TyphoonScreen() {
  return (
    <View style={styles.container}>
      <SafeAreaView style={{ flex: 1 }}>
        <View style={styles.header}>
          <View style={styles.alertBadge}>
            <Text style={styles.alertBadgeText}>⚠️ TYPHOON ALERT</Text>
          </View>
          <Text style={styles.title}>Typhoon Mode Active</Text>
          <Text style={styles.subtitle}>Signal #3 — Surge pricing may apply</Text>
        </View>

        <ScrollView contentContainerStyle={styles.scrollContent}>
          {ALERTS.map((alert) => (
            <Card
              key={alert.title}
              style={[styles.alertCard, { backgroundColor: alert.style.backgroundColor, borderColor: alert.style.borderColor }]}
            >
              <View style={styles.alertRow}>
                <Text style={{ fontSize: 22 }}>{alert.icon}</Text>
                <View style={{ flex: 1 }}>
                  <Text style={[styles.alertTitle, { color: alert.titleColor }]}>{alert.title}</Text>
                  <Text style={styles.alertSub}>{alert.subtitle}</Text>
                </View>
              </View>
            </Card>
          ))}
        </ScrollView>

        <View style={styles.footer}>
          <Button onPress={() => router.push('/booking/diagnose')} fullWidth>
            🚛 Book Emergency Tow
          </Button>
          <Button variant="ghost" onPress={() => router.replace('/(tabs)')} fullWidth>
            ← Back to Home
          </Button>
        </View>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#0B1D33' },
  header: { paddingHorizontal: spacing[5], paddingTop: spacing[2] },
  alertBadge: {
    backgroundColor: 'rgba(255,71,87,0.2)',
    borderRadius: 8,
    paddingHorizontal: 10,
    paddingVertical: 4,
    alignSelf: 'flex-start',
    marginBottom: spacing[3],
  },
  alertBadgeText: { fontFamily: fontFamily.bold, fontSize: 10, fontWeight: '700', color: colors.coral },
  title: { fontFamily: fontFamily.bold, fontSize: 20, fontWeight: '800', color: colors.white, marginBottom: spacing[1] },
  subtitle: { fontFamily: fontFamily.regular, fontSize: 11, color: 'rgba(255,255,255,0.5)', marginBottom: spacing[4] },
  scrollContent: { paddingHorizontal: spacing[4], gap: spacing[2] },
  alertCard: { borderWidth: 1 },
  alertRow: { flexDirection: 'row', alignItems: 'center', gap: 10 },
  alertTitle: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700' },
  alertSub: { fontFamily: fontFamily.regular, fontSize: 10, color: 'rgba(255,255,255,0.4)' },
  footer: { padding: spacing[4], gap: spacing[2] },
});
