import {
    View,
    Text,
    StyleSheet,
} from 'react-native';
import { useSms } from "../model/use-sms";
import ActionButton, { StylePreset } from '@shared/ui/action-button/layout';
import PhoneInput from '@shared/ui/phone-input/layout';

const LoginForm = () => {
    const { phone, setPhone, isSmsPending, error, handleRequestSmsCode, isRegister } = useSms();

    const notificationBlock = error && <Text style={styles.errorText}>{error}</Text>;
    const title = isRegister ? 'Зарегистрироваться' : 'Авторизоваться';

    return (
        <>
            <View style={styles.header}>
                <Text style={styles.title}>{title}</Text>
                <Text style={styles.subtitle}>
                    Укажите номер телефона, чтобы продолжить
                </Text>
            </View>

            <View style={styles.buttonContainer}>
                <PhoneInput 
                    value={phone} 
                    onChangeText={setPhone}
                />
                <ActionButton
                    text={title}
                    onPress={handleRequestSmsCode}
                    disabled={isSmsPending}
                    stylePreset={StylePreset.DEFAULT}
                />
            </View>
            
            {notificationBlock}
        </>
    )
}

const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: '#fff',
      padding: 20,
    },
    header: {
      flex: 1,
      justifyContent: 'center',
      alignItems: 'center',
    },
    title: {
      fontSize: 32,
      fontWeight: 'bold',
      color: '#333',
      marginBottom: 10,
    },
    subtitle: {
      fontSize: 16,
      color: '#666',
      textAlign: 'center',
      paddingHorizontal: 30,
    },
    buttonContainer: {
      flex: 1,
      justifyContent: 'flex-end',
      marginBottom: 50,
    },
    button: {
      paddingVertical: 15,
      borderRadius: 10,
      marginBottom: 15,
      alignItems: 'center',
    },
    loginButton: {
      backgroundColor: '#007AFF',
    },
    loginButtonText: {
      color: '#fff',
      fontSize: 18,
      fontWeight: '600',
    },
    registerButton: {
      backgroundColor: 'transparent',
      borderWidth: 2,
      borderColor: '#007AFF',
    },
    registerButtonText: {
      color: '#007AFF',
      fontSize: 18,
      fontWeight: '600',
    },
    errorText: {
      color: 'red',
      fontSize: 16,
      textAlign: 'center',
      marginBottom: 10,
    },
});

export default LoginForm;