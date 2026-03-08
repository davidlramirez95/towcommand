import { Pressable, StyleSheet, View, ViewProps } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { borderRadius, spacing } from '@/lib/theme/spacing';

interface CardProps extends ViewProps {
  children: React.ReactNode;
  onPress?: () => void;
  elevated?: boolean;
  selected?: boolean;
}

export function Card({ children, onPress, elevated = false, selected = false, style, ...props }: CardProps) {
  const content = (
    <View
      style={[
        styles.container,
        elevated && styles.elevated,
        selected && styles.selected,
        style,
      ]}
      {...props}
    >
      {children}
    </View>
  );

  if (onPress) {
    return (
      <Pressable onPress={onPress} style={({ pressed }) => pressed && styles.pressed}>
        {content}
      </Pressable>
    );
  }

  return content;
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.white,
    borderRadius: borderRadius.lg,
    padding: spacing[4],
    borderWidth: 1,
    borderColor: colors.greyLight,
  },
  elevated: {
    shadowColor: colors.navy,
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.08,
    shadowRadius: 30,
    elevation: 4,
  },
  selected: {
    borderWidth: 2,
    borderColor: colors.orange,
    backgroundColor: colors.cream,
  },
  pressed: {
    opacity: 0.9,
  },
});
