export type PageParams = {
  current?: number;
  pageIndex?: number;
  pageSize?: number;
};

export type PageResult<T> = {
  data?: T[];
  total?: number;
};

export type Id = string | number;

export type Org = {
  id?: Id;
  name?: string;
  creditCode?: string;
  address?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type Role = {
  id?: Id;
  name?: string;
  code?: string;
  orgId?: Id;
  orgName?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type User = {
  id?: Id;
  username?: string;
  realName?: string;
  workNo?: string;
  orgId?: Id;
  orgName?: string;
  enabled?: boolean;
  createdAt?: string;
  updatedAt?: string;
  roleIds?: Id[];
};

export type Menu = {
  id?: Id;
  pid?: Id;
  name?: string;
  path?: string;
  icon?: string;
  sort?: number;
  children?: Menu[];
  createdAt?: string;
  updatedAt?: string;
  resourceIds?: Id[];
};

export type Resource = {
  id?: Id;
  name?: string;
  method?: string;
  path?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type LoginParams = {
  username?: string;
  password?: string;
  codeId?: string;
  code?: string;
};

export type LoginResult = {
  accessToken?: string;
  expires?: string;
  refreshToken?: string;
};

export type CaptchaResult = {
  codeId?: string;
  code?: string;
};
