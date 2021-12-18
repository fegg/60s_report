import { Radio, Space } from 'antd-mobile';
import React, { useEffect, useState } from 'react';
import Department from '../department';
import styles from './index.module.scss';

const radioStyle = {
  '--icon-size': '16px',
  '--font-size': '12px',
  '--gap': '4px',
};

/**
 * 工作事项归属组件模块
 * @param props
 * @returns
 */
function RedioDepartment(props: any) {
  const { value, onChange, deptCode } = props;
  const [radioValue, setRadioValue] = useState(1);
  const [departVisible, setDepartVisible] = useState(false);
  const [currentDepart, setCurrentDepart] = useState<string>('');

  function onDepartmentChange(values: any, names: any) {
    setCurrentDepart(names.slice(-1)[0]);
    onChange(values.slice(-1)[0]);
  }

  function onInit(name: string) {
    value !== deptCode && setCurrentDepart(name);
  }

  useEffect(() => {
    value && value !== deptCode ? setRadioValue(2) : setRadioValue(1);
  }, [deptCode, value]);

  return (
    <div className={styles.radio}>
      <Radio.Group
        value={radioValue}
        onChange={(v) => {
          if (v === 2) {
            setDepartVisible(true);
          } else {
            setRadioValue(+v);
            setCurrentDepart('');
            onChange(deptCode);
          }
        }}
      >
        <Space>
          <Radio style={radioStyle} value={1}>
            归属本部门
          </Radio>
          <Radio style={radioStyle} value={2}>
            归属
            {currentDepart ? (
              <span
                onClick={() => {
                  setDepartVisible(true);
                }}
              >
                {currentDepart}
              </span>
            ) : (
              '其它部门'
            )}
          </Radio>
        </Space>
      </Radio.Group>

      <Department
        visible={departVisible}
        value={value}
        onInit={onInit}
        onChange={(values: any, names: any) => {
          onDepartmentChange(values, names);
        }}
        onCancel={(isCurrent: number) => {
          setDepartVisible(false);
          setRadioValue(isCurrent);
        }}
      />
    </div>
  );
}

export default RedioDepartment;
