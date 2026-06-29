import { PlusOutlined } from '@ant-design/icons';
import {
  ModalForm,
  PageContainer,
  ProFormText,
  ProFormTreeSelect,
  ProTable,
} from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { App, Button, Popconfirm } from 'antd';
import { useRef, useState } from 'react';
import {
  createOrg,
  deleteOrg,
  getOrgMenus,
  getOrgs,
  updateOrg,
} from '@/services/admin';
import type { Id, Org } from '@/services/admin/types';
import { loadMenuTreeData } from '../components/options';

type OrgForm = Org & { menuIds?: Id[] };

const OrgPage = () => {
  const actionRef = useRef<ActionType>(undefined);
  const { message } = App.useApp();
  const [open, setOpen] = useState(false);
  const [current, setCurrent] = useState<OrgForm>();

  const openForm = async (record?: Org) => {
    if (record?.id) {
      const menuIds = await getOrgMenus(record.id);
      setCurrent({ ...record, menuIds: menuIds.map(String) });
    } else {
      setCurrent(undefined);
    }
    setOpen(true);
  };

  const columns: ProColumns<Org>[] = [
    { title: '机构名称', dataIndex: 'name' },
    { title: '统一社会信用代码', dataIndex: 'creditCode', search: false },
    { title: '地址', dataIndex: 'address', search: false, ellipsis: true },
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
          title="确认删除该机构？"
          onConfirm={async () => {
            if (!record.id) return;
            await deleteOrg(record.id);
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
    <PageContainer title="机构管理">
      <ProTable<Org>
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        request={async (params) => {
          const result = await getOrgs(params);
          const data = Array.isArray(result) ? result : result.data || [];
          return { data, total: Array.isArray(result) ? data.length : result.total, success: true };
        }}
        toolBarRender={() => [
          <Button key="new" type="primary" icon={<PlusOutlined />} onClick={() => openForm()}>
            新建机构
          </Button>,
        ]}
      />
      <ModalForm<OrgForm>
        title={current?.id ? '编辑机构' : '新建机构'}
        open={open}
        modalProps={{ destroyOnHidden: true, onCancel: () => setOpen(false) }}
        initialValues={current}
        onFinish={async (values) => {
          if (current?.id) {
            await updateOrg(current.id, values);
          } else {
            await createOrg(values);
          }
          message.success('保存成功');
          setOpen(false);
          actionRef.current?.reload();
          return true;
        }}
      >
        <ProFormText name="name" label="机构名称" rules={[{ required: true }]} />
        <ProFormText name="creditCode" label="统一社会信用代码" rules={[{ required: true }]} />
        <ProFormText name="address" label="地址" />
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

export default OrgPage;
