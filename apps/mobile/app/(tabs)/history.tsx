import { View, Text, StyleSheet, FlatList } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { colors } from '@/lib/theme/colors';
import { textStyles } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

export default function HistoryScreen() {
  // Placeholder - will be replaced with TanStack Query + API
  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Activity</Text>
        <Text style={styles.subtitle}>Your booking history</Text>
      </View>

      <View style={styles.emptyState}>
        <Text style={styles.emptyIcon}>📋</Text>
        <Text style={styles.emptyTitle}>No bookings yet</Text>
        <Text style={styles.emptySubtitle}>Your tow truck bookings will appear here</Text>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  header: {
    paddingHorizontal: spacing[5],
    paddingTop: spacing[3],
    paddingBottom: spacing[4],
  },
  title: {
    ...textStyles.h2,
    color: colors.navy,
  },
  subtitle: {
    ...textStyles.bodySmall,
    color: colors.grey,
    marginTop: 2,
  },
  emptyState: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingBottom: spacing[16],
  },
  emptyIcon: {
    fontSize: 48,
    marginBottom: spacing[4],
  },
  emptyTitle: {
    ...textStyles.h3,
    color: colors.navy,
    marginBottom: spacing[2],
  },
  emptySubtitle: {
    ...textStyles.bodySmall,
    color: colors.grey,
  },
});
