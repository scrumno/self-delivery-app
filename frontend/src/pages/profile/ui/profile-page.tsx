import { StyleSheet, View } from 'react-native';
import { Profile } from '@features/profile';

export function ProfilePage() {

  return (
    <View style={styles.container}>
      <Profile />
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