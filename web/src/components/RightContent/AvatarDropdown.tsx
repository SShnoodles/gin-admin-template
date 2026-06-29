import { LogoutOutlined, SkinOutlined } from '@ant-design/icons';
import { history, useModel } from '@umijs/max';
import type { MenuProps } from 'antd';
import { Space, Spin } from 'antd';
import { createStyles } from 'antd-style';
import React, { startTransition } from 'react';
import { clearSession, getCurrentUser } from '@/services/admin/session';
import HeaderDropdown from '../HeaderDropdown';

type GlobalHeaderRightProps = {
  children?: React.ReactNode;
};

const useStyles = createStyles(({ token }) => ({
  account: {
    display: 'inline-flex',
    alignItems: 'center',
    height: '100%',
    paddingInline: 12,
    cursor: 'pointer',
    borderRadius: token.borderRadius,
    color: token.colorText,
    transition: `background ${token.motionDurationMid}`,
    '&:hover': {
      background: token.colorBgTextHover,
    },
  },
  name: {
    maxWidth: 120,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
}));

const menuItems: MenuProps['items'] = [
  {
    key: 'theme',
    icon: <SkinOutlined />,
    label: '主题设置',
  },
  {
    type: 'divider' as const,
  },
  {
    key: 'logout',
    icon: <LogoutOutlined />,
    label: '退出登录',
  },
];

const loginOut = () => {
  clearSession();
  const { search, pathname } = window.location;
  const urlParams = new URL(window.location.href).searchParams;
  const searchParams = new URLSearchParams({
    redirect: pathname + search,
  });
  const redirect = urlParams.get('redirect');
  if (window.location.pathname !== '/user/login' && !redirect) {
    history.replace({
      pathname: '/user/login',
      search: searchParams.toString(),
    });
  }
};

export const AvatarDropdown: React.FC<GlobalHeaderRightProps> = ({
  children,
}) => {
  const { initialState, setInitialState } = useModel('@@initialState');
  const { styles } = useStyles();

  const onMenuClick: MenuProps['onClick'] = (event) => {
    const { key } = event;
    if (key === 'logout') {
      startTransition(() => {
        setInitialState((s) => ({ ...s, currentUser: undefined }));
      });
      loginOut();
      return;
    }
    if (key === 'theme') {
      setInitialState((s) => ({ ...s, settingDrawerOpen: true }));
      return;
    }
  };

  if (!initialState) {
    return <Spin size="small" />;
  }

  const currentUser = initialState.currentUser ?? getCurrentUser();

  if (!currentUser) {
    return <Spin size="small" />;
  }

  return (
    <HeaderDropdown
      placement="bottomRight"
      menu={{
        selectedKeys: [],
        onClick: onMenuClick,
        items: menuItems,
      }}
      arrow
    >
      <span className={styles.account}>
        <Space size={8}>
          {children}
          <span className={styles.name}>{currentUser.name || 'Admin'}</span>
        </Space>
      </span>
    </HeaderDropdown>
  );
};
