import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing, borderRadius } from '@/lib/theme/spacing';

interface BackHeaderProps {
  title: string;
  onBack: () => void;
  right?: React.ReactNode;
}

export function BackHeader({ title, onBack, right }: BackHeaderProps) {
  return (
    <View style={styles.container}>
      <TouchableOpacity
        onPress={onBack}
        style={styles.backButton}
        accessibilityRole="button"
        accessibilityLabel="Go back"
      >
        <Text style={styles.backArrow}>←</Text>
      </TouchableOpacity>
      <Text style={styles.title} numberOfLines={1}>
        {title}
      </Text>
      {right && <View style={styles.right}>{right}</View>}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing[4],
    paddingTop: spacing[1],
    paddingBottom: spacing[2],
    gap: 10,
  },
  backButton: {
    width: 36,
    height: 36,
    borderRadius: borderRadius.md,
    backgroundColor: colors.light,
    alignItems: 'center',
    justifyContent: 'center',
  },
  backArrow: {
    fontSize: 16,
    color: colors.navy,
  },
  title: {
    flex: 1,
    fontFamily: fontFamily.bold,
    fontSize: 16,
    fontWeight: '700',
    color: colors.navy,
  },
  right: {
    flexShrink: 0,
  },
});
