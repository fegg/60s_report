export interface IDeptTrend {
    dayData: {
        day: string;
        deptList: IDeptListItem[]
    }[]
}

export interface IDepartments {
    total: string;
    num: string;
    deptList: IDeptListItem[];
}
interface IDeptListItem {
    deptName: string;
    deptCode: string;
    deptInfo: IDeptInfo;
}
export interface IDeptInfo {
    costInfo: ICostInfo;
    companyProjectInfo: ICompanyProjectInfo;
    deptProjectInfo: IDeptProjectInfo;
    matterInfo: IMatterInfo;
    managerInfo: IManagerInfo;
    lostInfo: ILostInfo;
    deptInfo: any[];
}
interface ICostInfo {
    loss: string;
    managerCost: string;
    companyProjectCost: string;
    deptProjectCost: string;
    matterCost: string;
    total: string;
}
interface ICompanyProjectInfo {
    list: IDeptListItem[];
}
interface IDeptProjectInfo {
    list: IDeptListItem[];
}
interface IMatterInfo {
    list: IDeptListItem[];
}
interface IDeptListItem {
    title: string;
    cost: string;
    user: IUserItem[];
    type?: string;
    projectId?: string;
    children?: IChildrenItem[];
}
interface IUserItem {
    userId: string;
    name: string;
    cost: string;
}
interface IManagerInfo {
    list: IDeptListItem[];
}
interface ILostInfo {
    user: IUserItem[];
}
interface IChildrenItem {
    title: string;
    projectId: string;
    cost: string;
    isParentSelf: boolean;
    user: IUserItem[];
}

