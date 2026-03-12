import { Text, StyleSheet, ViewStyle } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

interface SectionLabelProps {
  children: string;
  style?: ViewStyle;
}

export function SectionLabel({ children, style }: SectionLabelProps) {
  return <Text style={[styles.label, style]}>{children.toUpperCase()}</Text>;
}

const styles = StyleSheet.create({
  label: {
    fontFamily: fontFamily.semiBold,
    fontSize: 9,
    fontWeight: '700',
    color: colors.grey,
    letterSpacing: 1.5,
    marginBottom: spacing[2],
  },
});
