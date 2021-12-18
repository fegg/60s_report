import type { DataType } from './Line';
import { IDeptTrend } from "-/types/Departments";

export default function (dayData: IDeptTrend['dayData'], id?: string) {
  return dayData.reduce((acc, cur) => {
    const list = [] as (DataType & { percent: number })[];

    cur.deptList.forEach(dept => {
      const cost = dept.deptInfo.costInfo;
      const percent = (key: keyof typeof cost) => + (100 * (+cost[key] / +cost.total)).toFixed(2);
      list.push(
        {
          value: +cost.matterCost / 100,
          dept: dept.deptCode || id!,
          date: cur.day,
          type: '日常事项',
          percent: percent('matterCost')
        },
        {
          value: +cost.loss / 100,
          dept: dept.deptCode || id!,
          date: cur.day,
          type: '损耗',
          percent: percent('loss')
        }
      );
    });
    return [
      ...acc,
      ...list
    ];
  }, [] as DataType[]);
}