import { useMemo } from "react";
import { StyleSheet, Text, TouchableOpacity, TouchableOpacityProps } from "react-native";

export enum StylePreset {
    DEFAULT = 'default',
    PRIMARY = 'primary',
    SECONDARY = 'secondary',
    SUCCESS = 'success',
    DANGER = 'danger',
    WARNING = 'warning',
    INFO = 'info',
}

type ActionButtonProps = TouchableOpacityProps & {
    text: string;
    stylePreset: StylePreset;
};

export default function ActionButton({ text, stylePreset, style, ...touchableProps }: ActionButtonProps) {

    const presetStyle = useMemo(() => {
        switch (stylePreset) {
            case StylePreset.DEFAULT:
                return {
                    textColor: styles.registerButtonText,
                    buttonColor: styles.registerButton,
                };
            case StylePreset.PRIMARY:
                return {
                    textColor: styles.loginButtonText,
                    buttonColor: styles.loginButton,
                };
            default:
                return {
                    textColor: styles.loginButtonText,
                    buttonColor: styles.loginButton,
                }
        }
    }, [stylePreset]);

    return (
        <TouchableOpacity {...touchableProps} style={[presetStyle.buttonColor, styles.button, style]}>
            <Text style={presetStyle.textColor}>{text}</Text>
        </TouchableOpacity>
    )
}

const styles = StyleSheet.create({
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
});