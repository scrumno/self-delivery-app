import { StyleSheet, View } from "react-native"
import { useLogin } from "../model/use-login";
import { OtpInput } from "react-native-otp-entry";
import { Text } from "react-native";

const CodeForm = () => {
    const { requestHandler, error, setCode } = useLogin();
    const notificationBlock = error && <Text style={styles.errorText}>{error}</Text>;

    return (
        <>
            <View style={styles.header}>
                <Text style={styles.title}>Введите код</Text>
                <Text style={styles.subtitle}>
                    Введите код, который был отправлен на ваш номер телефона
                </Text>
            </View>

            <View style={styles.buttonContainer}>
                <OtpInput
                    numberOfDigits={4}
                    type="numeric"
                    onTextChange={setCode}
                    onFilled={requestHandler}
                    theme={{
                        containerStyle: styles.pinCodeContainer,
                        pinCodeContainerStyle: styles.pinCodeContainer,
                    }}
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
    pinCodeContainer: {
        alignItems: 'center',
        width: '15%'
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
      justifyContent: 'center',
      alignItems: 'center',
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


export default CodeForm;