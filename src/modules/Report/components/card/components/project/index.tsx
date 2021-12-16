import { Popup, Search } from "antd-mobile";
import React, { useState, useContext, useEffect, useCallback } from "react";
import { RightOutline } from "antd-mobile-icons";
import { AppContext, CardContext } from "../context";
import { debounce } from "lodash";
import styles from "../index.module.scss";
import TreeSelect from "../work/components/TreeSelect";
import "./index.scss";

function Project(props: any) {
  const { projects, projectsMap, getProjects, searchParams } =
    useContext(CardContext);
  const { timeScope } = useContext(AppContext);
  const [panelVisible, setPanelVisible] = useState(false);
  const [state, setState] = useState<any>([]);

  const { data = {}, onChange } = props;
  const func = debounce(getProjects, 800);
  const onSearchChange = useCallback(
    (value?: string) => {
      let params: any = { page: 1, pageSize: 999, searchDate: timeScope };

      if (value) {
        params.projectName = value;
      }

      func(params);
    },
    [func, timeScope]
  );

  useEffect(() => {
    let newProjects = [...projects];
    // -2 公司项目 -3 部门项目
    // 1. 按照公司项目 / 部门项目区分 projectRankType
    // 2. 按照项目等级区分 S A B C / 4 3 2 1 projectLevel
    let companyProject: any = [],
      departmentProject: any = [];
    Object.values(projectsMap).forEach((item: any) => {
      // 有子项目的父项目本身不能被选择
      // if (!item.hasChildrenProject) {
      const { children = [], ...data } = item;
      +item.projectRankType === 1 &&
        companyProject.push({ ...data, value: `~${data.value}` });
      +item.projectRankType === 2 &&
        departmentProject.push({ ...data, value: `~${data.value}` });
      // }
    });

    function formatLabel(label: string) {
      if (typeof label === "string") {
        return <span dangerouslySetInnerHTML={{ __html: label }}></span>;
      } else {
        return label;
      }
    }

    newProjects = newProjects.map((item: any) => {
      const label = formatLabel(item.label);
      if (item.children.length) {
        item.children = item.children.map((subItem: any) => {
          const label = formatLabel(subItem.label);
          return {
            ...subItem,
            label,
          };
        });
        const newItem = {
          ...item,
          label,
          value: `!${item.value}`,
          children: [{ ...item, children: [] }, ...item.children],
        };
        return newItem;
      }
      return { ...item, label };
    });
    // 父项目
    if (searchParams.projectName && projects.length) {
      // @ts-ignore
      document.querySelectorAll(".adm-tree-select-item")[3]?.click();
    }
    setState([
      { label: "公司级项目", value: "-2", children: companyProject },
      { label: "部门级项目", value: "-3", children: departmentProject },
      ...newProjects,
    ]);
  }, [projects, projectsMap, searchParams.projectName]);

  useEffect(() => {
    !data.projectId && setPanelVisible(true);
  }, [data]);

  return (
    <>
      <span
        className={styles.ellipsis}
        onClick={() => {
          setPanelVisible(true);
        }}
      >
        {data.projectName || <span className={styles.unselect}>请选择</span>}
        <RightOutline />
      </span>
      <Popup
        visible={panelVisible}
        onMaskClick={() => {
          setPanelVisible(false);
        }}
        className={styles.projectPopup}
        bodyStyle={{ borderRadius: "20px 20px 0 0" }}
      >
        <div className={styles.projectHeader}>
          项目名称
          <img
            onClick={() => setPanelVisible(false)}
            className={styles.projectPopupClose}
            src="https://pub-med-post.medlinker.com/_/prod/developer-panel/1462974496569823003.png"
          />
        </div>
        <Search
          onChange={onSearchChange}
          style={{
            backgroundColor: "#ffffff",
            width: "100%",
            padding: "10px 20px",
          }}
          placeholder="请输入项目名称过滤"
          defaultValue={searchParams?.projectName}
        />
        <div className={styles.checkList} style={{ height: "80vh" }}>
          <TreeSelect
            cacheKey="project_tree"
            options={state ?? projects}
            defaultValue={["-1", data.projectId]}
            onChange={(value) => {
              const item = projectsMap[value[1].replace("~", "")];

              let newItem: any = {
                projectType: item.projectType,
                projectName: item.label,
                projectId: item.value,
              };
              if (item.costDeptCode) {
                newItem["costDeptCode"] = item.costDeptCode;
              }
              if (newItem.projectId !== data.projectId) {
                newItem["remark"] = "";
              }
              onChange(newItem);
              setPanelVisible(false);
            }}
          />
        </div>
      </Popup>
    </>
  );
}

export default Project;
