import React from "react";


type ContextType = {
  recentOptions: { [key: string]: string },
  timeScope: string,
  total: number,
}

export const AppContext = React.createContext<ContextType>({
  recentOptions: {

  },
  timeScope: "",
  total: 0,
});

interface WorkType {
  label: string;
  value: string;
}

interface WorkMapType {
  projectName: string;
  projectId: string;
  notice: string;
  costDeptCode: string;
  isForceNotice: string;
  isRequiredComment: string;
}

export const CardContext = React.createContext<{
  works: WorkType[],
  projects: WorkMapType[],
  worksMap: WorkMapType[],
  projectsMap: Record<string, string>,
  getProjects: () => void,
  searchParams: {
    projectName: string
  }
}>({
  works: [],
  projects: [],
  worksMap: {},
  projectsMap: {},
  getProjects: () => { },
  searchParams: {
    projectName: ""
  }
});