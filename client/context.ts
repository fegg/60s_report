import { createContext } from 'react';
import { IProjectListItem } from '@/types/IProjectItem';

const AppContext = createContext<any>({
  total: 0,
  timeScope: '',
  recentOptions: {
    projectItems: [],
    workItems: [],
  },
  myProjects: [],
});

const CardContext = createContext<{
  works: IProjectListItem[];
  projects: any;
  projectsMap: any;
  getProjects: any;
  worksMap: any;
  searchParams: any;
}>({
  projects: [],
  projectsMap: {},
  works: [],
  worksMap: {},
  getProjects: () => {},
  searchParams: {},
});
export { AppContext, CardContext };
