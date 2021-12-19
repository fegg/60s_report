import React from "react";
import Line from "@/components/Line";
import transform from "@/components/Line/transform";
import data from "./data.json";
import { makeCardData } from "./helper";
import works from "./works.json";
import Card, { CardProgress } from "@/components/Card";

const lossRatio = (
  (100 * +works.data.costInfo.loss) /
  +works.data.costInfo.total
).toFixed(1);

const DashBoard = () => {
  const list = works.data.matterInfo.list;
  const total:number = list.reduce((pre, cur) => pre + (+cur.cost), 0);

  return (
    <>
      <Card className="dashboard-top-chart">
        <Line data={transform(data.data.dayData as any)} dept={"001-00A"} />
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
