export type UserBaseProfile = {
    firstName: string,
    secondName: string,
    birthDate: Date | null
};

export type UserFullProfile = UserBaseProfile & {
    address: string | null,
    email: string | null,
    fullName: string
}