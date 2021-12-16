export interface IProjectList {
  list: IListItem[];
  total: string;
  page: string;
}
interface IListItem {
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
  focusUser: IFocusUserItem[];
  startTime: string;
  finishTime: string;
  authControl: string;
  risks: any[];
  hasChildrenProject: boolean;
  hasNode: boolean;
  parentId: string;
  checkUser: ICheckUserItem[];
  parentName: string;
  officerList: IOfficerListItem[];
  commentUserList: ICommentUserListItem[];
  evaluateUser: null | IEvaluateUser;
  intro: string;
  result: string;
  projectSubject: null | IProjectSubject;
  checkList: ICheckListItem[];
  curCoin: string;
  isWaterProject: string;
  medCoinTotal: string;
  isWarning: boolean;
  subTotalScore: string;
  projectScore: string;
}
interface IFocusUserItem {
  userId: string;
  name: string;
  email: string;
}
interface IOfficerListItem {
  userId: string;
  name: string;
  email: string;
}
interface ICheckUserItem {
  userId: string;
  name: string;
  checkState: string;
  reason: string;
  email: string;
}
interface IEvaluateUser {
  userId: string;
  name: string;
  email: string;
}
interface ICheckListItem {
  id: string;
  checkTime: string;
  percent: string;
  checkContent: string;
  isDelete: boolean;
  lastCheckSuccessContent: string;
  checkState: string;
  verifyState: string;
  role: string;
  isComment: boolean;
  result: string;
  checkers: ICheckersItem[];
}
interface ICheckersItem {
  userId: string;
  name: string;
  checkState: string;
  reason: string;
}
interface IProjectSubject {
  ProjectSubjectId: string;
  projectSubjectName: string;
  status: string;
  checkGroups: ICheckGroupsItem[];
  groups: IGroupsItem[];
  groupScore: IGroupScoreItem[];
  description: string;
  url: string;
  createdAt: string;
}
interface ICheckGroupsItem {
  projectSubjectId: string;
  checkUser: ICheckUser;
  weight: string;
}
interface ICheckUser {
  userId: string;
  name: string;
  email: string;
}
interface IGroupsItem {
  projectSubjectId: string;
  groupId: string;
  groupName: string;
  captainUser: ICaptainUser | null;
  remark: string;
}
interface ICaptainUser {
  userId: string;
  name: string;
  email: string;
}
interface IGroupScoreItem {
  projectSubjectId: string;
  groupId: string;
  groupName: string;
  score: string;
  scoreList: IScoreListItem[];
}
interface IScoreListItem {
  projectSubjectId: string;
  groupId: string;
  checkUser: ICheckUser;
  score: string;
  weight: string;
}
interface ICommentUserListItem {
  userId: string;
  name: string;
  email: string;
}

