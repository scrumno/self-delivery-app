
import { Redirect, Stack } from 'expo-router'
import { useAuth } from '../../provider/auth-provider'
import { Text, View } from 'react-native';

export default function PublicLayout() {
  const { user, isLoading } = useAuth();
  const isAuthorized =
    Boolean(user.authorizeToken?.length) && Boolean(user.refreshToken?.length);

  if (isLoading) {
    return (
      <View>
        <Text>Loading..</Text>
      </View>
    );
  }

  if (isAuthorized) {
    return <Redirect href="/(private)" />;
  }

  return (
    <Stack screenOptions={{ headerShown: false }}>
      <Stack.Screen name="index" />
      <Stack.Screen name="authorize" />
      <Stack.Screen name="code" />
    </Stack>
  )
}
