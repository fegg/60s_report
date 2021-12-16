import { Toast, Dialog } from "antd-mobile";
import { AddCircleOutline } from "antd-mobile-icons";
import React, { useContext } from "react";
import Block from "../block";
import CardLayout from "./components/layout";
import styles from "./components/index.module.scss";
import { CardContext } from "./components/context";
/**
 *
 * @param props 卡片组件
 * @returns
 */
function Card(props: any) {
  const { deptCode, item, onChange, projectType, title, type } = props;
  let currentData: any = [],
    otherData: any = [],
    projects: any = [],
    projectsMap: any = {},
    works: any = [],
    worksMap: any = {},
    getProjects: any = () => {},
    searchParams: any = {};

  return (
    <CardContext.Provider
      value={{
        getProjects,
        projects,
        projectsMap,
        works,
        worksMap,
        searchParams,
      }}
    >
      <Block
        label={title}
        className={currentData.length ? styles.block : ""}
        extra={
          <div
            onClick={() => {
              let data: any = {
                time: +new Date(),
                cost: 0,
              };
              // 如果是工作事项需要按照类型区分数据
              if (type === "item") {
                data["projectType"] = projectType;
                if (currentData[0] && !currentData[0]["projectId"]) {
                  Toast.show({
                    content: "请先完善再新增",
                  });
                  return;
                }
              } else {
                if (currentData[0] && !currentData[0]["projectId"]) {
                  Toast.show({
                    content: "请先完善再新增",
                  });
                  return;
                }
              }
              // 如果第一项数据的 projectId 存在才能插入数据
              item[`${type}s`].unshift(data);
              onChange(item);
            }}
          >
            <AddCircleOutline style={{ color: "#1677ff" }} />
          </div>
        }
      >
        {!!currentData.length &&
          currentData.map((data: any, index: number) => (
            <CardLayout
              type={type}
              key={`${type}_${index}_${data["projectId"]}_${data.time}_${
                item[`${type}Type`]
              }`}
              data={data}
              projectType={projectType}
              deptCode={deptCode}
              onChange={(value: any) => {
                // 验证修改的数据是当前项目/或者工作事项时判断是否在项目中已经添加
                const info = currentData.filter(
                  (item: any) => item["projectId"] === value["projectId"]
                );
                if (value["projectId"] && info.length) {
                  // 如果重复选择当前项目/工作不提示
                  if (value["projectId"] === data["projectId"]) return;
                  Toast.show({
                    content: "当前工作已添加请勿重复添加",
                  });
                  console.log(data);
                  return;
                }

                //

                /**
                 * isRequiredComment * 评论是否必填,1必填,2非必填
                 * isForceNotice * 是否强制提醒用户必填,1必填,2非必填
                 * notice * 用户提醒文案
                 */
                if (value.isForceNotice === 1) {
                  Dialog.alert({
                    content: value.notice || "请填写备注",
                    closeOnMaskClick: false,
                  });
                }

                data = {
                  costDeptCode: deptCode,
                  ...data,
                  ...value,
                };
                if (type === "item") {
                  data["projectType"] = projectType;
                }

                currentData[index] = data;
                item[`${type}s`] = [...currentData, ...otherData];
                onChange(item);
              }}
              onRemove={() => {
                currentData.splice(index, 1);
                item[`${type}s`] = [...currentData, ...otherData];
                onChange(item);
              }}
            />
          ))}
      </Block>
    </CardContext.Provider>
  );
}

export default Card;
