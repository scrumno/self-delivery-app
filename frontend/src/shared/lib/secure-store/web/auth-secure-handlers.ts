function readWebUserJson(key: string): string | null {
    if (typeof sessionStorage === 'undefined') {
        return null;
    }
    return sessionStorage.getItem(key);
}

function writeWebUserJson(key: string, json: string): void {
    if (typeof sessionStorage === 'undefined') {
        return;
    }
    sessionStorage.setItem(key, json);
}

function removeWebUserJson(key: string): void {
    if (typeof sessionStorage === 'undefined') {
        return;
    }
    sessionStorage.removeItem(key);
}

export { readWebUserJson, writeWebUserJson, removeWebUserJson };