import React, { useState } from "react";
import styles from "../styles/report.module.less";
import { format } from "date-fns";
import { message } from "antd";
import { Button, SafeArea, NoticeBar, Checkbox, Toast } from "antd-mobile";
import ProjectOrWork from "./components/index";
import Loading from "./components/loading";
import PopCalendar from "./components/calendar";
import Collapse from "@/components/collapse";

import "./index.scss";

const Report = () => {
  const now = format(new Date(), "yyyy-MM-dd");
  const ref = React.useRef();
  const [total, setTotal] = React.useState<number>(0);
  const [isSubmit, setIsSubmit] = useState<any>();
  const [timeScope, setTimeScope] = React.useState(now);
  const [isAutoFill, setIsAutoFill] = useState<boolean>(false);
  const [items, setItems] = useState<any>([]);

  const depts = [
    {
      detpId: "10",
      deptCode: "001-00A",
      name: "我的部门",
    },
  ];

  function onSubmit() {
    Toast.show({
      content: "提交成功",
      getContainer: () => document.getElementById("report-container")!
    });
  }

  const getItem = React.useCallback(
    (dept: any) => {
      const index = items.findIndex(
        (item: any = {}) => item.userDeptCode === dept.deptCode
      );
      if (index !== -1) {
        return [items[index], index];
      } else {
        return [
          {
            items: [],
            projects: [],
            userDeptCode: dept.deptCode,
          },
          items.length,
        ];
      }
    },
    [items]
  );

  return (
    <div className={styles.page}>
      <div className={styles.wrapper}>
        <div className={styles.header}>
          <div className={styles.timer}>
            <PopCalendar
              ref={ref}
              timeScope={timeScope}
              onChange={(value: string) => {
                setTimeScope(value);
              }}
            />
            <img
              className={styles.arrowRight}
              src="https://pub-med-post.medlinker.com/_/prod/developer-panel/1463517605938594832.png"
            />
          </div>
          <span className={styles.total}>
            <span>剩余可用精力：</span>
            <span
              className={
                100 - total < 0 ? `${styles.warning} ${styles.num}` : styles.num
              }
            >
              {100 - total}%
            </span>
          </span>
        </div>
        <Loading visible={depts.length} loading={true}>
          {!!items.length &&
            +isSubmit === 1 &&
            timeScope === format(new Date(), "yyyy-MM-dd") && (
              <NoticeBar content="今日还未提交日报" color="alert" closeable />
            )}
          <Collapse
            className={styles.collapse}
            defaultActiveKey={Object.keys(depts)}
          >
            {depts.map((dept: any, index: number) => {
              const [item, itemIndex] = getItem(dept);

              return (
                <Collapse.Panel key={`${index}`} title={dept.name}>
                  <ProjectOrWork
                    onChange={(item: any) => {
                      items[itemIndex] = item;
                      setItems([...items]);
                    }}
                    timeScope={timeScope}
                    deptCode={dept.deptCode}
                    item={item}
                  />
                </Collapse.Panel>
              );
            })}
          </Collapse>
          <div className={styles.placeholder}></div>
          <SafeArea position="bottom" />
        </Loading>
      </div>
      <Loading visible={depts.length} laoding={false}>
        <div className={styles.footer}>
          <Checkbox
            style={{
              "--icon-size": "18px",
              "--font-size": "14px",
              "--gap": "6px",
              margin: "0 0 20px 0",
              textAlign: "center",
              color: isAutoFill ? "#666666" : "#999999",
            }}
            checked={isAutoFill}
            onChange={setIsAutoFill}
          >
            开启一周自动提交日报
          </Checkbox>
          <div style={{ display: "flex" }}>
            <Button
              className={`${styles.button} ${styles.look}`}
              block
              onClick={() => {
                // history.push(`/daily/statics${location.search}`);
              }}
            >
              日报记录
            </Button>
            <Button
              className={styles.button}
              block
              onClick={onSubmit}
              style={
                +isSubmit === 2 && timeScope === now
                  ? { backgroundColor: "#00b578" }
                  : {}
              }
              disabled={total > 100 || (+isSubmit === 2 && timeScope !== now)}
            >
              {total <= 100 ? "提交" : "工作精力超过 100%"}
            </Button>
          </div>

          <SafeArea position="bottom" />
        </div>
      </Loading>
      <div id="simulator-popup"></div>
    </div>
  );
};

export default Report;
