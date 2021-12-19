
export type OptionType = { label: string, value: string }

const cache: {
  [key: string]: OptionType[];
} = {};

export function loadCache(_key: string) {
  const key = `_tree_cache_${_key}`;
  if (!!localStorage.getItem(key)) {
    cache[key] = JSON.parse(localStorage.getItem(key)!);
  } else {
    cache[key] = [];
  }
  return cache[key];
}

export function updateCache(_key: string, value: OptionType) {
  const key = `_tree_cache_${_key}`;
  const list = cache[key];
  const old = list.findIndex(item => item.value === value.value);

  if (old >= 0) {
    list.splice(old, 1);
  }

  list.unshift(value);

  if (list.length > 10) {
    list.pop();
  }

  cache[key] = list;
  localStorage.setItem(key, JSON.stringify(list));
}