import { Alert, Platform } from "react-native";
import { ErrorCode } from "./variables";

export const handleUserNotFound = (phone: string, successCallback: () => void, unsuccessCallback: () => void) => {
    if (Platform.OS === 'web') {
        const shouldRegister = window.confirm(
            `Пользователь с номером ${phone} не найден. Хотите зарегистрироваться?`
        );
        
        if (shouldRegister) {
            successCallback();
        }

        unsuccessCallback();

        return;
    }
    
    Alert.alert('Пользователь не найден', 'Аккаунта с таким номером нет. Хотите зарегистрироваться?', [{
        text: 'Ввести другой номер', 
        style: 'cancel', 
        onPress: unsuccessCallback
    }, { 
        text: 'Зарегистрироваться',
        onPress: successCallback
    }]);
}

export const handleUserExist = (phone: string, successCallback: () => void, unsuccessCallback: () => void) => {
    if (Platform.OS === 'web') {
        const shouldRegister = window.confirm(
            `Пользователь с номером ${phone} найден. Хотите авторизоваться?`
        );
        
        if (shouldRegister) {
            successCallback();
        }

        unsuccessCallback();

        return;
    }
    
    Alert.alert('Пользователь найден', 'Аккаунта с таким номером уже существует. Хотите авторизоваться?', [{
        text: 'Ввести другой номер', 
        style: 'cancel', 
        onPress: unsuccessCallback
    }, { 
        text: 'Авторизоваться',
        onPress: successCallback
    }]);
}

export const handleError = (phone: string, successCallback: () => void, unsuccessCallback: () => void, code?: string) => {
    switch (code) {
        case ErrorCode.USER_NOT_FOUND:
            handleUserNotFound(phone, successCallback, unsuccessCallback)
            break;
        case ErrorCode.USER_EXIST:
            handleUserExist(phone, successCallback, unsuccessCallback)
            break;
        default:
            return;
    }
}