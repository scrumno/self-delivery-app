import { StyleSheet, View } from 'react-native';
import { CodeForm } from '@features/authorize';

export function CodePage() {

  return (
    <View style={styles.container}>
      <CodeForm />
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
