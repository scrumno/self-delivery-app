import { typography } from '@shared/ui/config/typography'
import { StyleProp, StyleSheet, Text, TextProps, TextStyle } from 'react-native'

type Variant = 'h1' | 'h2' | 'h3'

type Props = {
  variant: Variant
  color?: string
  content: string

  style?: StyleProp<TextStyle>
  textProps?: TextProps
}

const Heading = ({ variant, color, content, style, textProps }: Props) => {
  switch (variant) {
    case 'h1':
      return <Text style={[styles.h1, { color }]}>{content}</Text>
    case 'h2':
      return <Text style={[styles.h2, { color }]}>{content}</Text>
    case 'h3':
      return <Text style={[styles.h3, { color }]}>{content}</Text>
  }
}

const styles = StyleSheet.create({
  h3: {
    fontSize: typography.fontSize.md,
    fontWeight: 700,
    lineHeight: typography.lineHeight.md,
    letterSpacing: typography.letterSpacing.xxs,
  },
  h2: {
    fontSize: typography.fontSize.lg,
    fontWeight: 600,
    lineHeight: typography.lineHeight.md,
    letterSpacing: typography.letterSpacing.xxs,
  },
  h1: {
    fontSize: typography.fontSize.xxl,
    fontWeight: 700,
    lineHeight: typography.lineHeight.md,
    letterSpacing: typography.letterSpacing.xxs,
  },
})

export { Heading }
