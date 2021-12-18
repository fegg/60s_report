import React, { useEffect, useState } from 'react';
import { Card } from 'antd-mobile';
import useUserProjectStatics from '@/hooks/useUserProjectStatics';
import DeptTime from '@/modules/Department/components/Filter/DeptTime';
import Empty from '@/modules/Daily/components/empty';
import { mobileInit } from '@/utils';
import { setTitle } from '@/utils/dd';
import { getDefRangeTime } from '@/utils';
import moment from 'moment';
import Loading from '../components/loading';
import styles from './index.module.scss';

type RangeTime = { startTime: string; endTime: string };

function Statics() {
  const defaultDate = getDefRangeTime();
  const [rangeTime, setRangeTime] = useState<RangeTime>({
    startTime: moment(defaultDate.startTime).format('YYYY-MM-DD'),
    endTime: moment(defaultDate.endTime).format('YYYY-MM-DD'),
  });
  useEffect(() => {
    setTitle('医联一分钟日报');
    mobileInit();
  }, []);
  const userProjectStatics = useUserProjectStatics(rangeTime);
  return (
    <div className={styles.page}>
      <div className={styles.wrapper}>
        <div className={styles.header}>
          日报记录查询：
          <div className={styles.timer}>
            <DeptTime
              showYear={true}
              onChange={(value: any) => {
                setRangeTime(value);
              }}
            />
          </div>
        </div>
        <div className={styles['header-placeholder']}></div>
        <Loading visible={!!userProjectStatics} loading={true}>
          <div className={styles.projects}>
            {(userProjectStatics || []).map((item: any, index: number) => {
              return (
                <Card key={`${item.project.projectId}${index}`} className={styles.card}>
                  <div className={styles.cardContent}>
                    <span className={styles.title}>{item.project.projectName}</span>
                    <span className={styles.time}>{item.cost / 100}人/日</span>
                  </div>
                </Card>
              );
            })}
          </div>
          {!(userProjectStatics || []).length && <Empty title="未查询到时间范围内的填报记录" />}
        </Loading>
      </div>
    </div>
  );
}

export default Statics;
