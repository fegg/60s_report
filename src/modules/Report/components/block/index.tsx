import React from 'react';
import styles from './index.module.scss';

function Block(props: any) {
  const {
    children,
    label,
    extra,
    style = {},
    className = '',
  }: { children: any; extra: any; label?: string; style?: any; className?: string } = props;
  return (
    <div className={`${styles.block} ${className}`}>
      {label && (
        <div className={styles.label} style={style}>
          <span>{label}</span>
          {extra ? extra : ''}
        </div>
      )}

      {typeof children === 'function' ? children(props) : children}
    </div>
  );
}

export default Block;
