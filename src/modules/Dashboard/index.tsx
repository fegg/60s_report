import React from "react";
import Line, { DataType } from "@/components/Line";
import transform from "@/components/Line/transform";
import data from "./data.json";
import { startOfWeek, addDays, format } from "date-fns"
import works from "./works.json";
import Card, { CardProgress } from "@/components/Card";

const lossRatio = (
  (100 * +works.data.costInfo.loss) /
  +works.data.costInfo.total
).toFixed(1);

const DashBoard = () => {
  const list = works.data.matterInfo.list;
  const total:number = list.reduce((pre, cur) => pre + (+cur.cost), 0);
  const today = new Date();
  const start = startOfWeek(today, { weekStartsOn: 1});
  const dept ="001-00A";

  const rand = () => Array(list.length).fill(1).map(()=> Math.random() * 100);

  const chart_data = [
    start,
    addDays(start, 1),
    addDays(start, 2),
    addDays(start, 3),
    addDays(start, 4),
  ].reduce((pre, d) => {
    const data = rand();
    const day_total = data.reduce((a, b) => a + b, 0)
    pre.push(...list.map((item, index)  => {
      const value = data[index];
      return {
        value: +value.toFixed(1),
        dept, 
        date: format(d, 'yyyy-MM-dd'),
        type: item.title,
        percent: (100 * value /day_total).toFixed(2)
      }
    }))
    console.log('test: ', pre);
    return pre;
  }, [] as DataType[]);

  return (
    <>
      <Card className="dashboard-top-chart">
        <Line data={chart_data} dept={dept} />
      </Card>
      <div className="dashboard-detail-cards">
        {
          works.data.matterInfo.list.map((item, index) => {
            return <CardProgress title={item.title} progress={100*+item.cost/total} key={index} />
          })
        }
        <CardProgress title="日常损耗" progress={lossRatio} />
      </div>
    </>
  );
};

export default DashBoard;
