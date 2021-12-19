import React from 'react';
import { LEVELS, RANKS } from '@/constant/application';
import styles from './index.module.scss';

let LEVELS_MAP: any = {},
  RANKS_MAP: any = {};
LEVELS.forEach((v) => {
  LEVELS_MAP[v.value] = v.label;
});
RANKS.forEach((v) => {
  RANKS_MAP[v.value] = v.label;
});

function Block(props: {
  title: string;
  items: any[];
  filters: { rank: string; level: string };
  filterKey: 'rank' | 'level';
  ITEM_MAP: any;
  onFilterChange: any;
}) {
  const { title, items, filters, filterKey, ITEM_MAP, onFilterChange } = props;
  return (
    <div>
      <div className={styles.block}>
        <div className={styles.title}>{title}</div>
        <div
          className={styles.tag}
          onClick={(e: any) => {
            let { value } = e.target.dataset;
            if (value === filters[filterKey]) {
              value = '';
            }
            onFilterChange({ [filterKey]: value });
          }}
        >
          {items.map((v) => (
            <div key={v.value}>
              <div
                className={v.value === filters[filterKey] ? styles.checked : ''}
                data-value={v.value}
              >
                {ITEM_MAP[v.value]}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function Filter(props: any) {
  const { filterOptions, onChange: onFilterChange } = props;
  return (
    <div className={styles.filter}>
      <Block
        title="项目级别搜索"
        items={RANKS}
        filters={filterOptions}
        filterKey="rank"
        ITEM_MAP={RANKS_MAP}
        onFilterChange={onFilterChange}
      />
      <Block
        title="项目评级搜索"
        items={LEVELS}
        filters={filterOptions}
        filterKey="level"
        ITEM_MAP={LEVELS_MAP}
        onFilterChange={onFilterChange}
      />

      <div className={styles.block}>
        <div className={styles.title}>搜索结果</div>
        {props.children}
      </div>
    </div>
  );
}

export default Filter;
