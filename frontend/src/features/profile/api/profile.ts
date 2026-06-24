import api from "@shared/api/api";
import { UserFullProfile } from "../model/type";

export interface UserProfile {
    fullName?: string,
    birthDate?: string,
    isActive?: boolean,
    Email?: string,
}

export interface UserProfileResponse {
    isSuccess: boolean,
    error?: string,
}

export const updateProfile = async (data: Partial<UserFullProfile>): Promise<UserProfileResponse> => {
    try {
        const response = await api.put('/users/update-user-profile', data);
        return response.data;
    } catch (error: any) {
        if (!error.response) {
            throw new Error('Произошла неизвестная ошибка');
        }

        return error.response.data;        
    }
}

