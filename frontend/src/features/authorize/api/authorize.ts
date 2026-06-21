import api from "@shared/api/api";
import { AuthorizationCodeType } from "./variables";

export interface Code {
    phone: string;
    codeType: AuthorizationCodeType;
};

export interface LoginData {
    code: string;
    phone: string;
}

export interface RegisterData {
    code: string;
    phone: string;
}

export interface SmssResponseData {
    isSuccess: boolean;
    code?: string;
    error: string;
}

export interface AuthorizeResponseData {
    isSuccess: boolean;
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    error?: string;
    code?: string;
}

export const getSmsCode = async (data: Code): Promise<SmssResponseData> => {
    try {
        const response = await api.post('/auth/sms-code', data);
        return response.data;
    } catch (error: any) {
        if (!error.response) {
            throw new Error('Произошла неизвестная ошибка');
        }

        return error.response.data;        
    }
}

export const authorizeUser = async (data: LoginData): Promise<AuthorizeResponseData> => {
    try {
      const response = await api.post('/auth/authorize', data);
      return response.data;
    } catch (error: any) {
      if (!error.response) {
        throw new Error('Произошла неизвестная ошибка');
      }

      return error.response.data;
    }
}

export const registerUser = async (data: RegisterData): Promise<AuthorizeResponseData> => {
    try {
        const response = await api.post('/auth/registration', data);
        return response.data;
    } catch (error: any) {
        if (!error.response) {
            throw new Error('Произошла неизвестная ошибка');
        }

        return error.response.data;
    }
}

