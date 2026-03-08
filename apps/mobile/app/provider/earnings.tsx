import { View, Text, StyleSheet } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Card } from '@/components/ui/Card';
import { colors } from '@/lib/theme/colors';
import { textStyles, fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';

/**
 * Provider earnings summary screen.
 * Shows today, this week, this month, and all-time earnings.
 *
 * TODO (Week 6): Wire up to GET /providers/{id}/earnings API endpoint.
 */
export default function EarningsScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Earnings</Text>
      </View>

      <View style={styles.periods}>
        <EarningsPeriodCard
          label="Today"
          gross="PHP 0.00"
          net="PHP 0.00"
          jobs={0}
        />
        <EarningsPeriodCard
          label="This Week"
          gross="PHP 0.00"
          net="PHP 0.00"
          jobs={0}
        />
        <EarningsPeriodCard
          label="This Month"
          gross="PHP 0.00"
          net="PHP 0.00"
          jobs={0}
        />
        <EarningsPeriodCard
          label="All Time"
          gross="PHP 0.00"
          net="PHP 0.00"
          jobs={0}
        />
      </View>
    </SafeAreaView>
  );
}

function EarningsPeriodCard({
  label,
  gross,
  net,
  jobs,
}: {
  label: string;
  gross: string;
  net: string;
  jobs: number;
}) {
  return (
    <Card style={styles.periodCard}>
      <Text style={styles.periodLabel}>{label}</Text>
      <Text style={styles.periodNet}>{net}</Text>
      <View style={styles.periodRow}>
        <Text style={styles.periodGross}>Gross: {gross}</Text>
        <Text style={styles.periodJobs}>{jobs} jobs</Text>
      </View>
    </Card>
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
  periods: {
    paddingHorizontal: spacing[4],
    gap: spacing[3],
  },
  periodCard: {
    padding: spacing[4],
  },
  periodLabel: {
    ...textStyles.label,
    color: colors.grey,
    marginBottom: spacing[2],
  },
  periodNet: {
    fontFamily: fontFamily.bold,
    fontSize: 24,
    color: colors.navy,
    marginBottom: spacing[2],
  },
  periodRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  periodGross: {
    ...textStyles.caption,
    color: colors.grey,
  },
  periodJobs: {
    ...textStyles.caption,
    color: colors.teal,
  },
});
