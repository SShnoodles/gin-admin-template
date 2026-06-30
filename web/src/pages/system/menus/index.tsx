import { PlusOutlined } from '@ant-design/icons';
import {
  ModalForm,
  PageContainer,
  ProFormDigit,
  ProFormSelect,
  ProFormText,
  ProFormTreeSelect,
  ProTable,
} from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { App, Button, Popconfirm } from 'antd';
import { useRef, useState } from 'react';
import {
  createMenu,
  deleteMenu,
  getMenuResources,
  getMenus,
  updateMenu,
} from '@/services/admin';
import type { Id, Menu } from '@/services/admin/types';
import { loadParentMenuTreeData, loadResourceOptions } from '../components/options';

type MenuForm = Menu & { parentName?: string; resourceIds?: Id[] };

const MenuPage = () => {
  const actionRef = useRef<ActionType>(undefined);
  const { message } = App.useApp();
  const [open, setOpen] = useState(false);
  const [current, setCurrent] = useState<MenuForm>();
  const [parentLocked, setParentLocked] = useState(false);
  const formKey = `${current?.id || 'create'}-${current?.pid || '0'}-${parentLocked}`;

  const openForm = async (record?: Menu) => {
    if (record?.id) {
      const resourceIds = await getMenuResources(record.id);
      setCurrent({ ...record, resourceIds: resourceIds.map(String) });
      setParentLocked(false);
    } else {
      setCurrent({ parentName: '顶级菜单', pid: '0', sort: 0 });
      setParentLocked(true);
    }
    setOpen(true);
  };

  const columns: ProColumns<Menu>[] = [
    { title: '菜单名称', dataIndex: 'name' },
    { title: '路径', dataIndex: 'path' },
    { title: '图标', dataIndex: 'icon', search: false },
    { title: '排序', dataIndex: 'sort', search: false, width: 90 },
    { title: '创建时间', dataIndex: 'createdAt', valueType: 'dateTime', search: false },
    {
      title: '操作',
      valueType: 'option',
      render: (_, record) => [
        <Button
          key="addChild"
          type="link"
          onClick={() => {
            setCurrent({
              parentName: record.name || record.path || String(record.id),
              pid: record.id,
              sort: 0,
            });
            setParentLocked(true);
            setOpen(true);
          }}
        >
          添加下级
        </Button>,
        <Button key="edit" type="link" onClick={() => openForm(record)}>
          编辑
        </Button>,
        <Popconfirm
          key="delete"
          title="确认删除该菜单？子菜单也会被删除。"
          onConfirm={async () => {
            if (!record.id) return;
            await deleteMenu(record.id);
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
    <PageContainer title="菜单管理">
      <ProTable<Menu>
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={false}
        request={async () => {
          const data = await getMenus();
          return { data, success: true };
        }}
        pagination={false}
        toolBarRender={() => [
          <Button key="new" type="primary" icon={<PlusOutlined />} onClick={() => openForm()}>
            新建菜单
          </Button>,
        ]}
      />
      <ModalForm<MenuForm>
        key={formKey}
        title={current?.id ? '编辑菜单' : '新建菜单'}
        open={open}
        modalProps={{ destroyOnHidden: true, onCancel: () => setOpen(false) }}
        initialValues={current}
        onFinish={async (values) => {
          const payload = parentLocked ? { ...values, pid: current?.pid } : values;
          if (current?.id) {
            await updateMenu(current.id, { ...current, ...payload });
          } else {
            await createMenu(payload);
          }
          message.success('保存成功');
          setOpen(false);
          actionRef.current?.reload();
          return true;
        }}
      >
        <ProFormText name="name" label="菜单名称" rules={[{ required: true }]} />
        {parentLocked ? (
          <ProFormText name="parentName" label="上级菜单" disabled />
        ) : (
          <ProFormTreeSelect
            name="pid"
            label="上级菜单"
            rules={[{ required: true }]}
            fieldProps={{
              allowClear: false,
              showSearch: true,
              treeDefaultExpandAll: true,
              treeNodeFilterProp: 'label',
            }}
            request={() => loadParentMenuTreeData(current?.id)}
          />
        )}
        <ProFormText name="path" label="路径" rules={[{ required: true }]} />
        <ProFormText name="icon" label="图标" rules={[{ required: true }]} />
        <ProFormDigit name="sort" label="排序" min={0} rules={[{ required: true }]} />
        <ProFormSelect
          name="resourceIds"
          label="关联资源"
          fieldProps={{ mode: 'multiple' }}
          request={loadResourceOptions}
        />
      </ModalForm>
    </PageContainer>
  );
};

export default MenuPage;
