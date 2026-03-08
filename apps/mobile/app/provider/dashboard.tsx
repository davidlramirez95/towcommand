import { View, Text, StyleSheet } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { colors } from '@/lib/theme/colors';
import { textStyles, fontFamily } from '@/lib/theme/typography';
import { spacing, borderRadius } from '@/lib/theme/spacing';

export default function ProviderDashboardScreen() {
  // Placeholder - will wire up to WebSocket for real-time job notifications
  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Provider Dashboard</Text>
        <Badge variant="success">Online</Badge>
      </View>

      <Card elevated style={styles.statsCard}>
        <Text style={styles.statsTitle}>Today's Summary</Text>
        <View style={styles.statsRow}>
          <StatItem label="Jobs" value="0" />
          <StatItem label="Earnings" value="PHP 0" />
          <StatItem label="Rating" value="--" />
        </View>
      </Card>

      <View style={styles.emptyState}>
        <Text style={styles.emptyIcon}>🔔</Text>
        <Text style={styles.emptyTitle}>Waiting for jobs</Text>
        <Text style={styles.emptySubtitle}>New job requests will appear here</Text>
      </View>
    </SafeAreaView>
  );
}

function StatItem({ label, value }: { label: string; value: string }) {
  return (
    <View style={styles.statItem}>
      <Text style={styles.statValue}>{value}</Text>
      <Text style={styles.statLabel}>{label}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing[5],
    paddingVertical: spacing[3],
  },
  title: {
    ...textStyles.h2,
    color: colors.navy,
  },
  statsCard: {
    marginHorizontal: spacing[4],
    marginBottom: spacing[4],
  },
  statsTitle: {
    ...textStyles.label,
    color: colors.grey,
    marginBottom: spacing[3],
  },
  statsRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
  },
  statItem: {
    alignItems: 'center',
  },
  statValue: {
    fontFamily: fontFamily.bold,
    fontSize: 20,
    color: colors.navy,
  },
  statLabel: {
    ...textStyles.caption,
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
