import { Platform } from 'react-native';
import { readWebUserJson, writeWebUserJson, removeWebUserJson } from './web/auth-secure-handlers';
import { readNativeRaw, writeNativeRaw, removeNativeKey } from './native/auth-secure-handlers';

export type UserSecureStore = {
    authorizeToken: string | null;
    refreshToken: string | null;
    expiresIn: number | null;
};

export const getStoredUserJson = async (key: string): Promise<UserSecureStore | null> => {
    
    const raw = Platform.OS === 'web' ? readWebUserJson(key) : await readNativeRaw(key);
    return raw ? JSON.parse(raw) : null;
};

export const persistUser = async (key: string, value: UserSecureStore): Promise<void> => {
    const json = JSON.stringify(value);

    if (Platform.OS === 'web') {
        writeWebUserJson(key, json);
        return;
    }
    await writeNativeRaw(key, json);
};

export const clearStoredUser = async (key: string): Promise<void> => {
    if (Platform.OS === 'web') {
        removeWebUserJson(key);
        return;
    }
    await removeNativeKey(key);
};

/** Отдельный ключ (не JSON): нужно дополнить профиль перед главной. */
export const PROFILE_ONBOARDING_KEY = 'profile_onboarding_pending';

export async function getProfileOnboardingPending(): Promise<boolean> {
    const raw =
        Platform.OS === 'web'
            ? readWebUserJson(PROFILE_ONBOARDING_KEY)
            : await readNativeRaw(PROFILE_ONBOARDING_KEY);
    return raw === '1';
}

export async function persistProfileOnboardingPending(pending: boolean): Promise<void> {
    if (!pending) {
        await clearProfileOnboardingPending();
        return;
    }
    if (Platform.OS === 'web') {
        writeWebUserJson(PROFILE_ONBOARDING_KEY, '1');
        return;
    }
    await writeNativeRaw(PROFILE_ONBOARDING_KEY, '1');
}

export async function clearProfileOnboardingPending(): Promise<void> {
    if (Platform.OS === 'web') {
        removeWebUserJson(PROFILE_ONBOARDING_KEY);
        return;
    }
    await removeNativeKey(PROFILE_ONBOARDING_KEY);
}
