import React, { useContext, useEffect } from 'react';
import { TreeSelect } from 'antd-mobile';
import type { TreeSelectOption } from 'antd-mobile/es/components/tree-select';
import { AppContext } from '@/pages/$daily/context';
// import { loadCache, updateCache, OptionType } from './cache';

export interface ITreeSelectProps {
  options: TreeSelectOption[];
  columns?: number;
  defaultValue?: string[];
  cacheKey?: string;
  onChange: (value: TreeSelectOption) => void;
}

const cacheKeyMap: any = {
  work_tree: 'workItems',
  project_tree: 'projectItems',
};

// 从 children 中寻找 key
function findParentKey(collection: TreeSelectOption[], key: string) {
  const parent = collection.find((item) =>
    item.children.map((item: TreeSelectOption) => item.value).includes(key)
  );
  return parent?.value;
}

const ITreeSelect: React.FC<ITreeSelectProps> = (props) => {
  React.useEffect(() => {
    const elms = document.querySelectorAll('.wrap-tree-select .adm-tree-select-column');
    if (elms.length === 2) {
      (elms[0] as HTMLDivElement).style.width = '40%';
      (elms[1] as HTMLDivElement).style.width = '60%';
    }
  }, []);
  const { columns = 2, cacheKey = 'work_tree' } = props;
  const { recentOptions } = useContext(AppContext);

  const options = props.options.map((item) => {
    const node = {
      label: item.label,
      value: item.value,
    };

    if (item.children.length) {
      function formatLabel(label: string) {
        if (typeof label === 'string') {
          return <span dangerouslySetInnerHTML={{ __html: label }}></span>;
        } else {
          return label;
        }
      }
      return {
        ...node,
        children: item.children.map((child: TreeSelectOption) => ({
          label: formatLabel(child.label),
          value: child.value,
        })),
      };
    }

    if (+item.value < 0) {
      return node;
    }

    return {
      ...node,
      // 特殊处理，解决 tree-select 性能问题
      children: [{ ...node, value: `|${node.value}` }],
    };
  });

  // 加入最近选择
  const recent = {
    label: '最近选择',
    value: '-1',
    children: recentOptions[cacheKeyMap[cacheKey]],
  };
  options.unshift(recent);

  // 处理选择历史
  const values = props.defaultValue as [string, string];
  const parentKeys = props.options.map((item) => item.value);
  if (parentKeys.includes(values[1])) {
    values[0] = values[1];
    values[1] = `|${values[1]}`;
  } else if (values[1]) {
    values[0] = findParentKey(props.options, values[1]);
  }

  useEffect(() => {
    document.querySelector('.adm-tree-select-item-active')?.scrollIntoView();
  }, [values]);

  return (
    <TreeSelect
      options={options}
      defaultValue={values}
      className="wrap-tree-select"
      onChange={(value, nodes) => {
        if (value.length === columns) {
          const [, selected] = nodes.options;

          // 必须在截断 | 之前缓存
          // updateCache(cacheKey, selected as OptionType);
          selected.value = selected.value.replace(/\||~|!/, '');
          props.onChange([value[0], selected.value]);
        }
      }}
    ></TreeSelect>
  );
};

export default ITreeSelect;
