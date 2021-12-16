import { canUseDom } from './can-use-dom';

export let supportsPassive = false;

if (canUseDom) {
  try {
    const opts = {};
    Object.defineProperty(opts, 'passive', {
      // eslint-disable-next-line getter-return
      get() {
        supportsPassive = true;
      },
    });
    window.addEventListener('test-passive', null as any, opts);
  } catch (e) {}
}
