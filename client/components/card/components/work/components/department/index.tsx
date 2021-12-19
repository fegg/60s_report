import React, { useCallback, useEffect, useState } from 'react';
import useDepartmentList from '@/hooks/useDepartmentList';
import { CascadePicker } from 'antd-mobile';

function Department(props: any) {
  const { value, onChange, onCancel, onInit } = props;
  const departments: any = useDepartmentList();

  const department2options = useCallback((items: any) => {
    let departs: any = {};
    const data = items.map((item: any) => {
      departs[item.deptCode] = item;
      const data: any = {
        value: item.deptCode,
        label: item.name,
      };
      if (item.children && item.level < 2) {
        const [options, subDeparts] = department2options(item.children);
        data.children = options;
        departs = { ...departs, ...subDeparts };
      }
      return data;
    });
    return [data, departs];
  }, []);

  useEffect(() => {
    const [options, departs] = department2options(departments);
    setOptions(options);
    setDeparts(departs);
    departs[value] && onInit(departs[value].name);
  }, [department2options, departments, onInit, value]);

  const [options, setOptions] = useState<any>([]);
  const [departs, setDeparts] = useState<any>({});

  return (
    <CascadePicker
      options={options}
      visible={props.visible}
      onConfirm={(values) => {
        const names = values.filter((item) => item).map((v: any) => departs[v].name);
        values = values.filter((item) => item);
        onChange(values, names);
        onCancel(2);
      }}
      onCancel={() => onCancel(1)}
    />
  );
}

export default Department;
