import { Stack } from 'expo-router';
import { colors } from '@/lib/theme/colors';

export default function ProviderLayout() {
  return (
    <Stack
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: colors.bg },
      }}
    >
      <Stack.Screen name="dashboard" />
      <Stack.Screen name="earnings" />
    </Stack>
  );
}
