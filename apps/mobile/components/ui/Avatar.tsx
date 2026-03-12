import { View, Text, StyleSheet } from 'react-native';
import { colors } from '@/lib/theme/colors';
import { fontFamily } from '@/lib/theme/typography';

interface AvatarProps {
  name: string;
  size?: number;
}

function getInitials(name: string): string {
  return name
    .split(' ')
    .filter(Boolean)
    .map((w) => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);
}

export function Avatar({ name, size = 50 }: AvatarProps) {
  const initials = getInitials(name);
  return (
    <View
      style={[
        styles.container,
        {
          width: size,
          height: size,
          borderRadius: size * 0.32,
        },
      ]}
      accessibilityLabel={`Avatar for ${name}`}
    >
      <Text
        style={[
          styles.initials,
          { fontSize: size * 0.38 },
        ]}
      >
        {initials}
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.orange,
    alignItems: 'center',
    justifyContent: 'center',
  },
  initials: {
    fontFamily: fontFamily.bold,
    fontWeight: '800',
    color: colors.white,
  },
});
