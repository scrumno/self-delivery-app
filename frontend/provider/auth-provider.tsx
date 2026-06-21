import React, {
  createContext,
  useState,
  useContext,
  useEffect,
  useCallback,
} from 'react';
import { Platform } from 'react-native';
import {
  getStoredUserJson,
  persistUser,
  clearStoredUser,
  getProfileOnboardingPending,
  persistProfileOnboardingPending,
  clearProfileOnboardingPending,
  type UserSecureStore,
} from '@shared/lib/secure-store/handlers';

const initialUser: UserSecureStore = {
  authorizeToken: null,
  refreshToken: null,
  expiresIn: null,
};

const USER_KEY = 'user';

/**
 * Веб и HttpOnly / Secure cookie
 * --------------------------------
 * JavaScript не может записать HttpOnly-куку: её выставляет только сервер
 * в ответе (Set-Cookie: …; HttpOnly; Secure; SameSite=Strict|Lax).
 *
 * Тогда на web:
 * 1) не кладите access/refresh в sessionStorage/localStorage;
 * 2) axios/fetch с withCredentials: true (у вас в shared/api/api.ts для web);
 * 3) бэкенд принимает сессию по куке и при логауте сбрасывает её тем же механизмом;
 * 4) фронт узнаёт «залогинен ли пользователь» по успешному GET /me или 204 после логина,
 *    а не по чтению токена из document.
 *
 * Пока API отдаёт токены в JSON, для web остаётся sessionStorage (видно из JS при XSS).
 * После перевода авторизации на куки — замените readWebUserJson/writeWebUserJson на
 * вызов /me и уберите запись токенов из saveStoredData.
 */

type AuthContextType = {
  user: UserSecureStore;
  isLoading: boolean;
  profileOnboardingPending: boolean;
  saveStoredData: (key: string, value: UserSecureStore) => Promise<void>;
  clearUserTokenFromStorage: () => Promise<void>;
  startProfileOnboarding: () => Promise<void>;
  completeProfileOnboarding: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<UserSecureStore>(initialUser);
  const [isLoading, setIsLoading] = useState(true);
  const [profileOnboardingPending, setProfileOnboardingPending] = useState(false);

  const syncUserFromStorage = useCallback(async () => {
    try {
      const [raw, onboarding] = await Promise.all([
        getStoredUserJson(USER_KEY),
        getProfileOnboardingPending(),
      ]);

      setProfileOnboardingPending(onboarding);

      if (!raw) {
        setUser(initialUser);
        return;
      }

      setUser({
        authorizeToken: raw.authorizeToken ?? null,
        refreshToken: raw.refreshToken ?? null,
        expiresIn: raw.expiresIn ?? null,
      });
    } catch (error) {
      console.error('Произошла ошибка при загрузке данных авторизации', error);
      setUser(initialUser);
      setProfileOnboardingPending(false);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const saveStoredData = async (key: string, value: UserSecureStore) => {
    try {
      await persistUser(key, value);
      setUser(value);
    } catch (error) {
      console.error('Произошла ошибка при сохранении данных авторизации', error);
    }
  };

  const clearUserTokenFromStorage = useCallback(async () => {
    try {
      await clearStoredUser(USER_KEY);
      await clearProfileOnboardingPending();
      setProfileOnboardingPending(false);
      setUser(initialUser);
    } catch (error) {
      console.error('Произошла ошибка при выходе из системы', error);
    }
  }, []);

  const startProfileOnboarding = useCallback(async () => {
    await persistProfileOnboardingPending(true);
    setProfileOnboardingPending(true);
  }, []);

  const completeProfileOnboarding = useCallback(async () => {
    await clearProfileOnboardingPending();
    setProfileOnboardingPending(false);
  }, []);

  useEffect(() => {
    void syncUserFromStorage();
  }, [syncUserFromStorage]);

  useEffect(() => {
    if (Platform.OS !== 'web' || typeof window === 'undefined') {
      return;
    }

    const onFocus = () => {
      void syncUserFromStorage();
    };

    window.addEventListener('focus', onFocus);
    document.addEventListener('visibilitychange', onFocus);

    return () => {
      window.removeEventListener('focus', onFocus);
      document.removeEventListener('visibilitychange', onFocus);
    };
  }, [syncUserFromStorage]);

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        profileOnboardingPending,
        saveStoredData,
        clearUserTokenFromStorage,
        startProfileOnboarding,
        completeProfileOnboarding,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error('useAuth должен вызываться внутри AuthProvider');
  }

  return context;
};
