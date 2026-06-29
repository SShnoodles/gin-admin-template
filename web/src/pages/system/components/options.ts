import type { DataNode } from 'antd/es/tree';
import { getMenus, getOrgRoles, getOrgs, getResources, getRoles } from '@/services/admin';
import type { Id, Menu, Org, Resource, Role } from '@/services/admin/types';

export type SelectOption = {
  label: string;
  value: Id;
};

export const normalizeList = <T,>(result: T[] | { data?: T[] }) => {
  return Array.isArray(result) ? result : result.data || [];
};

export const toMenuTreeData = (menus: Menu[] = []): DataNode[] => {
  return menus.map((menu) => ({
    title: menu.name || menu.path || menu.id,
    key: String(menu.id),
    value: String(menu.id),
    children: toMenuTreeData(menu.children || []),
  }));
};

export const flattenMenus = (menus: Menu[] = []): Menu[] => {
  return menus.flatMap((menu) => [menu, ...flattenMenus(menu.children || [])]);
};

export async function loadOrgOptions(): Promise<SelectOption[]> {
  const result = await getOrgs();
  return normalizeList<Org>(result).map((item) => ({
    label: item.name || String(item.id),
    value: String(item.id),
  }));
}

export async function loadRoleOptions(orgId?: Id): Promise<SelectOption[]> {
  const result = orgId ? await getOrgRoles(orgId) : await getRoles({ pageSize: 1000 });
  return normalizeList<Role>(result).map((item) => ({
    label: item.name || item.code || String(item.id),
    value: String(item.id),
  }));
}

export async function loadResourceOptions(): Promise<SelectOption[]> {
  const result = await getResources();
  return normalizeList<Resource>(result).map((item) => ({
    label: `${item.name || item.path} ${item.method ? `(${item.method})` : ''}`,
    value: String(item.id),
  }));
}

export async function loadMenuTreeData() {
  const menus = await getMenus();
  return toMenuTreeData(menus);
}

export async function loadParentMenuOptions(currentId?: Id): Promise<SelectOption[]> {
  const menus = flattenMenus(await getMenus()).filter(
    (menu) => String(menu.id) !== String(currentId),
  );
  return [
    { label: '顶级菜单', value: '0' },
    ...menus.map((menu) => ({
      label: menu.name || menu.path || String(menu.id),
      value: String(menu.id),
    })),
  ];
}
