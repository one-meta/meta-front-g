import {%sColumns } from '@/columns/%s';
import { delete%sId, get%s, post%s, post%sBulk__openAPI__delete, put%sId } from '@/services/meta/%s';
import { ActionType, ProColumns, ProForm } from '@ant-design/pro-components';
import { useRef, useState } from 'react';
import ProDescriptionsLayout from '@/layouts/extend/ProDescriptionsLayout';
import ProTableLayout from '@/layouts/extend/ProTableLayout';
import PopconfirmDeleteLayout from '@/layouts/options/PopconfirmDeleteLayout';
import ShowRowDetailLayout from '@/layouts/options/ShowRowDetailLayout';
import EditRowLayout from '@/layouts/options/EditRowLayout';
import LinkLayout from '@/layouts/options/LinkLayout';
import { transformTime } from '@/utils/transform';
import RemarkProFormTextArea from '@/layouts/form/RemarkProFormTextArea';


const IndexPage: React.FC = () => {
  const actionRef = useRef<ActionType>();

  //多行选择
  const [selectedRowsState, setSelectedRows] = useState < API.%s[] > ([]);

  // 新建 使用ModalForm
  const [addModalDetail, setAddModalDetail] = useState<boolean>(false);
  const [addDetailPage, setAddDetailPage] = useState<any>();

  // 查看 使用Drawer
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [showDetailPage, setShowDetailPage] = useState<any>();

  // 更新 使用ModalForm
  const [updateModalDetail, setUpdateModalDetail] = useState<boolean>(false);
  // 更新 使用DrawerForm
  const [updateDrawerDetail, setUpdateDrawerDetail] = useState<boolean>(false);
  // 更新的值
  const [updateDetailPage, setUpdateDetailPage] = useState<any>();
  const [updateValues, setUpdateValues] = useState<any>();

  //新建页面的字段
  //如果编辑的字段也和新建页面一样（可操作的字段相同），可以共用
  const newPageDetail = (
    <ProForm.Group>
      <RemarkProFormTextArea />
      {/*TODO：根据需求增加表单*/}
      {/* <ProFormText name="username" label="用户名" rules={[{required: true}]} */}
      {/* <ProFormText
        rules={[
          {
            required: true,
          },
          {
            min: 8,
            message: '密码不能少于8个字符',
          },
        ]}
        name="password"
        label="密码"
        tooltip="密码不能少于8个字符"
      /> */}
    </ProForm.Group>
  )


  //增加行操作：编辑、查看、删除
  const columns: ProColumns<API.%s > [] =[
    ...%sColumns,
    {
      title: '操作',
      valueType: 'option',
      hideInDescriptions: true,
      render: (_text, record) => [
        //编辑
        <EditRowLayout
          key="edit"
          setUpdateValues={setUpdateValues}
          setUpdateDetailPage={setUpdateDetailPage}
          //在这里控制使用drawer还是modal，二选一
          // setUpdateDrawerDetail={setUpdateDrawerDetail}
          setUpdateModalDetail={setUpdateModalDetail}

          //需要更新的值，字段名与columns中相同
          //updateValues中id必须设置
          updateValues={
            {//id必须设置
              id: record.id,
              //将数组转换成单行数据，如果有
              // dataList: record.dataList?.join('\n'),
              //如果没有备注，可以去掉
              remark: record.remark,
              //其他需要更新的字段可以在这里添加
            }
          }
          //更新页面表单
          updateDetailPage={<>
            {/*TODO：其他需要更新的字段的表单可以在这里添加*/}

            {/*如果没有备注，可以去掉*/}
            <RemarkProFormTextArea />
          </>}
        />,
        //详情（跳转页面）
        // TODO：根据需求修改路径
        <LinkLayout key="link" to={`/%s/detail/${record.id}`} />,
        //查看
        <ShowRowDetailLayout
          key="view"
          setShowDetailPage={setShowDetailPage}
          setShowDetail={setShowDetail}
          detailPage={<>
            <ProDescriptionsLayout data={record} columns={columns} column={2} />
          </>}
        />,
        //删除
        <PopconfirmDeleteLayout
          key="delete"
          actionRef={actionRef}
          api={delete%sId}
          id={record.id}
        />,
      ],
    }
  ];



  //POST 数据转换，可以在提交时对字段进行特殊处理
  // const transformPostData = (data: any) => {
  //   return {
  //   // 提交时，将单行数据转换成数组
  //     dataList: data['dataList']?.split(/[(\r\n)\r\n]+/)
  //   }
  // }

  //GET 数据转换，可以在查询时对字段进行特殊处理
  // const transformGetData = () => {
  //   return {
  //   // 查询时，增加字段type，值为g
  //     type: 'g'
  //   }
  // }

  return <ProTableLayout
    actionRef={actionRef}
    columns={columns}

    //Service 方法
    getMethod={get%s}
    editMethod={put%sId}
    newMethod={post%s}
    deleteBulkMethod={post%sBulk__openAPI__delete}


    //默认false，如果true，需要同时传入转换函数
    transform={true}
    transformTime={transformTime}
    //提交时（新建/编辑），数据转换：将单行数据转换成数组
    // transformPostData={transformPostData}
    // 查询时，数据转换
    // transformGetData={transformGetData}


    //多选操作
    //显示多选操作；默认显示
    showRowSelect={true}
    selectedRowsState={selectedRowsState}
    setSelectedRows={setSelectedRows}

    //单击行数据，增加背景色，双击取消；默认不生效
    //每次都会遍历所有行，有点消耗资源？
    // rowSelectBackground={true}

    //显示分页栏；默认显示
    pagination={true}
    //显示搜索框；默认显示
    showSearch={true}
    //显示工具栏；默认显示
    showToolbar={true}

    //新建
    //显示新建按钮；默认显示
    showNew={true}
    addModalDetail={addModalDetail}
    setAddModalDetail={setAddModalDetail}
    addDetailPage={addDetailPage}
    setAddDetailPage={setAddDetailPage}
    newPageDetail={newPageDetail}

    // 显示
    showDetail={showDetail}
    setShowDetail={setShowDetail}
    showDetailPage={showDetailPage}
    setShowDetailPage={setShowDetailPage}

    // 更新
    updateDetailPage={updateDetailPage}
    setUpdateDetailPage={setUpdateDetailPage}
    updateValues={updateValues}
    setUpdateValues={setUpdateValues}
    // 编辑提交时，二次确认，默认否
    // secondConfirm={true}
    // 二次确认文本
    // secondConfirmContent={"是否继续"}

    // 使用ModalForm或者使用DrawerForm更新
    updateDrawerDetail={updateDrawerDetail}
    setUpdateDrawerDetail={setUpdateDrawerDetail}
    updateModalDetail={updateModalDetail}
    setUpdateModalDetail={setUpdateModalDetail}
  />

};
export default IndexPage;
