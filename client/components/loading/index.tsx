import { Loading } from 'antd-mobile';
import React from 'react';
import styles from './index.module.scss';

function LoadingLayout(props: any) {
  const { visible, loading, children }: any = props;
  return visible ? (
    children
  ) : loading ? (
    <div className={styles.loading}>
      <Loading />
    </div>
  ) : (
    ''
  );
}

export default LoadingLayout;
