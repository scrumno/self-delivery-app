import MaterialIcons from '@expo/vector-icons/MaterialIcons'
import { BottomTabBar, type BottomTabBarProps } from '@react-navigation/bottom-tabs'
import { Redirect, Tabs } from 'expo-router'
import { ActivityIndicator, StyleSheet, Text, View } from 'react-native'
import { useAuth } from '../../provider/auth-provider'

export default function PrivateLayout() {
  const { user, isLoading, profileOnboardingPending } = useAuth();

  const isAuthorized =
    Boolean(user.authorizeToken?.length) && Boolean(user.refreshToken?.length);

  if (!isAuthorized) {
    return <Redirect href="/(public)" />;
  }

  if (isLoading) {
    return (
      <View style={styles.loader}>
        <ActivityIndicator size="large" />
        <Text style={styles.loaderText}>Загрузка…</Text>
      </View>
    );
  }

  return (
    <Tabs
      initialRouteName={profileOnboardingPending ? 'profile' : undefined}
      tabBar={(props: BottomTabBarProps) =>
        profileOnboardingPending ? null : <BottomTabBar {...props} />
      }
      screenOptions={{ headerShown: false }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: 'Главная',
          href: profileOnboardingPending ? null : undefined,
          tabBarIcon: ({ color, size }) => (
            <MaterialIcons name="home" size={size ?? 24} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: 'Профиль',
          tabBarIcon: ({ color, size }) => (
            <MaterialIcons name="person" size={size ?? 24} color={color} />
          ),
        }}
      />
    </Tabs>
  )
}

const styles = StyleSheet.create({
  loader: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 12,
  },
  loaderText: {
    fontSize: 16,
    color: '#666',
  },
})
