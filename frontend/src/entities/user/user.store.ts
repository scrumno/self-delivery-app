import { makeAutoObservable } from 'mobx';
import { UserNames } from './types';

class UserStore {

    firstName: string = '';
    secondName: string = '';
    birthDate: Date | null = null;
    address: string | null = null;
    email: string | null = null;

    constructor() {
        makeAutoObservable(this);
    }

    setNames(firstName: string, secondName: string) {
        this.firstName = firstName;
        this.secondName = secondName;
    }

    hydrate(partial: Partial<UserNames>) {
        if (partial.firstName !== undefined) {
            this.firstName = partial.firstName;
        }
        if (partial.secondName !== undefined) {
            this.secondName = partial.secondName;
        }
    }

    reset() {
        this.firstName = '';
        this.secondName = '';
        this.birthDate = null;
        this.address = null;
        this.email = null;
    }
    
    get fullName(): string {
        return [this.firstName, this.secondName]
            .map((s) => s.trim())
            .filter(Boolean)
            .join(' ');
    }
}

export const userStore = new UserStore();
