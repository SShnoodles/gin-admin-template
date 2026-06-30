import type { TreeSelectProps } from 'antd';
import { getMenus, getOrgRoles, getOrgs, getResources, getRoles } from '@/services/admin';
import type { Id, Menu, Org, Resource, Role } from '@/services/admin/types';

type TreeSelectDataNode = NonNullable<TreeSelectProps['treeData']>[number];

export type SelectOption = {
  label: string;
  value: Id;
};

export const normalizeList = <T,>(result: T[] | { data?: T[] }) => {
  return Array.isArray(result) ? result : result.data || [];
};

export const toMenuTreeData = (menus: Menu[] = []): TreeSelectDataNode[] => {
  return menus.map((menu) => ({
    title: menu.name || menu.path || menu.id,
    label: menu.name || menu.path || menu.id,
    key: String(menu.id),
    value: String(menu.id),
    children: toMenuTreeData(menu.children || []),
  }));
};

export const toParentMenuTreeData = (menus: Menu[] = [], currentId?: Id): TreeSelectDataNode[] => {
  const buildNode = (menu: Menu): TreeSelectDataNode | null => {
    if (currentId && String(menu.id) === String(currentId)) {
      return null;
    }

    return {
      title: menu.name || menu.path || menu.id,
      label: menu.name || menu.path || menu.id,
      key: String(menu.id),
      value: String(menu.id),
      children: (menu.children || [])
        .map(buildNode)
        .filter((item): item is TreeSelectDataNode => Boolean(item)),
    };
  };

  return [
    {
      title: '顶级菜单',
      label: '顶级菜单',
      key: '0',
      value: '0',
      isLeaf: true,
    },
    ...menus.map(buildNode).filter((item): item is TreeSelectDataNode => Boolean(item)),
  ];
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

export async function loadParentMenuTreeData(currentId?: Id) {
  const menus = await getMenus();
  return toParentMenuTreeData(menus, currentId);
}
