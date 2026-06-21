import { Label } from "@react-navigation/elements";
import { useMemo } from "react";
import { StyleSheet, TextInput, View } from "react-native";

export enum StylePreset {
    DEFAULT = 'default',
}

type CustomTextInputProps = {
    label: string;
    value: string;
    onChangeText: (text: string) => void;
    stylePreset: StylePreset;
};

export default function CustomTextInput({ stylePreset, value, label, onChangeText }: CustomTextInputProps) {
    switch (stylePreset) {
        case StylePreset.DEFAULT:
            return (
                <View style={styles.buttonWrapper}>
                    <Label style={styles.registerButtonText}>{label}</Label>
                    <TextInput style={[styles.button, styles.registerButton]} value={value} onChangeText={onChangeText}></TextInput>
                </View>
            );
        default:
            return (
                <View style={styles.buttonWrapper}>
                    <Label style={styles.registerButtonText}>{label}</Label>
                    <TextInput style={[styles.button, styles.registerButton]} value={value} onChangeText={onChangeText}></TextInput>
                </View>
            );
    }
}

const styles = StyleSheet.create({
    button: {
        minHeight: 48,
        padding: 15,
        borderRadius: 10,
        marginBottom: 15,
        alignItems: 'center',
        width: '100%'
    },
    buttonWrapper: {
        width: '100%',
        alignItems: 'flex-start',
        gap: '8px'
    },
    label: {
        flex: 1,
        display: 'flex',
        color: '#000',
        fontSize: 18,
        fontWeight: '600',
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
});