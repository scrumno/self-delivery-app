import * as SecureStore from 'expo-secure-store';

async function readNativeRaw(key: string): Promise<string | null> {
    return await SecureStore.getItemAsync(key);
}

async function writeNativeRaw(key: string, value: string): Promise<void> {
    await SecureStore.setItemAsync(key, value, { keychainAccessible: SecureStore.WHEN_UNLOCKED_THIS_DEVICE_ONLY });
}

async function removeNativeKey(key: string): Promise<void> {
    try {
        await SecureStore.deleteItemAsync(key);
    } catch {
        return;
    }
}

export { readNativeRaw, writeNativeRaw, removeNativeKey };