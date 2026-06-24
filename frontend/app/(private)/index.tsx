import { HomePage } from '@pages/home'
import { Redirect } from 'expo-router'
import { ActivityIndicator, StyleSheet, View } from 'react-native'
import { useAuth } from '../../provider/auth-provider'

export default function PrivateHomeRoute() {
  const { isLoading, profileOnboardingPending } = useAuth()

  if (isLoading) {
    return (
      <View style={styles.loader}>
        <ActivityIndicator size="large" />
      </View>
    )
  }

  if (profileOnboardingPending) {
    return <Redirect href="/(private)/profile" />
  }

  return <HomePage />
}

const styles = StyleSheet.create({
  loader: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
})
