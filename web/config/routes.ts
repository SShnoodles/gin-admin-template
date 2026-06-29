export default [
  {
    path: '/user',
    layout: false,
    routes: [
      {
        path: '/user/login',
    name: '登录',
        component: './user/login',
      },
      {
        path: '/user',
        redirect: '/user/login',
      },
      {
        component: './exception/404',
        path: '/user/*',
      },
    ],
  },
  {
    path: '/',
    redirect: '/home',
  },
  {
    path: '/home',
    name: '首页',
    icon: 'home',
    component: './home',
  },
  {
    path: '/system',
    name: '系统管理',
    icon: 'setting',
    routes: [
      {
        path: '/system',
        redirect: '/system/users',
      },
      {
        path: '/system/users',
        name: '用户管理',
        icon: 'user',
        component: './system/users',
      },
      {
        path: '/system/orgs',
        name: '机构管理',
        icon: 'apartment',
        component: './system/orgs',
      },
      {
        path: '/system/roles',
        name: '角色管理',
        icon: 'team',
        component: './system/roles',
      },
      {
        path: '/system/menus',
        name: '菜单管理',
        icon: 'menu',
        component: './system/menus',
      },
      {
        path: '/system/resources',
        name: '资源管理',
        icon: 'api',
        component: './system/resources',
      },
    ],
  },
  {
    component: './exception/404',
    path: '/*',
  },
];
