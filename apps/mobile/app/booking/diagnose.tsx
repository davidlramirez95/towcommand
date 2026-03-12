import { useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, TouchableOpacity } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { BackHeader, Card, InfoTip, SectionLabel, Button } from '@/components/ui';
import { StatusPill } from '@/components/ui/StatusPill';
import { colors } from '@/lib/theme/colors';
import { fontFamily, textStyles } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';
import { useDiagnosisStore } from '@/stores/diagnosis';

const SYMPTOMS = [
  { id: 'fuel', icon: '⛽', label: 'Empty fuel' },
  { id: 'tire', icon: '💨', label: 'Flat tire' },
  { id: 'battery', icon: '🔋', label: 'Dead battery' },
  { id: 'lockout', icon: '🔑', label: 'Locked out' },
  { id: 'overheat', icon: '🌡️', label: 'Overheating' },
  { id: 'accident', icon: '🚨', label: 'Accident' },
  { id: 'engine', icon: '⚙️', label: 'Engine trouble' },
  { id: 'leak', icon: '💧', label: 'Leaking fluid' },
  { id: 'electrical', icon: '🔌', label: 'Electrical' },
  { id: 'unknown', icon: '❓', label: 'Not sure' },
];

function getSymptomHint(symptoms: string[]): string | null {
  if (symptoms.includes('fuel')) return 'This looks like you just need fuel delivery!';
  if (symptoms.includes('tire')) return 'Flat tire? On-site mechanic can change this in 15 mins.';
  if (symptoms.includes('accident'))
    return "Accident detected. We'll prioritize emergency recovery.";
  if (symptoms.length > 0) return 'Analyzing your combination of symptoms...';
  return null;
}

// Analyzing screen
function AnalyzingScreen() {
  const { selectedSymptoms, setResults } = useDiagnosisStore();

  useEffect(() => {
    const timer = setTimeout(() => {
      setResults({
        needsTow: false,
        confidence: 92,
        recommendedService: 'On-Site Mechanic',
        estimatedSavings: 2500,
        alternatives: [
          { service: 'Fuel Delivery', priceRange: '₱200–₱500', eta: '20 min' },
          { service: 'Full Tow', priceRange: '₱1,800–₱3,500', eta: '12 min' },
        ],
      });
    }, 2000);
    return () => clearTimeout(timer);
  }, [setResults]);

  return (
    <View style={styles.analyzingContainer}>
      <View style={styles.analyzingCircleOuter}>
        <View style={styles.analyzingCircleInner}>
          <Text style={styles.analyzingEmoji}>🤖</Text>
        </View>
      </View>
      <Text style={styles.analyzingTitle}>Analyzing your symptoms...</Text>
      <Text style={styles.analyzingSub}>
        AI is matching {selectedSymptoms.length} symptom
        {selectedSymptoms.length > 1 ? 's' : ''} to the best service
      </Text>
    </View>
  );
}

// Results screen
function ResultsScreen() {
  const { diagnosisResult, reset } = useDiagnosisStore();

  if (!diagnosisResult) return null;

  return (
    <SafeAreaView style={styles.container}>
      <BackHeader
        title="AI Diagnosis Result"
        onBack={() => reset()}
        right={<StatusPill label="AI-Powered" variant="info" />}
      />
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Card elevated style={styles.resultCard}>
          <View style={styles.resultRow}>
            <View style={styles.resultIcon}>
              <Text style={{ fontSize: 24 }}>✅</Text>
            </View>
            <View style={{ flex: 1 }}>
              <Text style={styles.resultTitle}>
                {diagnosisResult.needsTow ? 'You need a tow' : "You DON'T need a tow!"}
              </Text>
              <Text style={styles.resultConfidence}>
                {diagnosisResult.confidence}% confidence • Save up to ₱
                {diagnosisResult.estimatedSavings.toLocaleString()}
              </Text>
            </View>
          </View>
        </Card>

        <SectionLabel>Recommended for you</SectionLabel>
        <Card
          elevated
          selected
          onPress={() => router.push('/booking/vehicle')}
          style={styles.recommendedCard}
        >
          <View style={styles.resultRow}>
            <View style={[styles.serviceIcon, { backgroundColor: '#FFF3EB' }]}>
              <Text style={{ fontSize: 22 }}>🔧</Text>
            </View>
            <View style={{ flex: 1 }}>
              <Text style={styles.serviceName}>{diagnosisResult.recommendedService}</Text>
              <Text style={styles.serviceDesc}>Tire change + basic repair</Text>
              <View style={styles.priceRow}>
                <Text style={styles.price}>₱500–₱1,200</Text>
                <Text style={styles.eta}>ETA 15 min</Text>
              </View>
            </View>
          </View>
        </Card>

        <SectionLabel style={{ marginTop: spacing[4] }}>Other options</SectionLabel>
        {diagnosisResult.alternatives.map((alt, i) => (
          <Card key={i} onPress={() => router.push('/booking/vehicle')} style={styles.altCard}>
            <View style={styles.resultRow}>
              <View style={[styles.serviceIcon, { backgroundColor: colors.light }]}>
                <Text style={{ fontSize: 18 }}>{i === 0 ? '⛽' : '🚛'}</Text>
              </View>
              <View style={{ flex: 1 }}>
                <Text style={styles.altName}>{alt.service}</Text>
              </View>
              <View style={{ alignItems: 'flex-end' }}>
                <Text style={styles.altPrice}>{alt.priceRange}</Text>
                <Text style={styles.altEta}>{alt.eta}</Text>
              </View>
            </View>
          </Card>
        ))}

        <InfoTip icon="💰">
          <Text style={styles.infoText}>
            Users who used AI Diagnosis saved avg <Text style={{ fontWeight: '700' }}>₱2,100</Text>{' '}
            per incident
          </Text>
        </InfoTip>

        <TouchableOpacity
          onPress={() => router.push('/booking/service')}
          style={styles.skipLink}
        >
          <Text style={styles.skipText}>Skip AI → Choose service manually</Text>
        </TouchableOpacity>
      </ScrollView>
    </SafeAreaView>
  );
}

