import axios, { AxiosRequestConfig } from 'axios';
import { Platform } from 'react-native';
import { handleError } from './exceptions';
import { getRefreshToken, getToken, clearUserTokens, resolveApiBaseUrl, setUserTokens } from './base';


const BASE_URL = resolveApiBaseUrl();
const TIMEOUT = 30000;

const createInstance = (config: AxiosRequestConfig = {}) => {
  const instance = axios.create({
    baseURL: BASE_URL,
    timeout: TIMEOUT,
    // withCredentials: Platform.OS === 'web',
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      ...config.headers,
    },
    ...config,
  });

  instance.interceptors.request.use(async (request) => {
    const token = await getToken();
    if (token) {
        request.headers.Authorization = `Bearer ${token}`;
    }
    
    return request;
    }, (error) => {
        return Promise.reject(error);
    });

  instance.interceptors.response.use((response) => {
      return response;
    },
    async (error) => {
      const originalRequest = error.config;
      
      if (error.response?.status === 401 && !originalRequest._retry) {        
        try {
          originalRequest._retry = true;
          const refreshToken = await getRefreshToken();
          if (refreshToken) {
            const response = await axios.post(
              `${BASE_URL}/auth/refresh-tokens`,
              { refreshToken },
              { withCredentials: Platform.OS === 'web' }
            );
            
            const { accessToken, refreshToken: nextRefreshToken, expiresIn } = response.data;
            if (!accessToken || !nextRefreshToken) {
              await clearUserTokens();
              throw error;
            }

            await setUserTokens(accessToken, nextRefreshToken, expiresIn);
            
            originalRequest.headers.Authorization = `Bearer ${accessToken}`;
            return instance(originalRequest);
          }
        } catch (refreshError) {
          await clearUserTokens();
        }
      }
      
      if (error.response) {
        const { status, data } = error.response;
        
        handleError(status, data);
        
        throw error;
      } else if (error.request) {
        console.error('Network Error:', error.request);
        throw new Error('Нет подключения к серверу');
      } else {
        console.error('Request Error:', error.message);
        throw error;
      }
    }
  );
  
  return instance;
};

const api = createInstance();

export default api;