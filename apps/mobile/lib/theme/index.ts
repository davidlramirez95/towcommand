export { colors } from './colors';
export { fontFamily, fontSize, textStyles } from './typography';
export { spacing, borderRadius } from './spacing';

export const theme = {
  colors: require('./colors').colors,
  fonts: require('./typography'),
  spacing: require('./spacing').spacing,
  borderRadius: require('./spacing').borderRadius,
} as const;
