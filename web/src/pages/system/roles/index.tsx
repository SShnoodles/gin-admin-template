import { PlusOutlined } from '@ant-design/icons';
import {
  ModalForm,
  PageContainer,
  ProFormSelect,
  ProFormText,
  ProFormTreeSelect,
  ProTable,
} from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { App, Button, Popconfirm } from 'antd';
import { useRef, useState } from 'react';
import {
  createRole,
  deleteRole,
  getRoleMenus,
  getRoles,
  updateRole,
} from '@/services/admin';
import type { Id, Role } from '@/services/admin/types';
import { loadMenuTreeData, loadOrgOptions } from '../components/options';

type RoleForm = Role & { menuIds?: Id[] };

const RolePage = () => {
  const actionRef = useRef<ActionType>(undefined);
  const { message } = App.useApp();
  const [open, setOpen] = useState(false);
  const [current, setCurrent] = useState<RoleForm>();

  const openForm = async (record?: Role) => {
    if (record?.id) {
      const menuIds = await getRoleMenus(record.id);
      setCurrent({ ...record, menuIds: menuIds.map(String) });
    } else {
      setCurrent(undefined);
    }
    setOpen(true);
  };

  const columns: ProColumns<Role>[] = [
    { title: '角色名称', dataIndex: 'name' },
    { title: '角色编码', dataIndex: 'code' },
    { title: '机构', dataIndex: 'orgName', search: false },
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
          title="确认删除该角色？"
          onConfirm={async () => {
            if (!record.id) return;
            await deleteRole(record.id);
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
    <PageContainer title="角色管理">
      <ProTable<Role>
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        request={async (params) => {
          const result = await getRoles(params);
          return { data: result.data || [], total: result.total, success: true };
        }}
        toolBarRender={() => [
          <Button key="new" type="primary" icon={<PlusOutlined />} onClick={() => openForm()}>
            新建角色
          </Button>,
        ]}
      />
      <ModalForm<RoleForm>
        title={current?.id ? '编辑角色' : '新建角色'}
        open={open}
        modalProps={{ destroyOnHidden: true, onCancel: () => setOpen(false) }}
        initialValues={current}
        onFinish={async (values) => {
          if (current?.id) {
            await updateRole(current.id, values);
          } else {
            await createRole(values);
          }
          message.success('保存成功');
          setOpen(false);
          actionRef.current?.reload();
          return true;
        }}
      >
        <ProFormText name="name" label="角色名称" rules={[{ required: true }]} />
        <ProFormText name="code" label="角色编码" rules={[{ required: true }]} />
        <ProFormSelect name="orgId" label="机构" rules={[{ required: true }]} request={loadOrgOptions} />
        <ProFormTreeSelect
          name="menuIds"
          label="授权菜单"
          fieldProps={{
            multiple: true,
            treeCheckable: true,
            showCheckedStrategy: 'SHOW_PARENT',
          }}
          request={loadMenuTreeData}
        />
      </ModalForm>
    </PageContainer>
  );
};

export default RolePage;
