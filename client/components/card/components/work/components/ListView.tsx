import React from 'react';

const ListView = () => {
  return (
    <div>
      
    </div>
  );
};

export default ListView;

// <CheckList
//   className="project-checkList"
//   style={{ padding: '0 0 20px' }}
//   value={[data.projectId]}
//   onChange={(value: any) => {
//     if (!value.length) {
//       return;
//     }
//     const item = worksMap[value[0]];
//     // debugger;
//     let newItem: any = {
//       projectName: item.projectName,
//       projectId: item.projectId,
//       isRequiredComment: +item.isRequiredComment,
//       isForceNotice: +item.isForceNotice,
//       notice: item.notice,
//     };
//     // 消耗部门存在才传递
//     if (item.costDeptCode) {
//       newItem['costDeptCode'] = item.costDeptCode;
//     }

//     // 重选以后重置 remark
//     if (newItem.projectId !== data.projectId) {
//       newItem['remark'] = '';
//     }
//     onChange(newItem);
//     setPanelVisible(false);
//   }}
// >
//   <>
//     {works.map((item: any) => {
//       if (item.children.length) {
//         return (
//           <Collapse
//             key={item.value}
//             style={{ borderBottom: '1px solid #eeeeee', color: '#333333' }}
//           >
//             <Collapse.Panel key={`${item.value}-1`} title={item.label}>
//               {item.children.map((subItem: any) => (
//                 <CheckList.Item key={subItem.value} value={subItem.value}>
//                   {subItem.label}
//                 </CheckList.Item>
//               ))}
//             </Collapse.Panel>
//           </Collapse>
//         );
//       } else {
//         return (
//           <CheckList.Item key={item.value} value={item.value}>
//             {item.label}
//           </CheckList.Item>
//         );
//       }
//     })}
//     {!works.length && <Empty />}
//   </>
// </CheckList>

