import { useEffect } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { router } from 'expo-router';
import { Button } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';
import { useBookingStore } from '@/stores/booking';

export default function MatchingScreen() {
  const { setMatchingState, setMatchedProvider, setOtp } = useBookingStore();

  useEffect(() => {
    const timer = setTimeout(() => {
      setMatchingState('found');
      setMatchedProvider({
        id: 'prov-001',
        name: 'Juan Reyes',
        rating: 4.9,
        jobCount: 847,
        plate: 'ABC 1234',
        eta: 8,
        verified: true,
      });
      setOtp('482917');
      router.replace('/booking/matched');
    }, 3000);
    return () => clearTimeout(timer);
  }, [setMatchingState, setMatchedProvider, setOtp]);

  const handleCancel = () => {
    setMatchingState('idle');
    router.replace('/(tabs)');
  };

  return (
    <View style={styles.container}>
      <View style={styles.circleOuter}>
        <View style={styles.circleInner}>
          <Text style={styles.emoji}>🚛</Text>
        </View>
      </View>

      <Text style={styles.title}>Finding nearby trucks...</Text>
      <Text style={styles.subtitle}>Matching you with the best provider</Text>

      <View style={styles.tags}>
        {['📍 Makati', '🚛 3 available', '⏱ ~2 min'].map((tag) => (
          <View key={tag} style={styles.tag}>
            <Text style={styles.tagText}>{tag}</Text>
          </View>
        ))}
      </View>

      <View style={styles.cancelContainer}>
        <Button variant="ghost" onPress={handleCancel} small>
          Cancel Search
        </Button>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.navy,
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing[4],
  },
  circleOuter: {
    width: 100,
    height: 100,
    borderRadius: 50,
    borderWidth: 3,
    borderColor: 'rgba(245,166,35,0.15)',
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: spacing[5],
  },
  circleInner: {
    width: 64,
    height: 64,
    borderRadius: 32,
    backgroundColor: colors.orange,
    alignItems: 'center',
    justifyContent: 'center',
  },
  emoji: { fontSize: 38 },
  title: {
    fontFamily: fontFamily.bold,
    fontSize: 18,
    fontWeight: '700',
    color: colors.white,
  },
  subtitle: {
    fontFamily: fontFamily.regular,
    fontSize: 12,
    color: 'rgba(255,255,255,0.5)',
  },
  tags: { flexDirection: 'row', gap: 10, marginTop: spacing[4] },
  tag: {
    backgroundColor: 'rgba(255,255,255,0.08)',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 8,
  },
  tagText: {
    fontFamily: fontFamily.regular,
    fontSize: 10,
    color: 'rgba(255,255,255,0.4)',
  },
  cancelContainer: { position: 'absolute', bottom: 60 },
});
