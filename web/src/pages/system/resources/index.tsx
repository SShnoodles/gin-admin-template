import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { getResources } from '@/services/admin';
import type { Resource } from '@/services/admin/types';

const methodValueEnum = {
  GET: { text: 'GET', status: 'Success' },
  POST: { text: 'POST', status: 'Processing' },
  PUT: { text: 'PUT', status: 'Warning' },
  DELETE: { text: 'DELETE', status: 'Error' },
};

const ResourcePage = () => {
  const columns: ProColumns<Resource>[] = [
    { title: '资源名称', dataIndex: 'name' },
    { title: '请求方法', dataIndex: 'method', valueEnum: methodValueEnum },
    { title: '路径', dataIndex: 'path', ellipsis: true },
    { title: '创建时间', dataIndex: 'createdAt', valueType: 'dateTime', search: false },
  ];

  return (
    <PageContainer title="资源管理">
      <ProTable<Resource>
        rowKey="id"
        columns={columns}
        request={async (params) => {
          const result = await getResources(params);
          if (Array.isArray(result)) {
            return { data: result, total: result.length, success: true };
          }
          return { data: result.data || [], total: result.total, success: true };
        }}
      />
    </PageContainer>
  );
};

export default ResourcePage;
