import { TextStyle } from 'react-native';

/**
 * Font family names registered by @expo-google-fonts/poppins.
 * These names work on both native (via expo-font) and web (injected @font-face).
 */
export const fontFamily = {
  regular: 'Poppins_400Regular',
  semiBold: 'Poppins_600SemiBold',
  bold: 'Poppins_700Bold',
} as const;

export const fontSize = {
  xs: 10,
  sm: 12,
  md: 14,
  lg: 16,
  xl: 18,
  '2xl': 20,
  '3xl': 24,
  '4xl': 28,
  '5xl': 32,
} as const;

export const lineHeight = {
  xs: 14,
  sm: 16,
  md: 20,
  lg: 22,
  xl: 24,
  '2xl': 28,
  '3xl': 32,
  '4xl': 36,
  '5xl': 40,
} as const;

export const textStyles = {
  h1: {
    fontFamily: fontFamily.bold,
    fontSize: fontSize['4xl'],
    lineHeight: lineHeight['4xl'],
  } satisfies TextStyle,
  h2: {
    fontFamily: fontFamily.bold,
    fontSize: fontSize['2xl'],
    lineHeight: lineHeight['2xl'],
  } satisfies TextStyle,
  h3: {
    fontFamily: fontFamily.semiBold,
    fontSize: fontSize.xl,
    lineHeight: lineHeight.xl,
  } satisfies TextStyle,
  body: {
    fontFamily: fontFamily.regular,
    fontSize: fontSize.md,
    lineHeight: lineHeight.md,
  } satisfies TextStyle,
  bodySmall: {
    fontFamily: fontFamily.regular,
    fontSize: fontSize.sm,
    lineHeight: lineHeight.sm,
  } satisfies TextStyle,
  caption: {
    fontFamily: fontFamily.regular,
    fontSize: fontSize.xs,
    lineHeight: lineHeight.xs,
  } satisfies TextStyle,
  button: {
    fontFamily: fontFamily.bold,
    fontSize: fontSize.md,
    lineHeight: lineHeight.md,
  } satisfies TextStyle,
  buttonSmall: {
    fontFamily: fontFamily.semiBold,
    fontSize: fontSize.sm,
    lineHeight: lineHeight.sm,
  } satisfies TextStyle,
  label: {
    fontFamily: fontFamily.semiBold,
    fontSize: fontSize.sm,
    lineHeight: lineHeight.sm,
  } satisfies TextStyle,
} as const;
