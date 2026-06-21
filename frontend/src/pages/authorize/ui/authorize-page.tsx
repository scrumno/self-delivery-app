import { StyleSheet, View } from 'react-native';
import LoginForm from '@features/authorize/ui/login-form';

export function AuthorizePage() {

  return (
    <View style={styles.container}>
      <LoginForm />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    padding: 20,
  },
});