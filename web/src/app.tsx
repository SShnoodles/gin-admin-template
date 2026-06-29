import type { Settings as LayoutSettings } from '@ant-design/pro-components';
import { SettingDrawer } from '@ant-design/pro-components';
import type { RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import type { RequestOptions } from '@@/plugin-request/request';
import { UserOutlined } from '@ant-design/icons';
import { history, Link } from '@umijs/max';
import { Avatar } from 'antd';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import React from 'react';
import { AvatarDropdown, ErrorBoundary, Footer, OfflineBanner } from '@/components';
import { clearSession, getCurrentUser, getToken } from '@/services/admin/session';
import defaultSettings from '../config/defaultSettings';
import { errorConfig } from './requestErrorConfig';

dayjs.extend(relativeTime);

const loginPath = '/user/login';

export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: API.CurrentUser;
  loading?: boolean;
  fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
  settingDrawerOpen?: boolean;
}> {
  const fetchUserInfo = async () => getCurrentUser();
  const currentUser = getCurrentUser();
  const { location } = history;

  if (!currentUser && location.pathname !== loginPath) {
    history.replace(
      `${loginPath}?redirect=${encodeURIComponent(
        location.pathname + location.search + location.hash,
      )}`,
    );
  }

  return {
    fetchUserInfo,
    currentUser,
    settings: defaultSettings as Partial<LayoutSettings>,
    settingDrawerOpen: false,
  };
}

export const layout: RunTimeLayoutConfig = ({ initialState, setInitialState }) => {
  const currentUser = initialState?.currentUser ?? getCurrentUser();

  return {
    menuItemRender: (item, dom) => {
      if (item.path) {
        return (
          <Link to={item.path} prefetch>
            {dom}
          </Link>
        );
      }
      return dom;
    },
    actionsRender: false,
    footerRender: () => <Footer />,
    onPageChange: () => {
      const { location } = history;
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.replace(
          `${loginPath}?redirect=${encodeURIComponent(
            location.pathname + location.search + location.hash,
          )}`,
        );
      }
    },
    ErrorBoundary,
    menuHeaderRender: undefined,
    childrenRender: (children) => (
      <>
        {children}
        {currentUser ? (
          <div
            style={{
              position: 'fixed',
              top: 8,
              right: 24,
              zIndex: 1001,
            }}
          >
            <AvatarDropdown>
              <Avatar size="small" icon={<UserOutlined />} />
            </AvatarDropdown>
          </div>
        ) : null}
        <SettingDrawer
          disableUrlParams
          enableDarkTheme
          collapse={initialState?.settingDrawerOpen}
          onCollapseChange={(open) => {
            setInitialState((s) => ({
              ...s,
              settingDrawerOpen: open,
            }));
          }}
          settings={initialState?.settings}
          onSettingChange={(settings) => {
            setInitialState((s) => ({
              ...s,
              settings,
            }));
          }}
        />
      </>
    ),
    ...initialState?.settings,
  };
};

export const request: RequestConfig = {
  baseURL: '/api',
  ...errorConfig,
  requestInterceptors: [
    (config: RequestOptions) => {
      const token = getToken();
      if (token) {
        config.headers = {
          ...config.headers,
          Authorization: token.startsWith('Bearer ') ? token : `Bearer ${token}`,
        };
      }
      return config;
    },
  ],
  responseInterceptors: [
    (response) => {
      if (response.status === 401) {
        clearSession();
        history.replace(loginPath);
      }
      return response;
    },
  ],
};

export function rootContainer(container: React.ReactNode) {
  return (
    <>
      <OfflineBanner />
      <ErrorBoundary>{container}</ErrorBoundary>
    </>
  );
}
