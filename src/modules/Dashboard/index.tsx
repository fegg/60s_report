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
  const DeptCard = makeCardData(works.data as any, () => {});
  return (
    <>
      <Card className="dashboard-top-chart">
        <Line data={transform(data.data.dayData as any)} dept={"001-010"} />
      </Card>
      <div className="dashboard-detail-cards">
        <DeptCard
          {...{
            title: "日常事项",
            factor: "matter",
          }}
        />
        <CardProgress title="日常损耗" progress={lossRatio} />
      </div>
    </>
  );
};

export default DashBoard;
