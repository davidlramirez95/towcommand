import { useState, useCallback } from 'react';
import { View, Text, StyleSheet, Pressable, Alert, Vibration } from 'react-native';
import { router } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Button } from '@/components/ui/Button';
import { colors } from '@/lib/theme/colors';
import { textStyles, fontFamily } from '@/lib/theme/typography';
import { spacing, borderRadius } from '@/lib/theme/spacing';
import { useAuth } from '@/hooks/useAuth';
import { api } from '@/lib/api';

export default function SOSScreen() {
  const [triggered, setTriggered] = useState(false);
  const [loading, setLoading] = useState(false);
  const { user } = useAuth();

  const handleTriggerSOS = useCallback(async () => {
    setLoading(true);
    Vibration.vibrate([0, 100, 50, 100]); // Haptic confirmation

    try {
      await api.post('/sos/trigger', {
        userId: user?.id,
        lat: 0, // Will use actual GPS
        lng: 0,
        method: 'MANUAL_BUTTON',
      });
      setTriggered(true);
    } catch (err) {
      Alert.alert('Error', 'Failed to send SOS. Please call 911 directly.');
    } finally {
      setLoading(false);
    }
  }, [user]);

  if (triggered) {
    return (
      <SafeAreaView style={styles.containerTriggered}>
        <View style={styles.triggeredContent}>
          <View style={styles.pulseCircle} />
          <Text style={styles.triggeredTitle}>Help is on the way</Text>
          <Text style={styles.triggeredSubtitle}>
            TowCommand safety team has been alerted.{'\n'}
            Stay calm and stay in a safe location.
          </Text>
          <View style={styles.emergencyInfo}>
            <Text style={styles.emergencyLabel}>Emergency Hotline</Text>
            <Text style={styles.emergencyNumber}>911</Text>
          </View>
          <Button variant="secondary" onPress={() => router.back()} fullWidth>
            Close
          </Button>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      {/* Close Button */}
      <Pressable
        style={styles.closeButton}
        onPress={() => router.back()}
        accessibilityLabel="Close SOS"
        accessibilityRole="button"
      >
        <Text style={styles.closeText}>X</Text>
      </Pressable>

      <View style={styles.content}>
        <Text style={styles.title}>Emergency SOS</Text>
        <Text style={styles.subtitle}>
          Press the button below to alert the TowCommand safety team and nearby authorities.
        </Text>

        {/* SOS Button */}
        <Pressable
          style={styles.sosButton}
          onPress={handleTriggerSOS}
          disabled={loading}
          accessibilityLabel="Trigger SOS Alert"
          accessibilityRole="button"
        >
          <Text style={styles.sosButtonText}>SOS</Text>
          <Text style={styles.sosButtonSubtext}>Press to alert</Text>
        </Pressable>

        <Text style={styles.disclaimer}>
          This will notify the TowCommand safety team with your current location.
          For life-threatening emergencies, call 911 directly.
        </Text>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.navy,
  },
  containerTriggered: {
    flex: 1,
    backgroundColor: colors.coral,
  },
  closeButton: {
    position: 'absolute',
    top: spacing[12],
    right: spacing[5],
    width: 36,
    height: 36,
    borderRadius: borderRadius.full,
    backgroundColor: 'rgba(255,255,255,0.15)',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1,
  },
  closeText: {
    fontFamily: fontFamily.bold,
    fontSize: 16,
    color: colors.white,
  },
  content: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: spacing[8],
  },
  title: {
    ...textStyles.h1,
    color: colors.white,
    marginBottom: spacing[2],
  },
  subtitle: {
    ...textStyles.body,
    color: 'rgba(255,255,255,0.7)',
    textAlign: 'center',
    marginBottom: spacing[10],
  },
  sosButton: {
    width: 180,
    height: 180,
    borderRadius: 90,
    backgroundColor: colors.coral,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: colors.coral,
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.5,
    shadowRadius: 20,
    elevation: 8,
    marginBottom: spacing[10],
  },
  sosButtonText: {
    fontFamily: fontFamily.bold,
    fontSize: 40,
    color: colors.white,
    letterSpacing: 4,
  },
  sosButtonSubtext: {
    fontFamily: fontFamily.semiBold,
    fontSize: 12,
    color: 'rgba(255,255,255,0.7)',
    marginTop: 4,
  },
  disclaimer: {
    ...textStyles.caption,
    color: 'rgba(255,255,255,0.4)',
    textAlign: 'center',
    paddingHorizontal: spacing[4],
  },
  triggeredContent: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: spacing[8],
  },
  pulseCircle: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: 'rgba(255,255,255,0.2)',
    marginBottom: spacing[8],
  },
  triggeredTitle: {
    ...textStyles.h1,
    color: colors.white,
    marginBottom: spacing[3],
  },
  triggeredSubtitle: {
    ...textStyles.body,
    color: 'rgba(255,255,255,0.8)',
    textAlign: 'center',
    marginBottom: spacing[8],
  },
  emergencyInfo: {
    alignItems: 'center',
    marginBottom: spacing[8],
  },
  emergencyLabel: {
    ...textStyles.caption,
    color: 'rgba(255,255,255,0.6)',
    marginBottom: spacing[1],
  },
  emergencyNumber: {
    fontFamily: fontFamily.bold,
    fontSize: 32,
    color: colors.white,
  },
});
