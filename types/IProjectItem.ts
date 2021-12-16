export interface IIProjectItem {
    list: IProjectsListItem[];
    total: string;
}
export interface IProjectsListItem {
    parentId: string;
    projectId: string;
    userDeptCode: string;
    costDeptCode: string;
    costDeptName: string;
    projectName: string;
    projectType: string;
    remark: string;
    isRequiredComment: string;
    isForceNotice: string;
    notice: string;
    hasChildrenProject: boolean;
    createdAt: string;
    updatedAt: string;
    children: IProjectChildrenItem[];
}
interface IProjectChildrenItem {
    parentId: string;
    projectId: string;
    userDeptCode: string;
    costDeptCode: string;
    costDeptName: string;
    projectName: string;
    projectType: string;
    remark: string;
    isRequiredComment: string;
    isForceNotice: string;
    notice: string;
    hasChildrenProject: boolean;
    createdAt: string;
    updatedAt: string;
    children: any[];
}

