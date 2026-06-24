import { userStore } from "@entities/user/user.store";
import { useMutation } from "@tanstack/react-query";
import { router } from "expo-router";
import { useEffect, useState } from "react";
import { updateProfile } from "../api/profile";
import { UserFullProfile } from "./type";
import { useAuth } from "provider/auth-provider";

export const initialProfile: UserFullProfile = {
    firstName: "",
    secondName: "",
    birthDate: null,
    address: null,
    email: null,
    fullName: ""
};

export const usePersonalFields = () => {
    const { clearUserTokenFromStorage, profileOnboardingPending, completeProfileOnboarding } = useAuth();
    const [userProfile, setUserProfile] = useState<UserFullProfile>(initialProfile);
    const [error, setError] = useState('');

    const handleTextField = (field: keyof Pick<UserFullProfile, 'firstName' | 'secondName' | 'address' | 'email'>, value: string) => {
        setUserProfile(profile => ({
            ...profile,
            [field]: value
        }));
    }

    const { mutate: updateProfileRequest, isPending } = useMutation({
        mutationFn: (data: Partial<UserFullProfile>) => updateProfile(data),
        onSuccess: async (data: { isSuccess: boolean; error?: string }, variables) => {
            if (data.isSuccess) {
                if (variables) {
                    console.log('variables', variables);
                    userStore.hydrate({
                        firstName: variables.firstName ?? userStore.firstName,
                        secondName: variables.secondName ?? userStore.secondName,
                    });
                }

                if (profileOnboardingPending) {
                    await completeProfileOnboarding();
                    router.replace({ pathname: '/(private)' });
                }

                setError('');
            }

            setError(data.error ?? '');
        },
        onError: (error) => {
            setError(error.message);
        },
    })

    const createFormData = (): Partial<UserFullProfile> => {
        return Object.entries(userProfile).reduce<Partial<UserFullProfile>>((result, [key, value]) => {
            if (value && value != '' && value != null) {
                switch (key) {
                    case 'firstName':
                    case 'secondName':
                    case 'address':
                    case 'email':
                        result[key] = String(value);
                        break;
                    case 'birthDate':
                        result[key] = new Date(value);
                        break;
                }
            }
            return result;
        }, {});
    }

    const handleUpdateProfile = () => {
        updateProfileRequest(createFormData());
    }

    const handleLogout = async () => {
        await clearUserTokenFromStorage();

        userStore.reset();
        router.replace({ pathname: '/(public)' });
    }

    useEffect(() => {

        console.log('userStore.firstName', userStore.firstName);
        console.log('userStore.secondName', userStore.secondName);
        setUserProfile((prev) => ({
            ...prev,
            firstName: userStore.firstName || prev.firstName,
            secondName: userStore.secondName || prev.secondName,
        }));
    }, []);
       
    return {
        profileOnboardingPending,
        userProfile,
        handleTextField,
        handleUpdateProfile,
        handleLogout
    };
}