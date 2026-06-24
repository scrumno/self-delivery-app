import ProfilePreviewPending from "@features/profile/ui/profile-preview-pending"
import { useAuth } from "provider/auth-provider"
import ProfilePreview from "@features/profile/ui/profile-preview";

export const ProfileWidget = () => {
    const { profileOnboardingPending } = useAuth();

    if (profileOnboardingPending) {
        return <ProfilePreviewPending />
    }

    return <ProfilePreview />
}