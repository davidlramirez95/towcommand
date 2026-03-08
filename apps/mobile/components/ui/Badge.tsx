import { View, Text, StyleSheet, ViewStyle } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { textStyles } from '@/lib/theme/typography';
import { borderRadius, spacing } from '@/lib/theme/spacing';

type BadgeVariant = 'default' | 'success' | 'warning' | 'error' | 'info' | 'premium';

interface BadgeProps {
  children: string;
  variant?: BadgeVariant;
}

const variantStyles: Record<BadgeVariant, ViewStyle & { textColor: string }> = {
  default: { backgroundColor: colors.light, textColor: colors.navy },
  success: { backgroundColor: '#E8F5E9', textColor: colors.green },
  warning: { backgroundColor: '#FFF3E0', textColor: colors.gold },
  error: { backgroundColor: '#FFEBEE', textColor: colors.coral },
  info: { backgroundColor: '#E3F2FD', textColor: colors.blue },
  premium: { backgroundColor: '#FFF8E1', textColor: colors.gold },
};

export function Badge({ children, variant = 'default' }: BadgeProps) {
  const style = variantStyles[variant];

  return (
    <View style={[styles.container, { backgroundColor: style.backgroundColor }]}>
      <Text style={[styles.text, { color: style.textColor }]}>{children}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    borderRadius: borderRadius.pill,
    paddingVertical: spacing[1],
    paddingHorizontal: spacing[2],
    alignSelf: 'flex-start',
  },
  text: {
    ...textStyles.caption,
    fontWeight: '700',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
});
