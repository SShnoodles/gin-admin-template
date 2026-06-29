import { request } from '@umijs/max';
import type {
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

const pageParams = (params?: PageParams & Record<string, unknown>) => {
  const { current, pageIndex, ...rest } = params || {};
  return {
    pageIndex: pageIndex ?? current,
    ...rest,
  };
};

export async function login(body: LoginParams) {
  return request<LoginResult>('/login/account', {
    method: 'POST',
    data: body,
  });
}

export async function getCaptcha() {
  return request<CaptchaResult>('/login/captcha', {
    method: 'POST',
  });
}

export async function getUsers(params?: PageParams & Partial<User>) {
  return request<PageResult<User>>('/users', {
    method: 'GET',
    params: pageParams(params),
  });
}

export async function getUser(id: Id) {
  return request<User>(`/users/${id}`, {
    method: 'GET',
  });
}

export async function getUserRoles(id: Id) {
  return request<Id[]>(`/users/${id}/roles`, {
    method: 'GET',
  });
}

export async function createUser(body: User) {
  return request<{ id: Id }>('/users', {
    method: 'POST',
    data: body,
  });
}

export async function updateUser(id: Id, body: User) {
  return request<{ message: string }>(`/users/${id}`, {
    method: 'PUT',
    data: body,
  });
}

export async function deleteUser(id: Id) {
  return request<{ message: string }>(`/users/${id}`, {
    method: 'DELETE',
  });
}

export async function toggleUserEnabled(id: Id) {
  return request<{ message: string }>(`/users/${id}/enabled`, {
    method: 'PUT',
  });
}

export async function changePassword(body: {
  oldPassword?: string;
  newPassword?: string;
}) {
  return request<{ message: string }>('/users/change-password', {
    method: 'PUT',
    data: body,
  });
}

export async function getOrgs(params?: PageParams & Partial<Org>) {
  return request<PageResult<Org> | Org[]>('/orgs', {
    method: 'GET',
    params: pageParams(params),
  });
}

export async function getOrg(id: Id) {
  return request<Org>(`/orgs/${id}`, {
    method: 'GET',
  });
}

export async function getOrgMenus(id: Id) {
  return request<Id[]>(`/orgs/${id}/menus`, {
    method: 'GET',
  });
}

export async function createOrg(body: Org & { menuIds?: Id[] }) {
  return request<{ id: Id }>('/orgs', {
    method: 'POST',
    data: body,
  });
}

export async function updateOrg(id: Id, body: Org & { menuIds?: Id[] }) {
  return request<{ message: string }>(`/orgs/${id}`, {
    method: 'PUT',
    data: body,
  });
}

export async function deleteOrg(id: Id) {
  return request<{ message: string }>(`/orgs/${id}`, {
    method: 'DELETE',
  });
}

export async function getRoles(params?: PageParams & Partial<Role>) {
  return request<PageResult<Role>>('/roles', {
    method: 'GET',
    params: pageParams(params),
  });
}

export async function getOrgRoles(orgId: Id) {
  return request<Role[]>(`/roles/orgs/${orgId}`, {
    method: 'GET',
  });
}

export async function getRole(id: Id) {
  return request<Role>(`/roles/${id}`, {
    method: 'GET',
  });
}

export async function getRoleMenus(id: Id) {
  return request<Id[]>(`/roles/${id}/menus`, {
    method: 'GET',
  });
}

export async function createRole(body: Role & { menuIds?: Id[] }) {
  return request<{ id: Id }>('/roles', {
    method: 'POST',
    data: body,
  });
}

export async function updateRole(id: Id, body: Role & { menuIds?: Id[] }) {
  return request<{ message: string }>(`/roles/${id}`, {
    method: 'PUT',
    data: body,
  });
}

export async function deleteRole(id: Id) {
  return request<{ message: string }>(`/roles/${id}`, {
    method: 'DELETE',
  });
}

export async function getMenus() {
  return request<Menu[]>('/menus', {
    method: 'GET',
  });
}

export async function getMenu(id: Id) {
  return request<Menu>(`/menus/${id}`, {
    method: 'GET',
  });
}

export async function getMenuResources(id: Id) {
  return request<Id[]>(`/menus/${id}/resources`, {
    method: 'GET',
  });
}

export async function createMenu(body: Menu & { resourceIds?: Id[] }) {
  return request<{ id: Id }>('/menus', {
    method: 'POST',
    data: body,
  });
}

export async function updateMenu(id: Id, body: Menu & { resourceIds?: Id[] }) {
  return request<{ message: string }>(`/menus/${id}`, {
    method: 'PUT',
    data: body,
  });
}

export async function deleteMenu(id: Id) {
  return request<{ message: string }>(`/menus/${id}`, {
    method: 'DELETE',
  });
}

export async function getResources(params?: PageParams & Partial<Resource>) {
  return request<PageResult<Resource> | Resource[]>('/resources', {
    method: 'GET',
    params: pageParams(params),
  });
}
