import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card } from '@/components/ui/Card';
import { colors } from '@/lib/theme/colors';
import { textStyles, fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

const SERVICE_TYPES = [
  { id: 'FLATBED_TOWING', label: 'Flatbed Towing', desc: 'Best for sedans, SUVs, luxury vehicles', icon: '🚛' },
  { id: 'WHEEL_LIFT_TOWING', label: 'Wheel Lift Towing', desc: 'Quick hookup for short distances', icon: '🔗' },
  { id: 'MOTORCYCLE_TOWING', label: 'Motorcycle Towing', desc: 'Safe transport for 2-wheelers', icon: '🏍️' },
  { id: 'JUMPSTART', label: 'Jumpstart', desc: 'Dead battery? We\'ll get you running', icon: '🔋' },
  { id: 'TIRE_CHANGE', label: 'Tire Change', desc: 'Flat tire replacement on-site', icon: '🔧' },
  { id: 'LOCKOUT', label: 'Lockout Service', desc: 'Locked out of your vehicle?', icon: '🔑' },
  { id: 'FUEL_DELIVERY', label: 'Fuel Delivery', desc: 'Ran out of gas? We deliver', icon: '⛽' },
  { id: 'WINCH_RECOVERY', label: 'Winch Recovery', desc: 'Stuck vehicle extraction', icon: '⚙️' },
] as const;

export default function ServiceSelectionScreen() {
  const handleSelect = (serviceType: string) => {
    router.push({ pathname: '/booking/vehicle', params: { service: serviceType } });
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.backButton} onPress={() => router.back()}>←</Text>
        <Text style={styles.title}>Select Service</Text>
      </View>

      <ScrollView contentContainerStyle={styles.list}>
        {SERVICE_TYPES.map((service) => (
          <Card
            key={service.id}
            onPress={() => handleSelect(service.id)}
            style={styles.card}
          >
            <View style={styles.cardContent}>
              <Text style={styles.icon}>{service.icon}</Text>
              <View style={styles.cardText}>
                <Text style={styles.cardTitle}>{service.label}</Text>
                <Text style={styles.cardDesc}>{service.desc}</Text>
              </View>
              <Text style={styles.arrow}>›</Text>
            </View>
          </Card>
        ))}
      </ScrollView>
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
    alignItems: 'center',
    paddingHorizontal: spacing[4],
    paddingVertical: spacing[3],
    gap: spacing[3],
  },
  backButton: {
    fontSize: 20,
    color: colors.navy,
    padding: spacing[2],
  },
  title: {
    ...textStyles.h3,
    color: colors.navy,
  },
  list: {
    paddingHorizontal: spacing[4],
    paddingBottom: spacing[4],
    gap: spacing[3],
  },
  card: {
    padding: spacing[4],
  },
  cardContent: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[3],
  },
  icon: {
    fontSize: 28,
  },
  cardText: {
    flex: 1,
  },
  cardTitle: {
    ...textStyles.label,
    color: colors.navy,
  },
  cardDesc: {
    ...textStyles.caption,
    color: colors.grey,
    marginTop: 2,
  },
  arrow: {
    fontSize: 24,
    color: colors.grey,
  },
});
