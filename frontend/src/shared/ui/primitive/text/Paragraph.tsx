import { typography } from '@shared/ui/config/typography'
import { StyleProp, StyleSheet, Text, TextProps, TextStyle } from 'react-native'

type Variant = 'small' | 'medium' | 'large' | 'verySmall'

type Props = {
  variant: Variant
  color?: string
  content: string

  style?: StyleProp<TextStyle>
  textProps?: TextProps
}

const Paragraph = ({ variant, color, content, style, textProps }: Props) => {
  switch (variant) {
    case 'verySmall':
      return (
        <Text
          style={[styles.verySmall, { color }, style]}
          {...textProps}
        >
          {content}
        </Text>
      )
    case 'small':
      return (
        <Text
          style={[styles.small, { color }, style]}
          {...textProps}
        >
          {content}
        </Text>
      )
    case 'medium':
      return (
        <Text
          style={[styles.medium, { color }, style]}
          {...textProps}
        >
          {content}
        </Text>
      )
    case 'large':
      return (
        <Text
          style={[styles.large, { color }, style]}
          {...textProps}
        >
          {content}
        </Text>
      )
    default:
      throw new Error(`Unknown variant: ${variant}`)
  }
}

const styles = StyleSheet.create({
  verySmall: {
    fontSize: typography.fontSize.xxs,
    fontWeight: 500,
    letterSpacing: typography.letterSpacing.xs,
  },
  small: {
    fontSize: typography.fontSize.xs,
    fontWeight: 400,
    letterSpacing: typography.letterSpacing.xxs,
  },
  base: {
    fontSize: typography.fontSize.xs,
    fontWeight: 400,
    letterSpacing: typography.letterSpacing.xxs,
  },
  medium: {
    fontSize: typography.fontSize.md,
    fontWeight: 400,
    letterSpacing: typography.letterSpacing.xxs,
  },
  large: {
    fontSize: typography.fontSize.md,
    letterSpacing: typography.letterSpacing.xxs,
  },
})

export { Paragraph }
