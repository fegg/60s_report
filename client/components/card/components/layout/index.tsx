import { TextArea, Toast } from 'antd-mobile';
import { CloseCircleOutline } from 'antd-mobile-icons';
import React, { useContext } from 'react';
import { AppContext } from '@/pages/$daily/context';
import Percent from '../../../percent';
import Project from '../project';
import Work from '../work';
import styles from './index.module.scss';

function CardLayout(props: any) {
  const { data = {}, onChange, onRemove, type } = props;
  const { total } = useContext(AppContext);

  function onChangeValue(val: any, isToast = false) {
    if (total + (val - +data.cost) > 100) {
      val = 100 - total + +data.cost;
      isToast &&
        Toast.show({
          content: '工作精力已达到100%',
        });
    }
    return onChange({ cost: val });
  }
  return (
    <>
      <div className={styles.card}>
        <div className={styles.header}>
          {type === 'project' && <Project {...props} />}
          {type === 'item' && <Work {...props} />}
          <div className={styles.remove} onClick={onRemove}>
            <CloseCircleOutline />
          </div>
        </div>
        <Percent
          onChange={(val: any) => onChangeValue(val)}
          onAfterChange={(val: any) => onChangeValue(val, true)}
          value={data.cost}
        />
        <div className={styles.remark}>
          <div className={styles.title}>{data.isRequiredComment === 1 && <span>*</span>}备注：</div>
          <TextArea
            value={data.remark}
            onChange={(val: any) => onChange({ remark: val })}
            placeholder={`${data.isRequiredComment === 1 ? '请输入备注' : '（选填）'}`}
            rows={1}
            autoSize={{ minRows: 1, maxRows: 5 }}
          />
        </div>
        {data.isRequiredComment === 1 && !data.remark && (
          <div className={styles.tip}>
            <span>*</span>
            {data.notice}
          </div>
        )}
      </div>
    </>
  );
}

export default CardLayout;
