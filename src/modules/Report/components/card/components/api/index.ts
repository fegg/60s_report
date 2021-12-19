import data from "./api.json";

function makeWorkMap(res: any) {
  const maps: any = {};

  function foreachList(items: any) {
    return items.map((item: any) => {
      let data: any = {
        ...item,
        label: item.projectName,
        value: item.projectId,
      };
      if (item.children.length) {
        data['children'] = foreachList(item.children);
      }
      maps[item.projectId] = item;
      return data;
    });
  }
  const items = foreachList(res.list || []);

  return [items, maps] as const;
}

export default makeWorkMap(data.data);