// Main symptom selection screen
export default function DiagnoseScreen() {
  const { step, selectedSymptoms, toggleSymptom, startAnalysis } = useDiagnosisStore();

  if (step === 'analyzing') return <AnalyzingScreen />;
  if (step === 'results') return <ResultsScreen />;

  const hint = getSymptomHint(selectedSymptoms);
  const canAnalyze = selectedSymptoms.length > 0;

  return (
    <SafeAreaView style={styles.container}>
      <BackHeader
        title="Smart Diagnosis"
        onBack={() => router.back()}
        right={<StatusPill label="AI-Powered" variant="info" />}
      />
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Card elevated style={styles.heroCard}>
          <View style={styles.heroRow}>
            <Text style={{ fontSize: 22 }}>🤖</Text>
            <Text style={styles.heroTitle}>What's happening to your vehicle?</Text>
          </View>
          <Text style={styles.heroSub}>
            Select all symptoms so I can recommend the cheapest & fastest solution.
          </Text>
        </Card>

        <View style={styles.quickActions}>
          <Card style={styles.quickAction}>
            <Text style={{ fontSize: 22, textAlign: 'center' }}>🎤</Text>
            <Text style={styles.quickLabel}>Voice Describe</Text>
            <Text style={styles.quickSub}>Tagalog or English</Text>
          </Card>
          <Card style={styles.quickAction}>
            <Text style={{ fontSize: 22, textAlign: 'center' }}>📷</Text>
            <Text style={styles.quickLabel}>Take Photo</Text>
            <Text style={styles.quickSub}>AI visual analysis</Text>
          </Card>
        </View>

        <SectionLabel>Or tap your symptoms</SectionLabel>
        <View style={styles.symptomsGrid}>
          {SYMPTOMS.map((s) => {
            const selected = selectedSymptoms.includes(s.id);
            return (
              <TouchableOpacity
                key={s.id}
                onPress={() => toggleSymptom(s.id)}
                style={[styles.symptomChip, selected && styles.symptomChipSelected]}
                accessibilityRole="checkbox"
                accessibilityState={{ checked: selected }}
                accessibilityLabel={s.label}
              >
                <Text style={{ fontSize: 17 }}>{s.icon}</Text>
                <Text style={[styles.symptomLabel, selected && styles.symptomLabelSelected]}>
                  {s.label}
                </Text>
              </TouchableOpacity>
            );
          })}
        </View>

        {hint && (
          <InfoTip icon="💡">
            <Text style={styles.infoText}>{hint}</Text>
          </InfoTip>
        )}
      </ScrollView>

      <View style={styles.footer}>
        <Button
          onPress={() => canAnalyze && startAnalysis()}
          fullWidth
          disabled={!canAnalyze}
        >
          {`🤖 Diagnose My Problem (${selectedSymptoms.length})`}
        </Button>
        <TouchableOpacity
          onPress={() => router.push('/booking/service')}
          style={styles.skipLink}
        >
          <Text style={styles.skipText}>I already know what I need →</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  scrollContent: { padding: spacing[4], gap: spacing[3] },
  // Analyzing
  analyzingContainer: {
    flex: 1,
    backgroundColor: colors.navy,
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing[4],
  },
  analyzingCircleOuter: {
    width: 90,
    height: 90,
    borderRadius: 45,
    borderWidth: 3,
    borderColor: 'rgba(245,166,35,0.2)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  analyzingCircleInner: {
    width: 62,
    height: 62,
    borderRadius: 31,
    backgroundColor: colors.orange,
    alignItems: 'center',
    justifyContent: 'center',
  },
  analyzingEmoji: { fontSize: 34 },
  analyzingTitle: {
    fontFamily: fontFamily.bold,
    fontSize: 16,
    fontWeight: '700',
    color: colors.white,
  },
  analyzingSub: {
    fontFamily: fontFamily.regular,
    fontSize: 11,
    color: 'rgba(255,255,255,0.5)',
    textAlign: 'center',
    paddingHorizontal: 50,
  },
  // Hero card
  heroCard: {
    backgroundColor: colors.navy,
    borderWidth: 0,
  },
  heroRow: { flexDirection: 'row', alignItems: 'center', gap: 10, marginBottom: spacing[2] },
  heroTitle: {
    fontFamily: fontFamily.bold,
    fontSize: 13,
    fontWeight: '700',
    color: colors.white,
  },
  heroSub: {
    fontFamily: fontFamily.regular,
    fontSize: 11,
    color: 'rgba(255,255,255,0.6)',
    lineHeight: 16,
  },
  // Quick actions
  quickActions: { flexDirection: 'row', gap: spacing[2] },
  quickAction: { flex: 1, alignItems: 'center', padding: spacing[3] },
  quickLabel: {
    fontFamily: fontFamily.semiBold,
    fontSize: 10,
    fontWeight: '600',
    color: colors.navy,
    marginTop: spacing[1],
  },
  quickSub: { fontFamily: fontFamily.regular, fontSize: 9, color: colors.grey },
  // Symptoms
  symptomsGrid: { flexDirection: 'row', flexWrap: 'wrap', gap: 7 },
  symptomChip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 7,
    paddingVertical: 10,
    paddingHorizontal: 13,
    borderRadius: 12,
    borderWidth: 1.5,
    borderColor: colors.greyLight,
    backgroundColor: colors.white,
  },
  symptomChipSelected: {
    borderWidth: 2,
    borderColor: colors.orange,
    backgroundColor: colors.cream,
  },
  symptomLabel: {
    fontFamily: fontFamily.regular,
    fontSize: 11,
    fontWeight: '500',
    color: colors.navy,
  },
  symptomLabelSelected: {
    fontWeight: '700',
    color: colors.orange,
  },
  // Results
  resultCard: {
    backgroundColor: '#E8F8F0',
    borderWidth: 1.5,
    borderColor: 'rgba(0,196,140,0.3)',
  },
  resultRow: { flexDirection: 'row', alignItems: 'center', gap: 10 },
  resultIcon: {
    width: 44,
    height: 44,
    borderRadius: 14,
    backgroundColor: 'rgba(0,196,140,0.15)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  resultTitle: {
    fontFamily: fontFamily.bold,
    fontSize: 14,
    fontWeight: '700',
    color: '#1A7F5A',
  },
  resultConfidence: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.green },
  recommendedCard: { marginBottom: spacing[3] },
  serviceIcon: {
    width: 46,
    height: 46,
    borderRadius: 14,
    alignItems: 'center',
    justifyContent: 'center',
  },
  serviceName: { fontFamily: fontFamily.bold, fontSize: 14, fontWeight: '700', color: colors.navy },
  serviceDesc: {
    fontFamily: fontFamily.regular,
    fontSize: 11,
    color: colors.grey,
    marginTop: 2,
  },
  priceRow: { flexDirection: 'row', gap: 10, marginTop: 5 },
  price: { fontFamily: fontFamily.bold, fontSize: 12, fontWeight: '700', color: colors.green },
  eta: {
    fontFamily: fontFamily.regular,
    fontSize: 11,
    fontWeight: '600',
    color: colors.orange,
  },
  altCard: { marginBottom: spacing[2] },
  altName: {
    fontFamily: fontFamily.semiBold,
    fontSize: 12,
    fontWeight: '600',
    color: colors.navy,
  },
  altPrice: { fontFamily: fontFamily.bold, fontSize: 11, fontWeight: '700', color: colors.navy },
  altEta: { fontFamily: fontFamily.regular, fontSize: 9, color: colors.grey },
  infoText: { fontFamily: fontFamily.regular, fontSize: 10, color: '#5D6D7E' },
  // Footer
  footer: { padding: spacing[4], paddingBottom: spacing[1] },
  skipLink: { alignItems: 'center', marginTop: 10 },
  skipText: {
    fontFamily: fontFamily.regular,
    fontSize: 11,
    color: colors.grey,
    textDecorationLine: 'underline',
  },
});
