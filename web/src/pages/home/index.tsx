import { PageContainer, ProCard, StatisticCard } from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Button, Space, Typography } from 'antd';

const HomePage = () => {
  return (
    <PageContainer title="首页">
      <ProCard>
        <Space direction="vertical" size={16}>
          <Typography.Title level={4} style={{ margin: 0 }}>
            Gin Admin Template
          </Typography.Title>
          <Typography.Text type="secondary">
            当前前端已对接后端 Gin API，系统管理页面覆盖用户、机构、角色、菜单和资源接口。
          </Typography.Text>
          <Space wrap>
            <Button type="primary" onClick={() => history.push('/system/users')}>
              用户管理
            </Button>
            <Button onClick={() => history.push('/system/orgs')}>机构管理</Button>
            <Button onClick={() => history.push('/system/roles')}>角色管理</Button>
          </Space>
        </Space>
      </ProCard>
      <StatisticCard.Group style={{ marginTop: 16 }}>
        <StatisticCard statistic={{ title: '用户接口', value: 8 }} />
        <StatisticCard statistic={{ title: '机构接口', value: 6 }} />
        <StatisticCard statistic={{ title: '角色接口', value: 7 }} />
        <StatisticCard statistic={{ title: '菜单/资源接口', value: 7 }} />
      </StatisticCard.Group>
    </PageContainer>
  );
};

export default HomePage;
