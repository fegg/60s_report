import { Popup } from 'antd-mobile';
import React, { useContext, useEffect, useState } from 'react';
import { RightOutline } from 'antd-mobile-icons';
import { CardContext } from '@/pages/$daily/context';
import styles from '../index.module.scss';
import TreeSelect from './components/TreeSelect';

/**
 *
 * 工作事项卡片头部组件模块
 * @param props
 * @returns
 */
function Work(props: any) {
  const { data = {}, onChange } = props;
  const [panelVisible, setPanelVisible] = useState(false);
  const { works, worksMap } = useContext(CardContext);

  let item = data.projectName || <span className={styles.unselect}>请选择</span>;

  useEffect(() => {
    !data.projectId && setPanelVisible(true);
  }, [data]);

  return (
    <>
      <span
        className={styles.ellipsis}
        onClick={() => {
          setPanelVisible(true);
        }}
      >
        {item}
        <RightOutline />
      </span>

      <Popup
        visible={panelVisible}
        onMaskClick={() => {
          setPanelVisible(false);
        }}
        className={styles.projectPopup}
        bodyStyle={{ borderRadius: '20px 20px 0 0' }}
      >
        <div className={styles.projectHeader}>
          工作事项
          <img
            onClick={() => setPanelVisible(false)}
            className={styles.projectPopupClose}
            src="/_/prod/developer-panel/1462974496569823003.png"
          />
        </div>

        <div className={styles.checkList} style={{ height: '80vh' }}>
          <TreeSelect
            options={works}
            defaultValue={['-1', data.projectId]}
            onChange={(value) => {
              const item = worksMap[value[1]];
              // debugger;
              let newItem: any = {
                projectName: item.projectName,
                projectId: item.projectId,
                isRequiredComment: +item.isRequiredComment,
                isForceNotice: +item.isForceNotice,
                notice: item.notice,
              };
              // 消耗部门存在才传递
              if (item.costDeptCode) {
                newItem['costDeptCode'] = item.costDeptCode;
              }

              // 重选以后重置 remark
              if (newItem.projectId !== data.projectId) {
                newItem['remark'] = '';
              }
              onChange(newItem);
              setPanelVisible(false);
            }}
          />
        </div>
      </Popup>
    </>
  );
}

export default Work;
