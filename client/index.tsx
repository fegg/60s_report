import { Button, SafeArea, Toast, NoticeBar, Checkbox } from 'antd-mobile';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import moment from 'moment';
import { getPersonDept, getProjectTaskDetail, postProjectTask } from '@/api';
import { QuestionCircleOutline } from 'antd-mobile-icons';
import Collapse from '@/components/collapse';
import * as dd from 'dingtalk-jsapi';
import { setTitle } from '@/utils/dd';
import { mobileInit } from '@/utils';
import { useHistory } from '@medlinker/fundamental';
import { getUserPermanentProject, getUserRecentOption } from '@/api/daily';
import ProjectOrWork from './components/index';
import Loading from './components/loading';
import PopCalendar from './components/calendar';
import { AppContext } from './context';
import styles from './index.module.scss';
import './index.scss';

function Daily() {
  const history = useHistory();
  const [items, setItems] = useState<any>([]);
  const [timeScope, setTimeScope] = useState(moment(+new Date()).format('YYYY-MM-DD'));
  const [depts, setDepts] = useState([]);
  const [total, setTotal] = useState<number>(0);
  const [pending, setPending] = useState<boolean>(false);
  const [isSubmit, setIsSubmit] = useState<any>();
  const [isAutoFill, setIsAutoFill] = useState<boolean>(false);
  const [recentOptions, setRecentOptions] = useState<any>([]);
  const [myProjects, setMyProjects] = useState<any>([]);
  const ref = useRef();

  const updateTask = useCallback(() => {
    return getProjectTaskDetail({ timeScope })
      .then((res) => {
        res.list = (res.list || []).map((item: any) => {
          let items: any = [],
            projects: any = [];
          item.projects.forEach((project: any) => {
            // 0项目 1日常事项 2日常管理
            switch (+project.projectType) {
              case 1:
              case 2:
                projects.push(project);
                break;
              case 3:
              case 4:
                items.push(project);
                break;
            }
          });
          return {
            items,
            projects,
            userDeptCode: item.userDeptCode,
          };
        });
        return res;
      })
      .then((res) => {
        setItems(res.list);
        setIsSubmit(res.isSubmit);
        setIsAutoFill(!!res.isAutoFill);
      });
  }, [timeScope]);

  useEffect(() => {
    setTitle('医联一分钟日报');
    mobileInit();

    // 判断当前地址栏有没有日期
    // 通过日期查询日常信息
    getPersonDept({}).then((res) => {
      setDepts(res.list || []);
    });
    updateTask();
  }, [timeScope, updateTask]);

  useEffect(() => {
    const total = items.reduce((total: number, data: any) => {
      if (data) {
        total += data.projects.reduce((total: number, item: any) => {
          total += +item.cost || 0;
          return total;
        }, 0);
        total += data.items.reduce((total: number, item: any) => {
          total += +item.cost || 0;
          return total;
        }, 0);
      }
      return total;
    }, 0);
    setTotal(total);
  }, [items]);

  // 通过部门信息匹配部门对应的日常信息
  const getItem = useCallback(
    (dept: any) => {
      const index = items.findIndex((item: any = {}) => item.userDeptCode === dept.deptCode);
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

  function getButtonName() {
    if (+isSubmit === 1) {
      if (timeScope !== moment().format('YYYY-MM-DD')) {
        return '补填';
      } else {
        return '提交';
      }
    } else if (+isSubmit === 2) {
      if (timeScope !== moment().format('YYYY-MM-DD')) {
        return `${timeScope} 的数据已填写`;
      }
      return '更新';
    } else {
      return '提交';
    }
  }

  // 验证数据
  function validate(data: any) {
    let newItems: any = [];
    let errorInfo = '';
    data.forEach((item: any) => {
      if (item) {
        item.items = item.items.filter((item: any) => {
          const itemTypes: any = {
            '3': '工作',
            '4': '管理',
          };
          if (!item.projectId) {
            errorInfo = `请选择新增的日常${itemTypes[item.projectType]}`;
          } else if (item.isRequiredComment === 1 && !item.remark) {
            errorInfo = item.notice;
          } else {
            delete item.notice;
            delete item.isRequiredComment;
            delete item.isForceNotice;
            delete item.time;
          }
          return item;
        });
        item.projects = item.projects.filter((item: any) => {
          if (!item.projectId) {
            errorInfo = `请选择新增的项目`;
          }
          delete item.time;
          return item;
        });

        newItems.push({
          projects: [...item.items, ...item.projects],
          userDeptCode: item.userDeptCode,
        });
      }
    });
    if (errorInfo) {
      Toast.show({
        content: errorInfo || '请填写备注',
      });
      return false;
    }
    return newItems;
  }

  // 提交前验证数据
  function onSubmit() {
    if (total === 0) {
      Toast.show({
        content: '请完善个人工作填报',
      });
      return;
    }
    const data = validate(items);
    if (!data) {
      return;
    }
    console.log(data);
    if (pending) {
      return;
    }
    setPending(true);
    Toast.show({
      icon: 'loading',
      duration: 0,
    });
    postProjectTask({
      list: data,
      isAutoFill,
      timeScope,
    })
      .then(() => {
        Toast.clear();
        Toast.show({
          icon: 'success',
          content: '提交成功',
          afterClose() {
            dd.biz.navigation.close({}).catch((e) => {
              console.log(e.message);
            });
            // @ts-ignore
            ref.current.getUserFillLogCallback(timeScope);
            updateTask()
              .then(() => {
                setPending(false);
              })
              .catch(() => {
                setPending(false);
              });
          },
        });
      })
      .catch(() => {
        setPending(false);
        Toast.clear();
      });
  }

  useEffect(() => {
    // 查询最近选择
    getUserRecentOption().then(({ projectItems = [], workItems = [] }: any) => {
      projectItems = (projectItems || []).map((item: any) => ({
        label: item.projectName,
        value: +item.parentId === 0 ? `|${item.id}` : item.id,
        projectType: item.projectRankType,
      }));
      workItems = (workItems || []).map((item: any) => ({
        label: item.projectName,
        value: +item.parentId === 0 ? `|${item.id}` : item.id,
      }));
      setRecentOptions({
        projectItems,
        workItems,
      });
    });

    // 查询我的常驻项目
    getUserPermanentProject({
      page: 1,
      limit: 999,
    }).then(({ list = [] }) => {
      list = list.map((item: any) => ({
        label: item.projectName,
        value: +item.parentId === 0 ? `|${item.id}` : item.id,
        projectType: item.projectRankType,
      }));
      setMyProjects(list);
    });
  }, []);

  return (
    <AppContext.Provider value={{ total, timeScope, recentOptions, myProjects }}>
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
            </div>
            <span className={styles.total}>
              <span>剩余可用精力：</span>
              <span className={100 - total < 0 ? `${styles.warning} ${styles.num}` : styles.num}>
                {100 - total}%
              </span>
            </span>
          </div>
          <div className={styles['header-placeholder']}></div>
          <Loading visible={depts.length} loading={true}>
            {!!items.length && +isSubmit === 1 && timeScope === moment().format('YYYY-MM-DD') && (
              <NoticeBar content="今日还未提交日报" color="alert" closeable />
            )}
            <Collapse className={styles.collapse} defaultActiveKey={Object.keys(depts)}>
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
                '--icon-size': '18px',
                '--font-size': '14px',
                '--gap': '6px',
                margin: '0 0 20px 0',
                textAlign: 'center',
                color: isAutoFill ? '#666666' : '#999999',
              }}
              checked={isAutoFill}
              onChange={setIsAutoFill}
            >
              开启一周自动提交日报
            </Checkbox>
            <div style={{ display: 'flex' }}>
              <Button
                className={`${styles.button} ${styles.look}`}
                block
                onClick={() => {
                  history.push(`/daily/statics${location.search}`);
                }}
              >
                日报记录
              </Button>
              <Button
                className={styles.button}
                block
                onClick={onSubmit}
                style={
                  +isSubmit === 2 && timeScope === moment().format('YYYY-MM-DD')
                    ? { backgroundColor: '#00b578' }
                    : {}
                }
                disabled={
                  total > 100 ||
                  pending ||
                  (+isSubmit === 2 && timeScope !== moment().format('YYYY-MM-DD'))
                }
              >
                {total <= 100 ? getButtonName() : '工作精力超过 100%'}
              </Button>
            </div>

            <SafeArea position="bottom" />
          </div>
        </Loading>
      </div>
    </AppContext.Provider>
  );
}

export default Daily;
