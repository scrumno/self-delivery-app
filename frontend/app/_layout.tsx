import { useOnestFonts } from '@shared/ui/config/fonts';
import { DarkTheme, DefaultTheme, ThemeProvider } from '@react-navigation/native';
import { Stack } from 'expo-router';
import * as SplashScreen from 'expo-splash-screen';
import { StatusBar } from 'expo-status-bar';
import { useEffect } from 'react';
import { useColorScheme } from 'react-native';
import { AuthProvider, useAuth } from '../provider/auth-provider';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { QueryProvider } from '../provider/query-provider';
import 'react-native-reanimated';

function RootStack() {
  const { user } = useAuth();
  const isAuthorized =
    Boolean(user.authorizeToken?.length) && Boolean(user.refreshToken?.length);

  if (isAuthorized) {
    return (
      <Stack key="private-stack" screenOptions={{ headerShown: false }}>
        <Stack.Screen name="(private)" />
      </Stack>
    );
  }

  return (
    <Stack key="public-stack" screenOptions={{ headerShown: false }}>
      <Stack.Screen name="(public)" />
    </Stack>
  );
}

export default function RootLayout() {
  const colorScheme = useColorScheme();
  const [fontsLoaded] = useOnestFonts();

  useEffect(() => {
    if (fontsLoaded) {
      SplashScreen.hideAsync();
    }
  }, [fontsLoaded]);

  if (!fontsLoaded) {
    return null;
  }

  let currentTheme = colorScheme === 'dark' ? DarkTheme : DefaultTheme;
  currentTheme = DefaultTheme;

  return (
    <SafeAreaProvider>
      <ThemeProvider value={currentTheme}>
        <StatusBar style="auto" />
        <AuthProvider>
          <QueryProvider>
            <RootStack />
          </QueryProvider>
        </AuthProvider>
      </ThemeProvider>
    </SafeAreaProvider>
  );
}
