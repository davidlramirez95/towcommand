import { Tabs } from 'expo-router';
import { View, StyleSheet } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';

export default function TabsLayout() {
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarActiveTintColor: colors.orange,
        tabBarInactiveTintColor: colors.grey,
        tabBarLabelStyle: {
          fontFamily: fontFamily.semiBold,
          fontSize: 10,
        },
        tabBarStyle: {
          backgroundColor: colors.white,
          borderTopColor: colors.greyLight,
          borderTopWidth: 1,
          paddingTop: 4,
        },
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: 'Home',
          tabBarIcon: ({ color, size }) => <TabIcon name="home" color={color} size={size} />,
        }}
      />
      <Tabs.Screen
        name="history"
        options={{
          title: 'Activity',
          tabBarIcon: ({ color, size }) => <TabIcon name="history" color={color} size={size} />,
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: 'Account',
          tabBarIcon: ({ color, size }) => <TabIcon name="profile" color={color} size={size} />,
        }}
      />
    </Tabs>
  );
}

/**
 * Minimal tab icons using View shapes.
 * Replace with proper SVG icons in production.
 */
function TabIcon({ name, color, size }: { name: string; color: string; size: number }) {
  return (
    <View
      style={[
        styles.iconPlaceholder,
        { width: size, height: size, borderColor: color },
      ]}
    />
  );
}

const styles = StyleSheet.create({
  iconPlaceholder: {
    borderWidth: 2,
    borderRadius: 6,
  },
});
