import type { LoginResult } from './types';

const TOKEN_KEY = 'gin_admin_access_token';
const EXPIRES_KEY = 'gin_admin_access_token_expires';
const USERNAME_KEY = 'gin_admin_username';

export function saveSession(username: string | undefined, result: LoginResult) {
  if (result.accessToken) {
    localStorage.setItem(TOKEN_KEY, result.accessToken);
  }
  if (result.expires) {
    localStorage.setItem(EXPIRES_KEY, result.expires);
  }
  if (username) {
    localStorage.setItem(USERNAME_KEY, username);
  }
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(EXPIRES_KEY);
  localStorage.removeItem(USERNAME_KEY);
}

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function getCurrentUser(): API.CurrentUser | undefined {
  const token = getToken();
  if (!token) return undefined;

  return {
    name: localStorage.getItem(USERNAME_KEY) || 'Admin',
    access: 'admin',
  };
}
