import React from 'react';
import Card from './card';

/**
 * 项目组件：包含项目 / 工作日常
 * @returns 项目组件
 */
function ProjectOrWork(props: any) {
  return (
    <div>
      <Card type="project" title="项目投入" {...props} />
      <Card type="item" title="日常工作" {...props} projectType={3} />
      {/* <Card type="item" title="日常管理" {...props} projectType={4} /> */}
    </div>
  );
}

export default ProjectOrWork;
