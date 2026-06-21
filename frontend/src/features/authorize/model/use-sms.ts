import { useMutation } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react"
import { Code, getSmsCode, SmssResponseData } from "../api/authorize";
import { router, useLocalSearchParams } from "expo-router";
import { AuthorizationCodeType } from "../api/variables";
import { handleError } from "../api/exception";

export const useSms = () => {
    const { codeType, phone: phoneParam } = useLocalSearchParams<{ codeType?: AuthorizationCodeType; phone?: string }>();
    const [phone, setPhone] = useState<string>(phoneParam ?? '')
    const [error, setError] = useState<string>('');

    const isRegister = useMemo(() => codeType === AuthorizationCodeType.REGISTER, [codeType]);

    const { mutate: requestSmsCode, isPending } = useMutation({
        mutationFn: (data: Code) => getSmsCode(data),
        onSuccess: (data: SmssResponseData, variables: Code) => {
            
            if (data.isSuccess) {
                router.push({
                    pathname: '/(public)/code',
                    params: { 
                        phone: variables.phone,
                        codeType: codeType,
                    },
                });
            }

            if (!data.isSuccess && variables.codeType) {
                const currentRequesttype = variables.codeType === AuthorizationCodeType.AUTHORIZE ? 
                AuthorizationCodeType.REGISTER : AuthorizationCodeType.AUTHORIZE;

                handleError(
                    variables.phone, 
                    () => router.push({
                        pathname: '/(public)/authorize',
                        params: { phone: phone, codeType: currentRequesttype },
                    }),
                    () => setError(''), 
                    data.code
                );
            }

            setError(data.error);
        },
        onError: (error) => {
            setError(error.message);
        },
    })

    const handleRequestSmsCode = () => {
        if (phone && codeType) requestSmsCode({ phone, codeType });
    }

    useEffect(() => {
        if (phoneParam) {
            setPhone(phoneParam);
        }
    }, [phoneParam]);

    return {
        phone,
        isRegister,
        setPhone,
        error,
        handleRequestSmsCode,
        isSmsPending: isPending,
    }
}