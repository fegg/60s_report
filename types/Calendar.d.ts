export interface ICalendar {
    list: IListItem[];
}
interface IListItem {
    items: IItemsItem[];
    date: string;
}
interface IItemsItem {
    id: string;
    projectId: string;
    type: string;
    projectRankType: string;
    title: string;
    remark: string;
    level: string;
    item: string;
    finishTime: string;
    isNeedReport: string;
    reportAid: string;
    needChoiceItems: string;
    needFocusItems: string;
    risks: any[];
    projectInfo: IProjectInfo;
}
interface IProjectInfo {
    id: string;
    title: string;
    projectId: string;
    projectRankType: string;
    projectLevel: string;
    managerId: string;
    managerName: string;
    projectLink: string;
    costDeptCode: string;
    costDeptName: string;
    focusUser: any[];
    startTime: string;
    finishTime: string;
    authControl: string;
    risks: any[];
    hasChildrenProject: boolean;
    hasNode: boolean;
    parentId: string;
    checkUser: any[];
    checkStdState: string;
    checkContent: string;
    parentName: string;
    officerList: any[];
}

