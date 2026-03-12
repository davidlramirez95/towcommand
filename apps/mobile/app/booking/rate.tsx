import { useState } from 'react';
import { View, Text, StyleSheet, ScrollView, TouchableOpacity, TextInput } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { BackHeader, Avatar, InfoTip, Button } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';
import { useBookingStore } from '@/stores/booking';

const POSITIVE_TAGS = ['Fast response', 'Professional', 'Careful handling', 'Great communication', 'Clean truck', 'Fair price'];

export default function RateScreen() {
  const [stars, setStars] = useState(0);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const { matchedProvider, reset } = useBookingStore();

  const provider = matchedProvider ?? { name: 'Juan Reyes' };

  const toggleTag = (tag: string) =>
    setSelectedTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag],
    );

  const handleSubmit = () => {
    reset();
    router.replace('/(tabs)');
  };

  return (
    <SafeAreaView style={styles.container}>
      <BackHeader title="Rate Your Experience" onBack={() => router.back()} />
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <View style={styles.providerSection}>
          <Avatar name={provider.name} size={60} />
          <Text style={styles.providerName}>{provider.name}</Text>
          <Text style={styles.providerJob}>Flatbed Tow • TC-2026-00847</Text>
        </View>

        <View style={styles.starsRow}>
          {[1, 2, 3, 4, 5].map((i) => (
            <TouchableOpacity
              key={i}
              onPress={() => setStars(i)}
              accessibilityRole="radio"
              accessibilityLabel={`${i} star${i > 1 ? 's' : ''}`}
              accessibilityState={{ selected: i <= stars }}
            >
              <Text style={[styles.star, i <= stars ? styles.starActive : styles.starInactive]}>
                ⭐
              </Text>
            </TouchableOpacity>
          ))}
        </View>

        {stars > 0 && (
          <>
            <Text style={styles.tagHeader}>
              {stars >= 4 ? 'WHAT WAS GREAT?' : 'WHAT COULD IMPROVE?'}
            </Text>
            <View style={styles.tagsGrid}>
              {POSITIVE_TAGS.map((tag) => {
                const selected = selectedTags.includes(tag);
                return (
                  <TouchableOpacity
                    key={tag}
                    onPress={() => toggleTag(tag)}
                    style={[styles.tag, selected && styles.tagSelected]}
                  >
                    <Text style={[styles.tagText, selected && styles.tagTextSelected]}>
                      {tag}
                    </Text>
                  </TouchableOpacity>
                );
              })}
            </View>

            <TextInput
              style={styles.commentInput}
              placeholder="Add a comment (optional)..."
              placeholderTextColor={colors.grey}
              multiline
              accessibilityLabel="Comment"
            />

            <InfoTip icon="💡">
              <Text style={styles.infoText}>
                Your review helps build a trusted community of Suki providers in your area.
              </Text>
            </InfoTip>
          </>
        )}
      </ScrollView>

      <View style={styles.footer}>
        <Button onPress={handleSubmit} fullWidth disabled={stars === 0}>
          Submit Review →
        </Button>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  scrollContent: { padding: spacing[5], gap: spacing[4] },
  providerSection: { alignItems: 'center', gap: spacing[2] },
  providerName: { fontFamily: fontFamily.bold, fontSize: 16, fontWeight: '700', color: colors.navy },
  providerJob: { fontFamily: fontFamily.regular, fontSize: 11, color: colors.grey },
  starsRow: { flexDirection: 'row', justifyContent: 'center', gap: 8 },
  star: { fontSize: 36 },
  starActive: { opacity: 1 },
  starInactive: { opacity: 0.3 },
  tagHeader: {
    fontFamily: fontFamily.semiBold,
    fontSize: 9,
    fontWeight: '700',
    color: colors.grey,
    letterSpacing: 1.5,
    textAlign: 'center',
  },
  tagsGrid: { flexDirection: 'row', flexWrap: 'wrap', gap: 6, justifyContent: 'center' },
  tag: {
    paddingVertical: 8,
    paddingHorizontal: 14,
    borderRadius: 10,
    borderWidth: 1.5,
    borderColor: colors.greyLight,
    backgroundColor: colors.white,
  },
  tagSelected: {
    borderWidth: 2,
    borderColor: colors.orange,
    backgroundColor: colors.cream,
  },
  tagText: { fontFamily: fontFamily.regular, fontSize: 11, fontWeight: '500', color: colors.navy },
  tagTextSelected: { fontWeight: '700', color: colors.orange },
  commentInput: {
    padding: spacing[3],
    borderRadius: 12,
    borderWidth: 1.5,
    borderColor: colors.greyLight,
    backgroundColor: colors.white,
    fontFamily: fontFamily.regular,
    fontSize: 12,
    color: colors.navy,
    minHeight: 60,
    textAlignVertical: 'top',
  },
  infoText: { fontFamily: fontFamily.regular, fontSize: 10, color: '#5D6D7E' },
  footer: { padding: spacing[5], paddingBottom: spacing[1] },
});
