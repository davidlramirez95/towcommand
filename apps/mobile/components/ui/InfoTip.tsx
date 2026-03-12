import { View, Text, StyleSheet } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing, borderRadius } from '@/lib/theme/spacing';

interface InfoTipProps {
  icon: string;
  children: React.ReactNode;
}

export function InfoTip({ icon, children }: InfoTipProps) {
  return (
    <View style={styles.container}>
      <Text style={styles.icon}>{icon}</Text>
      <View style={styles.content}>
        {typeof children === 'string' ? (
          <Text style={styles.text}>{children}</Text>
        ) : (
          children
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[2],
    backgroundColor: colors.cream,
    borderWidth: 1.5,
    borderColor: 'rgba(245,166,35,0.2)',
    borderRadius: borderRadius.lg,
    padding: spacing[4],
  },
  icon: {
    fontSize: 14,
  },
  content: {
    flex: 1,
  },
  text: {
    fontFamily: fontFamily.regular,
    fontSize: 10,
    color: '#5D6D7E',
  },
});
