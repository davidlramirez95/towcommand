import { Stack } from 'expo-router';
import { colors } from '@/lib/theme/colors';

export default function BookingLayout() {
  return (
    <Stack
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: colors.bg },
        animation: 'slide_from_right',
      }}
    >
      <Stack.Screen name="service" />
      <Stack.Screen name="vehicle" />
      <Stack.Screen name="dropoff" />
      <Stack.Screen name="price" />
      <Stack.Screen name="matching" />
      <Stack.Screen name="matched" />
      <Stack.Screen name="[id]" />
      <Stack.Screen name="chat" />
      <Stack.Screen name="condition" />
      <Stack.Screen name="complete" />
      <Stack.Screen name="rate" />
      <Stack.Screen name="diagnose" />
      <Stack.Screen name="payment" />
    </Stack>
  );
}
