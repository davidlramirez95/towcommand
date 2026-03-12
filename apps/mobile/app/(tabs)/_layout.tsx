import { Tabs } from 'expo-router';
import { Text } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';

const TAB_ICONS: Record<string, string> = {
  home: '🏠',
  history: '📋',
  suki: '⭐',
  profile: '👤',
};

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
          tabBarIcon: ({ focused, size }) => (
            <TabIcon name="home" focused={focused} size={size} />
          ),
        }}
      />
      <Tabs.Screen
        name="history"
        options={{
          title: 'Activity',
          tabBarIcon: ({ focused, size }) => (
            <TabIcon name="history" focused={focused} size={size} />
          ),
        }}
      />
      <Tabs.Screen
        name="suki"
        options={{
          title: 'Suki',
          tabBarIcon: ({ focused, size }) => (
            <TabIcon name="suki" focused={focused} size={size} />
          ),
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: 'Account',
          tabBarIcon: ({ focused, size }) => (
            <TabIcon name="profile" focused={focused} size={size} />
          ),
        }}
      />
    </Tabs>
  );
}

function TabIcon({ name, focused, size }: { name: string; focused: boolean; size: number }) {
  return (
    <Text style={{ fontSize: size, opacity: focused ? 1 : 0.4 }}>
      {TAB_ICONS[name] ?? '•'}
    </Text>
  );
}
