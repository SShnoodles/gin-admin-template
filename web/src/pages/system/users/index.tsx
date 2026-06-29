import { PlusOutlined } from '@ant-design/icons';
import {
  ModalForm,
  PageContainer,
  ProFormSelect,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { App, Button, Popconfirm, Switch } from 'antd';
import { useRef, useState } from 'react';
import {
  createUser,
  deleteUser,
  getUserRoles,
  getUsers,
  toggleUserEnabled,
  updateUser,
} from '@/services/admin';
import type { Id, User } from '@/services/admin/types';
import { loadOrgOptions, loadRoleOptions } from '../components/options';

type UserForm = User & { roleIds?: Id[] };

const UserPage = () => {
  const actionRef = useRef<ActionType>(undefined);
  const { message } = App.useApp();
  const [open, setOpen] = useState(false);
  const [current, setCurrent] = useState<UserForm>();
  const [orgId, setOrgId] = useState<Id>();

  const openForm = async (record?: User) => {
    if (record?.id) {
      const roleIds = await getUserRoles(record.id);
      setCurrent({ ...record, roleIds: roleIds.map(String) });
      setOrgId(record.orgId);
    } else {
      setCurrent(undefined);
      setOrgId(undefined);
    }
    setOpen(true);
  };

  const columns: ProColumns<User>[] = [
    { title: '用户名', dataIndex: 'username' },
    { title: '姓名', dataIndex: 'realName', search: false },
    { title: '工号', dataIndex: 'workNo', search: false },
    { title: '机构', dataIndex: 'orgName', search: false },
    {
      title: '状态',
      dataIndex: 'enabled',
      search: false,
      render: (_, record) => (
        <Switch
          checked={record.enabled}
          checkedChildren="启用"
          unCheckedChildren="禁用"
          onChange={async () => {
            if (!record.id) return;
            await toggleUserEnabled(record.id);
            message.success('状态已更新');
            actionRef.current?.reload();
          }}
        />
      ),
    },
    { title: '创建时间', dataIndex: 'createdAt', valueType: 'dateTime', search: false },
    {
      title: '操作',
      valueType: 'option',
      render: (_, record) => [
        <Button key="edit" type="link" onClick={() => openForm(record)}>
          编辑
        </Button>,
        <Popconfirm
          key="delete"
          title="确认删除该用户？"
          onConfirm={async () => {
            if (!record.id) return;
            await deleteUser(record.id);
            message.success('删除成功');
            actionRef.current?.reload();
          }}
        >
          <Button type="link" danger>
            删除
          </Button>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer title="用户管理">
      <ProTable<User>
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        request={async (params) => {
          const result = await getUsers(params);
          return { data: result.data || [], total: result.total, success: true };
        }}
        toolBarRender={() => [
          <Button key="new" type="primary" icon={<PlusOutlined />} onClick={() => openForm()}>
            新建用户
          </Button>,
        ]}
      />
      <ModalForm<UserForm>
        title={current?.id ? '编辑用户' : '新建用户'}
        open={open}
        modalProps={{ destroyOnHidden: true, onCancel: () => setOpen(false) }}
        initialValues={current}
        onFinish={async (values) => {
          if (current?.id) {
            await updateUser(current.id, values);
          } else {
            await createUser(values);
          }
          message.success('保存成功');
          setOpen(false);
          actionRef.current?.reload();
          return true;
        }}
      >
        <ProFormText name="username" label="用户名" rules={[{ required: true }]} />
        <ProFormText name="realName" label="姓名" rules={[{ required: true }]} />
        <ProFormText name="workNo" label="工号" rules={[{ required: true }]} />
        <ProFormSelect
          name="orgId"
          label="机构"
          rules={[{ required: true }]}
          request={loadOrgOptions}
          fieldProps={{ onChange: (value) => setOrgId(value as Id) }}
        />
        <ProFormSelect
          name="roleIds"
          label="角色"
          fieldProps={{ mode: 'multiple' }}
          params={{ orgId }}
          request={async () => loadRoleOptions(orgId)}
        />
      </ModalForm>
    </PageContainer>
  );
};

export default UserPage;
