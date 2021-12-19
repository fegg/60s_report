import React, { forwardRef, useEffect, useImperativeHandle, useState } from 'react';
import { Calendar } from 'react-date-range';
// @ts-ignore
import * as locales from 'react-date-range/dist/locale';
import usePop from '@/hooks/usePop';
import moment from 'moment';
import useGetUserFillLog from '@/hooks/useGetUserFillLog';
import './index.scss';

interface Props {
  timeScope?: string;
  onChange?: any;
}
function PopCalendar(props: Props, ref: any) {
  const { timeScope = '', onChange } = props;
  const [Wrapper, onWrapperChange, setShowText] = usePop(timeScope);
  const [currentDate, setCurrentDate] = useState(moment(timeScope).toDate());

  const { userLogs, getUserFillLogCallback }: any = useGetUserFillLog({ timeScope });
  useEffect(() => {
    setShowText(timeScope);
  }, [setShowText, timeScope]);

  useImperativeHandle(
    ref,
    () => ({
      getUserFillLogCallback,
    }),
    [getUserFillLogCallback]
  );

  function customDayContent(day: Date) {
    let extraDot = null;
    if (
      userLogs.indexOf(moment(day).format('YYYY-MM-DD')) === -1 &&
      moment(day) <= moment() &&
      [6, 7].indexOf(moment(day).isoWeekday()) === -1
    ) {
      extraDot = (
        <div
          style={{
            height: '5px',
            width: '5px',
            borderRadius: '100%',
            background: 'orange',
            position: 'absolute',
            top: 2,
            right: 2,
          }}
        />
      );
    }
    return (
      <div>
        {extraDot}
        <span>{moment(day).date()}</span>
      </div>
    );
  }

  return (
    <Wrapper style={{ width: '100%' }}>
      <Calendar
        showDateDisplay={false}
        showMonthAndYearPickers={false}
        locale={locales['zhCN']}
        maxDate={moment().toDate()}
        date={currentDate}
        dayContentRenderer={customDayContent}
        onChange={(value) => {
          setCurrentDate(moment(value).toDate());
          onChange(moment(value).format('YYYY-MM-DD'));
          onWrapperChange(); // close
        }}
        onShownDateChange={(day: Date) => {
          // 解决重复刷新组件的问题
          setCurrentDate(moment(day).toDate());
          getUserFillLogCallback(day);
        }}
      />
    </Wrapper>
  );
}

export default forwardRef(PopCalendar);
