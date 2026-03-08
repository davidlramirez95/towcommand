/**
 * Theme System Tests — 2nd Order Logic
 *
 * 1st order: "Colors exist" — trivial.
 * 2nd order: Every component spreads theme tokens into StyleSheet.create().
 *   If ANY token is undefined, NaN, or wrong type, it silently creates
 *   invisible/broken UI rather than crashing. These tests catch that class
 *   of bugs by validating the contract between theme and consumers.
 */
import { colors } from '@/lib/theme/colors';
import { fontFamily, fontSize, textStyles } from '@/lib/theme/typography';
import { spacing, borderRadius } from '@/lib/theme/spacing';

// --- Colors ---

describe('colors', () => {
  const HEX_REGEX = /^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$/;
  const RGBA_REGEX = /^rgba?\(/;

  it('every color value is a valid hex or rgba string', () => {
    const allColors = Object.entries(colors);
    expect(allColors.length).toBeGreaterThan(0);

    for (const [key, value] of allColors) {
      expect(typeof value).toBe('string');
      const isValid = HEX_REGEX.test(value as string) || RGBA_REGEX.test(value as string);
      expect(isValid).toBe(true);
    }
  });

  it('has all required brand colors (prevents silent fallback to system defaults)', () => {
    // These 4 are used in every screen's StyleSheet
    expect(colors.navy).toBeDefined();
    expect(colors.teal).toBeDefined();
    expect(colors.orange).toBeDefined();
    expect(colors.gold).toBeDefined();
  });

  it('has semantic aliases (components reference these, not raw colors)', () => {
    expect(colors.success).toBeDefined();
    expect(colors.warning).toBeDefined();
    expect(colors.error).toBeDefined();
    expect(colors.info).toBeDefined();
  });

  it('bg color is distinct from white (ensures cream background renders visually)', () => {
    expect(colors.bg).not.toBe(colors.white);
  });

  it('overlay has alpha channel (needed for modal backdrops)', () => {
    expect(colors.overlay).toMatch(/rgba/);
  });
});

// --- Typography ---

describe('typography', () => {
  it('fontFamily values are non-empty strings (RN falls back to system font on empty)', () => {
    expect(typeof fontFamily.regular).toBe('string');
    expect(typeof fontFamily.semiBold).toBe('string');
    expect(typeof fontFamily.bold).toBe('string');
    expect(fontFamily.regular.length).toBeGreaterThan(0);
    expect(fontFamily.semiBold.length).toBeGreaterThan(0);
    expect(fontFamily.bold.length).toBeGreaterThan(0);
  });

  it('all fontSize values are positive numbers (0 or NaN renders invisible text)', () => {
    for (const [key, value] of Object.entries(fontSize)) {
      expect(typeof value).toBe('number');
      expect(value).toBeGreaterThan(0);
    }
  });

  it('textStyles produce valid StyleSheet objects (spread into components)', () => {
    const styles = [
      textStyles.h1, textStyles.h2, textStyles.h3,
      textStyles.body, textStyles.bodySmall,
      textStyles.caption, textStyles.button, textStyles.label,
    ];

    for (const style of styles) {
      expect(style).toBeDefined();
      expect(typeof style.fontSize).toBe('number');
      expect(typeof style.fontFamily).toBe('string');
      expect(style.fontSize).toBeGreaterThan(0);
    }
  });

  it('heading hierarchy: h1 > h2 > h3 (visual hierarchy correctness)', () => {
    expect(textStyles.h1.fontSize).toBeGreaterThan(textStyles.h2.fontSize);
    expect(textStyles.h2.fontSize).toBeGreaterThan(textStyles.h3.fontSize);
  });

  it('body > caption (readable text should be larger than auxiliary text)', () => {
    expect(textStyles.body.fontSize).toBeGreaterThan(textStyles.caption.fontSize);
  });
});

// --- Spacing ---

describe('spacing', () => {
  it('spacing scale has no undefined gaps (spacing[n] used as array index in every screen)', () => {
    // Components use spacing[0] through spacing[10], some use spacing[12] and spacing[16]
    const requiredIndices = [0, 1, 2, 3, 4, 5, 6, 8, 10, 12, 16];
    for (const idx of requiredIndices) {
      expect(spacing[idx]).toBeDefined();
      expect(typeof spacing[idx]).toBe('number');
    }
  });

  it('spacing values are non-negative (negative padding/margin breaks layout)', () => {
    for (const value of Object.values(spacing)) {
      expect(value).toBeGreaterThanOrEqual(0);
    }
  });

  it('spacing is monotonically increasing (larger index = more space)', () => {
    const keys = Object.keys(spacing).map(Number).sort((a, b) => a - b);
    for (let i = 1; i < keys.length; i++) {
      expect(spacing[keys[i]]).toBeGreaterThanOrEqual(spacing[keys[i - 1]]);
    }
  });

  it('borderRadius values are positive numbers (0 is valid but unlikely for cards)', () => {
    expect(typeof borderRadius.sm).toBe('number');
    expect(typeof borderRadius.md).toBe('number');
    expect(typeof borderRadius.lg).toBe('number');
    expect(typeof borderRadius.xl).toBe('number');
    expect(typeof borderRadius.pill).toBe('number');
    expect(borderRadius.sm).toBeGreaterThan(0);
  });

  it('borderRadius hierarchy: sm < md < lg < xl < pill (design consistency)', () => {
    expect(borderRadius.sm).toBeLessThan(borderRadius.md);
    expect(borderRadius.md).toBeLessThan(borderRadius.lg);
    expect(borderRadius.lg).toBeLessThan(borderRadius.xl);
    expect(borderRadius.xl).toBeLessThan(borderRadius.pill);
  });
});
