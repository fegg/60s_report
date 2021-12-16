export const SERVER_STATUS_NORMAL = 1;
export const SERVER_STATUS_SHELVE = 3;
// 是数据库默认时间
export const DEFAULT_TIME = '0001-01-01';

export const SERVER_STATUS = [
  { name: '正常', value: SERVER_STATUS_NORMAL, color: 'green' },
  // { name: '下架', value: SERVER_STATUS_SHELVE, color: 'red' },
];

// 项目等级
export const LEVELS = [
  {
    value: '4',
    label: 'S',
  },
  {
    value: '3',
    label: 'A',
  },
  {
    value: '2',
    label: 'B',
  },
  {
    value: '1',
    label: 'C',
  },
];
let LEVEL_COLORS = ['magenta', 'orange', 'blue', null];
export const LEVEL_MAP = LEVELS.reduce((total: any, item: any, index: number): any => {
  total[item.value] = {
    label: item.label,
    color: LEVEL_COLORS[index],
  };
  return total;
}, {});

// 项目级别
export const RANKS = [
  {
    value: '1',
    label: '公司级别',
    // children: [],
  },
  {
    value: '2',
    label: '部门级别',
    /* children: [
      {
        value: '1',
        label: '日常运营',
      },
      {
        value: '2',
        label: '日常管理',
      },
      {
        value: '3',
        label: '日常工作',
      },
    ], */
  },
];
export const RANK_MAP = RANKS.reduce((total: any, item: any): any => {
  total[item.value] = item.label;
  return total;
}, {});

// 风险等级
export const RISK_LEVELS = [
  {
    value: '1',
    label: '低',
  },
  {
    value: '2',
    label: '中',
  },
  {
    value: '3',
    label: '高',
  },
];

export const AUTH_CONTROL_MAP: any = {
  '1': '公开',
  '2': '保密',
};

export const tag_colors = [
  'rgba(254, 141, 111, 40%)',
  'rgba(255, 157, 30, 15%)',
  'rgba(155, 219, 197, 40%)',
  'rgba(30, 116, 255, 15%)',
  'rgba(246, 215, 195, 48%)',
  'rgba(221, 221, 221, 29%)',
  'rgb(157, 223, 211)',
];

// alias for tag_colors;
export const colors = tag_colors;

export const factors = [
  ['companyProject', '公司项目'],
  ['deptProject', '部门项目'],
  ['matter', '日常事项'],
  ['manager', '管理事项'],
];

export const PROJECT_DATA_TYPE = {
  project: '1',
  tag: '2',
};

// 项目看板类型
export const PROJECT_TYPES = [
  {
    value: PROJECT_DATA_TYPE.project,
    label: '项目',
  },
  {
    value: PROJECT_DATA_TYPE.tag,
    label: '项目标签',
  },
];

export const rankMap = RANKS.reduce((total: any, item: any): any => {
  total[item.value] = item.label;
  return total;
}, {});
export const levelMap = LEVELS.reduce((total: any, item: any): any => {
  total[item.value] = item.label;
  return total;
}, {});

// 立项流程
export const BID_STATUS_INIT = '0'; // 默认
export const BID_STATUS_UNSTART = '1'; // 报名中
export const BID_STATUS_SCORE = '2'; // 评分中
export const BID_STATUS_FINISHED = '3'; // 已完成

export const projectStatusList = [
  {
    label: '竞标报名中',
    color: 'purple',
    value: BID_STATUS_UNSTART,
  },
  {
    label: '竞标评分中',
    color: 'orange',
    value: BID_STATUS_SCORE,
  },
  {
    label: '竞标完成',
    color: 'success',
    value: BID_STATUS_FINISHED,
  },
];
export const projectStatusMap = projectStatusList.reduce((total: any, item: any): any => {
  total[item.value] = item;
  return total;
}, {});

export const CHECK_STATUS = [
  {
    label: '未设立验收标准',
    color: 'red',
    value: '-3',
  },
  {
    label: '验收标准确立中',
    color: 'orange',
    value: '-2',
  },
  {
    label: '验收标准已确立',
    color: 'success',
    value: '-1',
  },
  {
    label: '部分验收已审核',
    color: 'green',
    value: '1',
  },
  {
    label: '项目考核结束',
    color: '#87d068',
    value: '2',
  },
];

export const CHECK_STATUS_MAP = CHECK_STATUS.reduce((total: any, item: any): any => {
  total[item.value] = item;
  return total;
}, {});

export const statusName: any = {
  '10': '待审核',
  '11': '待审核',
  '20': '已通过',
  '30': '已拒绝',
};
export const IS_SUBJECT_TYPE = 1; // 项目信息
export const IS_BID_TYPE = 2; // 开题信息
