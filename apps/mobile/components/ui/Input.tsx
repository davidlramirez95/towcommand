import { useState } from 'react';
import { TextInput, View, Text, StyleSheet, TextInputProps } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { fontFamily, fontSize, textStyles } from '@/lib/theme/typography';
import { borderRadius, spacing } from '@/lib/theme/spacing';

interface InputProps extends TextInputProps {
  label?: string;
  error?: string;
  helper?: string;
}

export function Input({ label, error, helper, style, ...props }: InputProps) {
  const [focused, setFocused] = useState(false);

  return (
    <View style={styles.wrapper}>
      {label && <Text style={styles.label}>{label}</Text>}
      <TextInput
        style={[
          styles.input,
          focused && styles.focused,
          error && styles.error,
          style,
        ]}
        placeholderTextColor={colors.grey}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
        accessibilityLabel={label}
        {...props}
      />
      {error && <Text style={styles.errorText}>{error}</Text>}
      {helper && !error && <Text style={styles.helperText}>{helper}</Text>}
    </View>
  );
}

const styles = StyleSheet.create({
  wrapper: {
    gap: spacing[1],
  },
  label: {
    ...textStyles.label,
    color: colors.navy,
  },
  input: {
    fontFamily: fontFamily.regular,
    fontSize: fontSize.md,
    color: colors.navy,
    backgroundColor: colors.white,
    borderWidth: 1.5,
    borderColor: colors.greyLight,
    borderRadius: borderRadius.md,
    paddingVertical: spacing[3],
    paddingHorizontal: spacing[4],
  },
  focused: {
    borderColor: colors.teal,
  },
  error: {
    borderColor: colors.coral,
  },
  errorText: {
    ...textStyles.caption,
    color: colors.coral,
  },
  helperText: {
    ...textStyles.caption,
    color: colors.grey,
  },
});
