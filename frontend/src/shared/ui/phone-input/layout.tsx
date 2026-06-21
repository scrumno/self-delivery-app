import MaskInput from 'react-native-mask-input'
import { StyleSheet } from 'react-native';

export const PhoneInputMask = {
    RU: ['+', '7', ' ', '(', /\d/, /\d/, /\d/, ')', ' ', /\d/, /\d/, /\d/, '-', /\d/, /\d/, '-', /\d/, /\d/],
};

type PhoneInputProps = {
    value: string, 
    onChangeText: (value: string, rawValue: string) => void
}

export default function PhoneInput({ value, onChangeText }: PhoneInputProps) {
    return (
        <MaskInput
            mask={PhoneInputMask.RU}
            placeholder={'+7'}
            keyboardType="numeric"
            value={value}
            onChangeText={(maskedValue, rawValue) => onChangeText(maskedValue, rawValue)}
            style={[styles.input, styles.defaultInput]}
        />
    )
}

const styles = StyleSheet.create({
    input: {
        padding: 15,
        borderRadius: 10,
        marginBottom: 15,
        alignItems: 'center',
        borderWidth: 1,
    },
    defaultInput: {
        backgroundColor: '#FFF',
        color: '#000',
        fontSize: 18,
        fontWeight: '600',
    }
});