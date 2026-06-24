import { Onest_300Light } from '@expo-google-fonts/onest/300Light';
import { Onest_400Regular } from '@expo-google-fonts/onest/400Regular';
import { Onest_500Medium } from '@expo-google-fonts/onest/500Medium';
import { Onest_600SemiBold } from '@expo-google-fonts/onest/600SemiBold';
import { Onest_700Bold } from '@expo-google-fonts/onest/700Bold';
import { useFonts } from '@expo-google-fonts/onest/useFonts';

export const onestFonts = {
  Onest_300Light,
  Onest_400Regular,
  Onest_500Medium,
  Onest_600SemiBold,
  Onest_700Bold,
} as const;

export function useOnestFonts() {
  return useFonts(onestFonts);
}
