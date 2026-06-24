import ActionButton, { StylePreset } from "@shared/ui/action-button/layout"
import CustomTextInput, { StylePreset as TextStylePreset } from "@shared/ui/text-input/layout"
import { observer } from "mobx-react-lite"
import { ScrollView, StyleSheet, Text, View } from "react-native"
import { usePersonalFields } from "../model/use-personal-field"

const ProfilePreviewPending = observer(() => {
    const { handleTextField, userProfile, handleUpdateProfile } = usePersonalFields();

    return (
        <ScrollView style={styles.scroll}>
            <View style={styles.header}>
                <Text style={styles.title}>Давайте познакомимся!</Text>
                <Text style={styles.subtitle}>
                    Расскажите о себе, чтобы мы могли лучше вас узнать
                </Text>
            </View>

            <View style={styles.buttonContainer}>
                <>
                    <CustomTextInput
                        onChangeText={(text: string) => handleTextField('firstName', text)}
                        value={userProfile.firstName}
                        label="Введите ваше имя"
                        stylePreset={TextStylePreset.DEFAULT}
                    />
                    <CustomTextInput
                        onChangeText={(text: string) => handleTextField('secondName', text)}
                        value={userProfile.secondName}
                        label="Введите вашу фамилию"
                        stylePreset={TextStylePreset.DEFAULT}
                    /> 
                </>
                
                <ActionButton
                    text="Сохранить"
                    onPress={handleUpdateProfile}
                    stylePreset={StylePreset.PRIMARY}
                    style={styles.buttonFullWidth}
                />
            </View>
        </ScrollView>
    )
})

const styles = StyleSheet.create({
    scroll: {
        flex: 1,
    },
    scrollContent: {
        flexGrow: 1,
    },
    container: {
      flex: 1,
      backgroundColor: '#fff',
      padding: 20,
    },
    pinCodeContainer: {
        alignItems: 'center',
        width: '15%'
    },
    buttonFullWidth: {
        width: '100%'
    },
    header: {
      paddingVertical: 24,
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
      justifyContent: 'flex-start',
      alignItems: 'center',
      paddingTop: 8,
      width: '100%',
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

export default ProfilePreviewPending;