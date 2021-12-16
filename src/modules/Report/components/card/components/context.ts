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

export const CardContext = React.createContext<{
  works: string[],
  projects: string[],
  worksMap: Record<string, string>,
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