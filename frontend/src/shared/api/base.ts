import { clearStoredUser, getStoredUserJson, persistUser } from '@shared/lib/secure-store/handlers';
import Constants from 'expo-constants';
import { Platform } from 'react-native';

const setUserTokens = async (accessToken: string, refreshToken: string, expiresIn: number) => {
    try {
        await persistUser('user', { authorizeToken: accessToken, refreshToken, expiresIn: expiresIn });
    } catch (error) {
        console.error('Ошибка сохранения user токенов:', error);
    }
}; 

const getToken = async () => {
    const user = await getStoredUserJson('user');
    
    if (!user?.authorizeToken) {
        return null;
    }

    return user.authorizeToken;
};

const getRefreshToken = async () => {
    const user = await getStoredUserJson('user');
    
    if (!user?.refreshToken) {
        return null;
    }

    return user.refreshToken;
};

const clearUserTokens = async () => {
    await clearStoredUser('user');
};

const resolveApiBaseUrl = (): string => {
    const fromEnv = process.env.EXPO_PUBLIC_API_URL?.trim();
    if (fromEnv) {
        const base = fromEnv.replace(/\/?$/, '');
        return base.endsWith('/api/v1') ? base : `${base}/api/v1`;
    }

    if (Platform.OS === 'web') {
        return 'http://localhost:8087/api/v1';
    }

    const hostUri = Constants.expoConfig?.hostUri;
    const metroHost =
        typeof hostUri === 'string' && hostUri.length > 0 ? hostUri.split(':')[0] : '';

    if (metroHost && metroHost !== 'localhost' && metroHost !== '127.0.0.1') {
        return `http://${metroHost}:8087/api/v1`;
    }

    if (Platform.OS === 'android') {
        return 'http://10.0.2.2:8087/api/v1';
    }

    return 'http://localhost:8087/api/v1';
}

export { setUserTokens, getRefreshToken, getToken, getStoredUserJson,resolveApiBaseUrl, clearUserTokens };