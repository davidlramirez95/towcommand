import { Pressable, StyleSheet, Text, ActivityIndicator, ViewStyle, TextStyle } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { textStyles } from '@/lib/theme/typography';
import { borderRadius, spacing } from '@/lib/theme/spacing';

type ButtonVariant = 'primary' | 'secondary' | 'teal' | 'danger' | 'ghost';

interface ButtonProps {
  children: string;
  onPress: () => void;
  variant?: ButtonVariant;
  fullWidth?: boolean;
  small?: boolean;
  loading?: boolean;
  disabled?: boolean;
  accessibilityLabel?: string;
}

const variantStyles: Record<ButtonVariant, { container: ViewStyle; text: TextStyle }> = {
  primary: {
    container: { backgroundColor: colors.orange },
    text: { color: colors.white },
  },
  secondary: {
    container: { backgroundColor: colors.light, borderWidth: 1.5, borderColor: colors.greyLight },
    text: { color: colors.navy },
  },
  teal: {
    container: { backgroundColor: colors.teal },
    text: { color: colors.white },
  },
  danger: {
    container: { backgroundColor: colors.coral },
    text: { color: colors.white },
  },
  ghost: {
    container: { backgroundColor: 'transparent' },
    text: { color: colors.orange },
  },
};

export function Button({
  children,
  onPress,
  variant = 'primary',
  fullWidth = false,
  small = false,
  loading = false,
  disabled = false,
  accessibilityLabel,
}: ButtonProps) {
  const styles = variantStyles[variant];

  return (
    <Pressable
      onPress={onPress}
      disabled={disabled || loading}
      accessibilityLabel={accessibilityLabel ?? children}
      accessibilityRole="button"
      style={({ pressed }) => [
        baseStyles.container,
        styles.container,
        small && baseStyles.small,
        fullWidth && baseStyles.fullWidth,
        (disabled || loading) && baseStyles.disabled,
        pressed && baseStyles.pressed,
      ]}
    >
      {loading ? (
        <ActivityIndicator size="small" color={styles.text.color} />
      ) : (
        <Text style={[small ? textStyles.buttonSmall : textStyles.button, styles.text]}>
          {children}
        </Text>
      )}
    </Pressable>
  );
}

const baseStyles = StyleSheet.create({
  container: {
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: borderRadius.md,
    paddingVertical: spacing[3],
    paddingHorizontal: spacing[6],
    flexDirection: 'row',
    gap: spacing[2],
  },
  small: {
    borderRadius: borderRadius.sm,
    paddingVertical: spacing[2],
    paddingHorizontal: spacing[4],
  },
  fullWidth: {
    width: '100%',
  },
  disabled: {
    opacity: 0.5,
  },
  pressed: {
    opacity: 0.8,
  },
});
