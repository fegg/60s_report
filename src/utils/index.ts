export function fixedNumber(value: number, length = 1) {
  return Number(value.toFixed(length));
}

export const getFixedValue = (value: number, fixed = 1) => {
  // @ts-ignore
  if (isNaN(value) || value === Infinity) return 0;
  const reg = /0.0*[^0]/;
  if (value > 0 && value < 0.1) {
    const r = value.toString().match(reg);
    if (r) return Number(r[0]);
  }
  return Number(value.toFixed(fixed));
};

export function isMobile() {
  let info = navigator.userAgent;
  let agents = ['Android', 'iPhone', 'SymbianOS', 'Windows Phone', 'iPod', 'iPad'];
  for (let i = 0; i < agents.length; i++) {
    if (info.indexOf(agents[i]) >= 0) return true;
  }

  // 新增 dd 平台判断
  // @ts-ignore
  if (window.dd && dd.env.platform !== 'notInDingTalk') {
    return true;
  }
  return false;
}