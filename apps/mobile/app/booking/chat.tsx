import { View, Text, StyleSheet, ScrollView, TextInput, TouchableOpacity } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Avatar } from '@/components/ui';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';
import { spacing } from '@/lib/theme/spacing';
import { useBookingStore } from '@/stores/booking';

const DEMO_MESSAGES = [
  { from: 'driver', msg: 'Magandang hapon po! On my way na. I\'ll be there in about 8 minutes 🚛' },
  { from: 'me', msg: 'Thank you! I\'m at the gas station beside McDonald\'s EDSA' },
  { from: 'driver', msg: 'Copy po! I can see the location. White Montero po ba?' },
  { from: 'me', msg: 'Yes correct! White Montero GLS' },
];

export default function ChatScreen() {
  const { matchedProvider } = useBookingStore();
  const provider = matchedProvider ?? { name: 'Juan Reyes', eta: 5 };

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity
          onPress={() => router.back()}
          style={styles.backBtn}
          accessibilityLabel="Go back"
        >
          <Text style={styles.backArrow}>←</Text>
        </TouchableOpacity>
        <Avatar name={provider.name} size={36} />
        <View style={{ flex: 1 }}>
          <Text style={styles.headerName}>{provider.name}</Text>
          <Text style={styles.headerStatus}>● Online • ETA {provider.eta} min</Text>
        </View>
        <Text style={{ fontSize: 16 }}>📞</Text>
      </View>

      {/* Messages */}
      <ScrollView style={styles.messages} contentContainerStyle={styles.messagesList}>
        <View style={styles.timestamp}>
          <Text style={styles.timestampText}>Today 2:34 PM</Text>
        </View>
        {DEMO_MESSAGES.map((m, i) => (
          <View
            key={i}
            style={[styles.bubble, m.from === 'me' ? styles.bubbleMe : styles.bubbleDriver]}
          >
            <Text
              style={[
                styles.bubbleText,
                m.from === 'me' ? styles.bubbleTextMe : styles.bubbleTextDriver,
              ]}
            >
              {m.msg}
            </Text>
          </View>
        ))}
      </ScrollView>

      {/* Input */}
      <View style={styles.inputBar}>
        <View style={styles.quickActions}>
          {['📍', '📷', '🎤'].map((e) => (
            <View key={e} style={styles.quickBtn}>
              <Text style={{ fontSize: 16 }}>{e}</Text>
            </View>
          ))}
        </View>
        <TextInput
          style={styles.textInput}
          placeholder="Type a message..."
          placeholderTextColor={colors.grey}
          accessibilityLabel="Message input"
        />
        <TouchableOpacity style={styles.sendBtn} accessibilityLabel="Send message">
          <Text style={{ fontSize: 16 }}>➤</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.bg },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: spacing[4],
    paddingVertical: spacing[2],
    gap: 10,
    borderBottomWidth: 1,
    borderBottomColor: colors.greyLight,
  },
  backBtn: {
    width: 32,
    height: 32,
    borderRadius: 10,
    backgroundColor: colors.light,
    alignItems: 'center',
    justifyContent: 'center',
  },
  backArrow: { fontSize: 14, color: colors.navy },
  headerName: { fontFamily: fontFamily.bold, fontSize: 13, fontWeight: '700', color: colors.navy },
  headerStatus: { fontFamily: fontFamily.regular, fontSize: 10, color: colors.green },
  messages: { flex: 1 },
  messagesList: { padding: spacing[4], gap: spacing[2] },
  timestamp: { alignItems: 'center', marginBottom: spacing[2] },
  timestampText: {
    fontFamily: fontFamily.regular,
    fontSize: 9,
    color: colors.grey,
    backgroundColor: colors.light,
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 6,
  },
  bubble: { maxWidth: '75%', padding: 10, paddingHorizontal: 14, borderRadius: 16 },
  bubbleMe: {
    alignSelf: 'flex-end',
    backgroundColor: colors.orange,
    borderBottomRightRadius: 4,
  },
  bubbleDriver: {
    alignSelf: 'flex-start',
    backgroundColor: colors.white,
    borderBottomLeftRadius: 4,
    borderWidth: 1,
    borderColor: colors.greyLight,
  },
  bubbleText: { fontFamily: fontFamily.regular, fontSize: 12, lineHeight: 17 },
  bubbleTextMe: { color: colors.white },
  bubbleTextDriver: { color: colors.navy },
  inputBar: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: spacing[4],
    paddingVertical: spacing[2],
    gap: spacing[2],
    borderTopWidth: 1,
    borderTopColor: colors.greyLight,
  },
  quickActions: { flexDirection: 'row', gap: 6 },
  quickBtn: {
    width: 36,
    height: 36,
    borderRadius: 12,
    backgroundColor: colors.light,
    alignItems: 'center',
    justifyContent: 'center',
  },
  textInput: {
    flex: 1,
    paddingVertical: 10,
    paddingHorizontal: 14,
    borderRadius: 12,
    backgroundColor: colors.white,
    borderWidth: 1.5,
    borderColor: colors.greyLight,
    fontFamily: fontFamily.regular,
    fontSize: 12,
    color: colors.navy,
  },
  sendBtn: {
    width: 36,
    height: 36,
    borderRadius: 12,
    backgroundColor: colors.teal,
    alignItems: 'center',
    justifyContent: 'center',
  },
});
