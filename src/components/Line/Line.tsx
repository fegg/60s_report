import React from "react";
import { Chart, View } from "@antv/g2";

import getContainerId from "./getContainerId";
export interface DataType {
  value: number;
  date: string;
  dept: string;
  type: string;
}

export const colors = [
  "#FF4646",
  "#21A97A",
  "#1874FF",
  "#FFC214",
  "#025DF4",
  "#DB6BCF",
  "#2498D1",
  "#BBBDE6",
  "#4045B2",
  "#21A97A",
  "#FF745A",
  "#007E99",
  "#FFA8A8",
  "#2391FF",
  "#FFC328",
  "#A0DC2C",
  "#946DFF",
  "#626681",
  "#EB4185",
  "#CD8150",
  "#36BCCB",
  "#327039",
  "#803488",
  "#83BC99",
];

const Line = (props: {
  data: DataType[];
  dept: string;
  showGrid?: boolean;
  flipPage?: boolean;
  legendWidth?: number;
  scale?: number;
}) => {
  const ref = React.useRef<View>();
  const {
    dept,
    data,
    showGrid = false,
    flipPage = true,
    legendWidth,
    scale = 1,
  } = props;

  const id = getContainerId();

  React.useEffect(() => {
    if (ref.current) {
      const chart = ref.current;
      chart.changeData(data);
      chart.filter("dept", (value) => value === dept);
      chart.render();
    }
  }, [dept, data]);

  React.useEffect(() => {
    const chart = new Chart({
      container: id,
      autoFit: true,
      renderer: "svg",
      height: 220 * scale,
    });

    chart.data(data);
    chart.filter("dept", (value) => value === dept);

    chart.axis("date", {
      label: {
        style: {
          fontSize: 12 * scale,
        },
        formatter: (val) => {
          return val.substr(5);
        },
      },
    });

    chart.axis("value", {
      label: {
        style: {
          fontSize: 12 * scale,
        },
      },
      grid: showGrid
        ? {
            line: {
              style: {
                stroke: "#1A3A73",
              },
            },
          }
        : null,
    });

    chart.legend({
      flipPage,
      itemWidth: legendWidth,
      label: {
        style: {
          fontSize: 20 * scale,
        },
      },
      marker: (_, index) => {
        return {
          symbol: "circle",
          style: {
            fill: colors[index],
          },
        };
      },
    });

    chart.tooltip({
      showCrosshairs: true,
      title: (_, datum) => {
        return `${datum["value"]}人日 (${datum["percent"]}%)`;
      },
    });

    // 部门项目绿色
    // 公司项目蓝色
    // 日常事项橙色
    // 损耗 红色
    chart.line().position("date*value").color("type", colors).shape("smooth");

    ref.current = chart;
    chart.render();
  }, []);

  return <div id={id} style={{ width: "98%" }}></div>;
};

export default Line;
