import { View, Text, StyleSheet } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';

type PillVariant = 'success' | 'warning' | 'danger' | 'info' | 'neutral';

interface StatusPillProps {
  label: string;
  variant: PillVariant;
}

const variantColors: Record<PillVariant, string> = {
  success: colors.green,
  warning: colors.gold,
  danger: colors.coral,
  info: colors.blue,
  neutral: colors.grey,
};

export function StatusPill({ label, variant }: StatusPillProps) {
  return (
    <View style={[styles.container, { backgroundColor: variantColors[variant] }]}>
      <Text style={styles.label}>{label.toUpperCase()}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 4,
    alignSelf: 'flex-start',
  },
  label: {
    fontFamily: fontFamily.bold,
    fontSize: 8,
    fontWeight: '700',
    color: colors.white,
  },
});
