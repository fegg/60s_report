import React from "react";
import { Card, ProgressBar } from "antd-mobile";
import type { CardProps } from "antd-mobile/es/components/card";
import classnames from "classnames";
import { isMobile } from "@/utils";
import "./styles.less";

interface MecCardProps extends CardProps {
  fill?: string;
  inline?: boolean;
  clickable?: boolean;
}

const defaultColor = "rgba(254, 141, 111, 40%)";

const MedCard: React.FC<MecCardProps> = (props) => {
  const { fill = "#fff", inline = false, clickable = false, ...other } = props;

  return (
    <Card
      {...other}
      className={classnames("med-card", props.className, {
        "med-card-clickable": clickable,
        "med-card-mobile": isMobile(),
        "med-card-inline": inline,
      })}
      style={{ backgroundColor: fill, ...props.style }}
    />
  );
};

export default MedCard;

interface CardProgressProps extends MecCardProps {
  progress: string | number;
  // 展示类型
  type?: string;
  color?: string;
  customExtra?: any;
}
export const CardProgress: React.FC<CardProgressProps & CardProps> = function ({
  color,
  customExtra,
  ...props
}) {
  const percent = Number(props.progress);

  return (
    <MedCard
      {...props}
      bodyStyle={{ display: "flex", ...props.bodyStyle }}
      extra={
        customExtra ? (
          customExtra
        ) : (
          <div style={{ display: "flex", alignItems: "center" }}>
            <div style={{ width: 84, marginRight: 8 }}>
              <ProgressBar
                percent={percent}
                style={{
                  "--fill-color": color || "#1677ff",
                }}
              />
            </div>
            <span className="card-title-progress-num">{`${percent}%`}</span>
          </div>
        )
      }
    />
  );
};

interface CardStaticsProps extends MecCardProps {
  quota: string;
  desc?: CardProps["extra"];
}

export const CardStatics: React.FC<CardStaticsProps> = function (props) {
  const { title, className, clickable, fill = defaultColor, ...other } = props;
  return (
    <MedCard
      className={classnames(
        "card-statics-main",
        { "med-card-clickable": clickable },
        className
      )}
      inline={true}
      fill={fill}
      {...other}
    >
      <p className="card-statics-title">{title}</p>
      <h4 className="card-statics-quota">{props.quota}</h4>
      <div className="card-statics-desc">{props.desc}</div>
    </MedCard>
  );
};
