export const colors = {
  // Primary
  navy: '#0B1D33',
  teal: '#00897B',
  gold: '#F5A623',
  orange: '#FF6B35',
  white: '#FFFFFF',

  // Secondary
  coral: '#FF4757',
  green: '#00C48C',
  blue: '#2196F3',

  // Neutrals
  cream: '#FFF9F0',
  light: '#F5F1EB',
  bg: '#FAF7F2',
  grey: '#8E9BAE',
  greyLight: '#E8E4DE',
  dark: '#1A1A2E',

  // Semantic
  success: '#00C48C',
  warning: '#F5A623',
  error: '#FF4757',
  info: '#2196F3',

  // Transparent
  overlay: 'rgba(11, 29, 51, 0.6)',
  shadow: 'rgba(11, 29, 51, 0.08)',
} as const;

export type ColorName = keyof typeof colors;
