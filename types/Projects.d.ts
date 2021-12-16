export interface IProjects {
    list: IProjectListItem[];
    cost: number;
}
export interface IProjectListItem {
    projectId: string;
    projectName: string;
    managerName: string;
    startTime: string;
    endTime: string;
    costDeptName: string;
    projectLevel: string;
    projectRankType: string;
    cost: string;
    ratio: number;
    hasChildrenProject: boolean;
    deptProjectStatics: IDeptProjectStaticsItem[];
    deptProjectStaticsMap: IDeptProjectStaticsMap;
}
interface IDeptProjectStaticsItem {
    deptCode: string;
    name: string;
    cost: string;
    ratio: number;
    user: IProjectUserItem[];
}
interface IProjectUserItem {
    userId: string;
    name: string;
    cost: string;
}
interface IDeptProjectStaticsMap {
}

