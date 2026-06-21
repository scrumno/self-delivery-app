import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';

export const nativeSecureStoreOptions: SecureStore.SecureStoreOptions =
    Platform.OS === 'ios'
        ? { keychainAccessible: SecureStore.WHEN_UNLOCKED_THIS_DEVICE_ONLY }
        : {};
