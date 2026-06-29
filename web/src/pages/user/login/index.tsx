import { LockOutlined, SafetyCertificateOutlined, UserOutlined } from '@ant-design/icons';
import { LoginForm, ProFormText } from '@ant-design/pro-components';
import { Helmet, history, useModel } from '@umijs/max';
import { Alert, App, Button } from 'antd';
import { createStyles } from 'antd-style';
import React, { startTransition, useEffect, useState } from 'react';
import { Footer } from '@/components';
import { getCaptcha, login } from '@/services/admin';
import { getCurrentUser, saveSession } from '@/services/admin/session';
import type { CaptchaResult, LoginParams, LoginResult } from '@/services/admin/types';
import Settings from '../../../../config/defaultSettings';

const getSafeRedirectUrl = (redirect: string | null): string => {
  if (!redirect?.startsWith('/') || redirect.startsWith('//')) return '/';

  try {
    const parsed = new URL(redirect, window.location.origin);
    if (parsed.origin !== window.location.origin) return '/';
    return `${parsed.pathname}${parsed.search}${parsed.hash}`;
  } catch {
    return '/';
  }
};

const getLoginResult = (result: LoginResult | { data?: LoginResult }): LoginResult => {
  if ('data' in result && result.data) return result.data;
  return result as LoginResult;
};

const getRedirectTarget = () => {
  const urlParams = new URL(window.location.href).searchParams;
  const target = getSafeRedirectUrl(urlParams.get('redirect'));
  return target.startsWith('/user') ? '/home' : target || '/home';
};

const useStyles = createStyles(({ token }) => ({
  container: {
    display: 'flex',
    flexDirection: 'column',
    minHeight: '100vh',
    overflow: 'auto',
    background:
      'linear-gradient(135deg, rgba(22,119,255,0.10), rgba(255,255,255,0.92))',
  },
  captcha: {
    width: 140,
    height: 48,
    borderRadius: token.borderRadius,
    cursor: 'pointer',
    objectFit: 'cover',
    border: `1px solid ${token.colorBorder}`,
  },
}));

const LoginMessage: React.FC<{ content: string }> = ({ content }) => (
  <Alert style={{ marginBottom: 24 }} message={content} type="error" showIcon />
);

const Login: React.FC = () => {
  const [captcha, setCaptcha] = useState<CaptchaResult>();
  const [error, setError] = useState<string>();
  const { initialState, setInitialState } = useModel('@@initialState');
  const { styles } = useStyles();
  const { message } = App.useApp();

  const refreshCaptcha = async () => {
    const result = await getCaptcha();
    setCaptcha(result);
  };

  useEffect(() => {
    refreshCaptcha();
  }, []);

  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      startTransition(() => {
        setInitialState((s) => ({
          ...s,
          currentUser: userInfo,
        }));
      });
    }
  };

  const handleSubmit = async (values: LoginParams) => {
    try {
      const result = getLoginResult(await login({ ...values, codeId: captcha?.codeId }));
      if (!result.accessToken) {
        throw new Error('登录成功但未返回 accessToken');
      }
      saveSession(values.username, result);
      const currentUser = getCurrentUser();
      message.success('登录成功');
      if (currentUser) {
        startTransition(() => {
          setInitialState((s) => ({
            ...s,
            currentUser,
          }));
        });
      } else {
        await fetchUserInfo();
      }
      history.replace(getRedirectTarget());
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败，请重试');
      refreshCaptcha();
    }
  };

  return (
    <div className={styles.container}>
      <Helmet>
        <title>登录 - {Settings.title}</title>
      </Helmet>
      <div style={{ flex: 1, padding: '56px 0' }}>
        <LoginForm<LoginParams>
          contentStyle={{ minWidth: 300, maxWidth: 420 }}
          logo={<img alt="logo" src="/logo.svg" />}
          title={Settings.title}
          subTitle="Gin 后端管理系统"
          onFinish={handleSubmit}
        >
          {error && <LoginMessage content={error} />}
          <ProFormText
            name="username"
            fieldProps={{ size: 'large', prefix: <UserOutlined /> }}
            placeholder="请输入用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          />
          <ProFormText.Password
            name="password"
            fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
            placeholder="请输入密码"
            rules={[{ required: true, message: '请输入密码' }]}
          />
          <ProFormText
            name="code"
            fieldProps={{
              size: 'large',
              prefix: <SafetyCertificateOutlined />,
              suffix: captcha?.code ? (
                <img
                  className={styles.captcha}
                  src={captcha.code}
                  alt="验证码"
                  onClick={refreshCaptcha}
                />
              ) : (
                <Button type="link" onClick={refreshCaptcha}>
                  刷新
                </Button>
              ),
            }}
            placeholder="请输入验证码"
            rules={[{ required: true, message: '请输入验证码' }]}
          />
        </LoginForm>
      </div>
      <Footer />
    </div>
  );
};

export default Login;
