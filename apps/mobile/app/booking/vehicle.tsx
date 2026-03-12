import { View, Text, StyleSheet, ScrollView, TouchableOpacity } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { BackHeader, Card, SectionLabel, Button } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';
import { useVehicleStore } from '@/stores/vehicle';

const CONDITIONS = [
  "Engine won't start",
  'Flat tire(s)',
  'Accident damage',
  'Other / Not sure',
];

const DEMO_VEHICLE = {
  id: 'demo-v1',
  make: 'Mitsubishi',
  model: 'Montero GLS',
  year: 2026,
  plate: 'ABC 1234',
  type: 'suv' as const,
  color: 'White',
};

export default function VehicleScreen() {
  const {
    vehicles,
    selectedVehicleId,
    selectedCondition,
    selectVehicle,
    selectCondition,
    addVehicle,
  } = useVehicleStore();

  // Add demo vehicle if none exist
  if (vehicles.length === 0) {
    addVehicle(DEMO_VEHICLE);
  }

  const canContinue = selectedVehicleId !== null && selectedCondition !== null;

  return (
    <SafeAreaView style={styles.container}>
      <BackHeader title="Vehicle Details" onBack={() => router.back()} />
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <SectionLabel>Saved vehicles</SectionLabel>
        {vehicles.map((v) => {
          const isSelected = v.id === selectedVehicleId;
          return (
            <Card
              key={v.id}
              onPress={() => selectVehicle(v.id)}
              selected={isSelected}
              style={styles.vehicleCard}
            >
              <View style={styles.vehicleRow}>
                <View style={styles.vehicleIcon}>
                  <Text style={{ fontSize: 22 }}>🚗</Text>
                </View>
                <View style={{ flex: 1 }}>
                  <Text style={styles.vehicleName}>
                    {v.year} {v.model}
                  </Text>
                  <Text style={styles.vehicleInfo}>
                    {v.plate} • {v.type.charAt(0).toUpperCase() + v.type.slice(1)} • {v.color}
                  </Text>
                </View>
                {isSelected && (
                  <View style={styles.checkCircle}>
                    <Text style={styles.checkMark}>✓</Text>
                  </View>
                )}
              </View>
            </Card>
          );
        })}

        <Card style={styles.addVehicleCard}>
          <Text style={{ fontSize: 22, textAlign: 'center' }}>➕</Text>
          <Text style={styles.addVehicleText}>Add New Vehicle</Text>
        </Card>

        <SectionLabel style={{ marginTop: spacing[3] }}>Vehicle condition</SectionLabel>
        {CONDITIONS.map((condition) => {
          const isSelected = selectedCondition === condition;
          return (
            <TouchableOpacity key={condition} onPress={() => selectCondition(condition)}>
              <Card
                selected={isSelected}
                style={styles.conditionCard}
              >
                <View style={styles.conditionRow}>
                  <Text
                    style={[styles.conditionText, isSelected && styles.conditionTextSelected]}
                  >
                    {condition}
                  </Text>
                  {isSelected && (
                    <View style={styles.checkCircleSmall}>
                      <Text style={styles.checkMarkSmall}>✓</Text>
                    </View>
                  )}
                </View>
              </Card>
            </TouchableOpacity>
          );
        })}
      </ScrollView>

      <View style={styles.footer}>
        <Button
          onPress={() => router.push('/booking/dropoff')}
          fullWidth
          disabled={!canContinue}
        >
          Continue →
        </Button>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  scrollContent: { padding: spacing[4], gap: spacing[2] },
  vehicleCard: { marginBottom: spacing[2] },
  vehicleRow: { flexDirection: 'row', alignItems: 'center', gap: 12 },
  vehicleIcon: {
    width: 46,
    height: 46,
    borderRadius: 14,
    backgroundColor: colors.cream,
    alignItems: 'center',
    justifyContent: 'center',
  },
  vehicleName: { fontFamily: fontFamily.bold, fontSize: 13, fontWeight: '700', color: colors.navy },
  vehicleInfo: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.grey },
  checkCircle: {
    width: 20,
    height: 20,
    borderRadius: 10,
    backgroundColor: colors.orange,
    alignItems: 'center',
    justifyContent: 'center',
  },
  checkMark: { color: colors.white, fontSize: 12, fontWeight: '800' },
  addVehicleCard: {
    borderWidth: 1.5,
    borderStyle: 'dashed',
    borderColor: colors.greyLight,
    alignItems: 'center',
    padding: spacing[4],
    marginBottom: spacing[3],
  },
  addVehicleText: {
    fontFamily: fontFamily.semiBold,
    fontSize: 12,
    fontWeight: '600',
    color: colors.navy,
    marginTop: spacing[1],
  },
  conditionCard: { marginBottom: spacing[1], padding: spacing[3] },
  conditionRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between' },
  conditionText: { fontFamily: fontFamily.regular, fontSize: 12, fontWeight: '500', color: colors.navy },
  conditionTextSelected: { fontWeight: '700', color: colors.orange },
  checkCircleSmall: {
    width: 18,
    height: 18,
    borderRadius: 9,
    backgroundColor: colors.orange,
    alignItems: 'center',
    justifyContent: 'center',
  },
  checkMarkSmall: { color: colors.white, fontSize: 10 },
  footer: { padding: spacing[4], paddingBottom: spacing[1] },
});
