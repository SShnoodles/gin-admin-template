import { request } from '@umijs/max';
import type { RequestOptions } from '@@/plugin-request/request';
import type {
  ApiResult,
  CaptchaResult,
  Id,
  LoginParams,
  LoginResult,
  Menu,
  Org,
  PageParams,
  PageResult,
  Resource,
  Role,
  User,
} from './types';

const isApiResult = <T,>(value: unknown): value is ApiResult<T> => {
  return Boolean(value && typeof value === 'object' && 'success' in value);
};

async function requestData<T>(url: string): Promise<T>;
async function requestData<T>(url: string, options: RequestOptions | undefined): Promise<T>;
async function requestData<T>(url: string, options: RequestOptions): Promise<T>;
async function requestData<T>(url: string, options?: RequestOptions): Promise<T> {
  const result = options
    ? await request<T | ApiResult<T>>(url, options)
    : await request<T | ApiResult<T>>(url);
  if (isApiResult<T>(result)) {
    return result.data as T;
  }
  return result;
}

async function requestArray<T>(url: string, options?: RequestOptions): Promise<T[]> {
  const result = await requestData<T[] | undefined>(url, options);
  return Array.isArray(result) ? result : [];
}

async function requestPage<T>(url: string, options?: RequestOptions): Promise<PageResult<T>> {
  const result = await requestData<PageResult<T> | undefined>(url, options);
  return {
    data: result?.data || [],
    total: result?.total || 0,
  };
}

async function requestListOrPage<T>(
  url: string,
  options?: RequestOptions,
): Promise<PageResult<T> | T[]> {
  const result = await requestData<PageResult<T> | T[] | undefined>(url, options);
  if (Array.isArray(result)) {
    return result;
  }
  return {
    data: result?.data || [],
    total: result?.total || 0,
  };
}

const pageParams = (params?: PageParams & Record<string, unknown>) => {
  const { current, pageIndex, ...rest } = params || {};
  return {
    pageIndex: pageIndex ?? current,
    ...rest,
  };
};

export async function login(body: LoginParams) {
  return requestData<LoginResult>('/login/account', {
    method: 'POST',
    data: body,
  });
}

export async function getCaptcha() {
  return requestData<CaptchaResult>('/login/captcha', {
    method: 'POST',
  });
}

export async function getUsers(params?: PageParams & Partial<User>) {
  return requestPage<User>('/users', {
    method: 'GET',
    params: pageParams(params),
  });
}

export async function getUser(id: Id) {
  return requestData<User>(`/users/${id}`, {
    method: 'GET',
  });
}

export async function getUserRoles(id: Id) {
  return requestArray<Id>(`/users/${id}/roles`, {
    method: 'GET',
  });
}

export async function createUser(body: User) {
  return requestData<{ id: Id }>('/users', {
    method: 'POST',
    data: body,
  });
}

export async function updateUser(id: Id, body: User) {
  return requestData<{ message: string }>(`/users/${id}`, {
    method: 'PUT',
    data: body,
  });
}

export async function deleteUser(id: Id) {
  return requestData<{ message: string }>(`/users/${id}`, {
    method: 'DELETE',
  });
}

export async function toggleUserEnabled(id: Id) {
  return requestData<{ message: string }>(`/users/${id}/enabled`, {
    method: 'PUT',
  });
}

export async function changePassword(body: {
  oldPassword?: string;
  newPassword?: string;
}) {
  return requestData<{ message: string }>('/users/change-password', {
    method: 'PUT',
    data: body,
  });
}

export async function getOrgs(params?: PageParams & Partial<Org>) {
  return requestListOrPage<Org>('/orgs', {
    method: 'GET',
    params: pageParams(params),
  });
}

export async function getOrg(id: Id) {
  return requestData<Org>(`/orgs/${id}`, {
    method: 'GET',
  });
}

export async function getOrgMenus(id: Id) {
  return requestArray<Id>(`/orgs/${id}/menus`, {
    method: 'GET',
  });
}

export async function createOrg(body: Org & { menuIds?: Id[] }) {
  return requestData<{ id: Id }>('/orgs', {
    method: 'POST',
    data: body,
  });
}

export async function updateOrg(id: Id, body: Org & { menuIds?: Id[] }) {
  return requestData<{ message: string }>(`/orgs/${id}`, {
    method: 'PUT',
    data: body,
  });
}

export async function deleteOrg(id: Id) {
  return requestData<{ message: string }>(`/orgs/${id}`, {
    method: 'DELETE',
  });
}

export async function getRoles(params?: PageParams & Partial<Role>) {
  return requestPage<Role>('/roles', {
    method: 'GET',
    params: pageParams(params),
  });
}

export async function getOrgRoles(orgId: Id) {
  return requestArray<Role>(`/roles/orgs/${orgId}`, {
    method: 'GET',
  });
}

export async function getRole(id: Id) {
  return requestData<Role>(`/roles/${id}`, {
    method: 'GET',
  });
}

export async function getRoleMenus(id: Id) {
  return requestArray<Id>(`/roles/${id}/menus`, {
    method: 'GET',
  });
}

export async function createRole(body: Role & { menuIds?: Id[] }) {
  return requestData<{ id: Id }>('/roles', {
    method: 'POST',
    data: body,
  });
}

export async function updateRole(id: Id, body: Role & { menuIds?: Id[] }) {
  return requestData<{ message: string }>(`/roles/${id}`, {
    method: 'PUT',
    data: body,
  });
}

export async function deleteRole(id: Id) {
  return requestData<{ message: string }>(`/roles/${id}`, {
    method: 'DELETE',
  });
}

export async function getMenus() {
  return requestArray<Menu>('/menus', {
    method: 'GET',
  });
}

export async function getMenu(id: Id) {
  return requestData<Menu>(`/menus/${id}`, {
    method: 'GET',
  });
}

export async function getMenuResources(id: Id) {
  return requestArray<Id>(`/menus/${id}/resources`, {
    method: 'GET',
  });
}

export async function createMenu(body: Menu & { resourceIds?: Id[] }) {
  return requestData<{ id: Id }>('/menus', {
    method: 'POST',
    data: body,
  });
}

export async function updateMenu(id: Id, body: Menu & { resourceIds?: Id[] }) {
  return requestData<{ message: string }>(`/menus/${id}`, {
    method: 'PUT',
    data: body,
  });
}

export async function deleteMenu(id: Id) {
  return requestData<{ message: string }>(`/menus/${id}`, {
    method: 'DELETE',
  });
}

export async function getResources(params?: PageParams & Partial<Resource>) {
  return requestListOrPage<Resource>('/resources', {
    method: 'GET',
    params: pageParams(params),
  });
}
