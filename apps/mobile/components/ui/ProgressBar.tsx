import { View, Text, StyleSheet } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

interface ProgressBarProps {
  steps: string[];
  currentStep: number;
}

export function ProgressBar({ steps, currentStep }: ProgressBarProps) {
  return (
    <View style={styles.container}>
      <View style={styles.barRow}>
        {steps.map((_, index) => (
          <View
            key={index}
            style={[styles.bar, index < currentStep ? styles.barActive : styles.barInactive]}
          />
        ))}
      </View>
      <View style={styles.labelRow}>
        {steps.map((label, index) => (
          <Text
            key={index}
            style={[styles.label, index < currentStep && styles.labelActive]}
            numberOfLines={1}
          >
            {label}
          </Text>
        ))}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    gap: spacing[1],
  },
  barRow: {
    flexDirection: 'row',
    gap: 4,
  },
  bar: {
    flex: 1,
    height: 3,
    borderRadius: 2,
  },
  barActive: {
    backgroundColor: colors.green,
  },
  barInactive: {
    backgroundColor: 'rgba(0,0,0,0.1)',
  },
  labelRow: {
    flexDirection: 'row',
    gap: 4,
  },
  label: {
    flex: 1,
    fontFamily: fontFamily.regular,
    fontSize: 9,
    color: colors.grey,
    textAlign: 'center',
  },
  labelActive: {
    color: colors.green,
  },
});
