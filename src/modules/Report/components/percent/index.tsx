import Slider from '@/components/slider/index';
import React from 'react';
import styles from './index.module.scss';

const marks = Array.from({ length: 21 }, (v, k) => k * 5).reduce((total: any, v: any) => {
  total[v] = v === 0 || v === 100 ? v : '';
  return total;
}, {});

function Percent({ onChange, onAfterChange, value }: any) {
  return (
    <div className={styles.sliderWrapper}>
      <Slider value={value} onChange={onChange} onAfterChange={onAfterChange} marks={marks} />
    </div>
  );
}

export default Percent;
