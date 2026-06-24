import { useMutation } from "@tanstack/react-query";
import { authorizeUser, LoginData, RegisterData, AuthorizeResponseData, registerUser } from "../api/authorize";
import { router, useLocalSearchParams } from "expo-router";
import { AuthorizationCodeType } from "../api/variables";
import { useAuth } from "../../../../provider/auth-provider";
import { useState } from "react";

export const useLogin = () => {
    const { phone, codeType } = useLocalSearchParams<{ phone?: string, codeType?: AuthorizationCodeType }>();
    const { saveStoredData, startProfileOnboarding } = useAuth();
    const [error, setError] = useState<string>('');
    const [code, setCode] = useState<string>('');

    const handleAuthorize = async (data: AuthorizeResponseData, variables: LoginData) => {
        if (!data.isSuccess && data.error) {
            setError(data.error);
            return;
        }
        
        if (data.isSuccess) {
            await saveStoredData('user', {
                authorizeToken: data.accessToken,
                refreshToken: data.refreshToken,
                expiresIn: data.expiresIn,
            });

            if (codeType === AuthorizationCodeType.REGISTER) {
                await startProfileOnboarding();
                router.replace({ pathname: '/(private)/profile' });
                return;
            }

            router.replace('/(private)');
        }
    }

    const { mutate: requestAuthorizeUser } = useMutation({
        mutationFn: (data: LoginData) => authorizeUser(data),
        onSuccess: async (data: AuthorizeResponseData, variables: LoginData) => handleAuthorize(data, variables),
        onError: (error: any) => {
            if (!error.response?.data?.isSuccess) {
                setError(error.response?.data?.error);
            }
        },
    })

    const { mutate: requestRegisterUser } = useMutation({
        mutationFn: (data: RegisterData) => registerUser(data),
        onSuccess: async (data:  AuthorizeResponseData, variables: RegisterData) => handleAuthorize(data, variables),
        onError: (error: any) => {
            if (!error.response?.data?.isSuccess) {
                setError(error.response?.data?.error);
            }
        },
    })

    const requestHandler = (code: string) => {
        switch (codeType) {
            case AuthorizationCodeType.AUTHORIZE:
                requestAuthorizeUser({ code, phone: phone ?? '' });
                break;
            case AuthorizationCodeType.REGISTER:
                requestRegisterUser({ code, phone: phone ?? '' });
                break;
        }
    }

    return {
        phone,
        requestHandler,
        code,
        setCode,
        error,
    }
}