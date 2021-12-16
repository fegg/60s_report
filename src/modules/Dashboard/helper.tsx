import React from "react";
import MedCard, { CardProgress, CardStatics } from "@/components/Card";
import { colors, factors } from "@/constant/application";
import { getFixedValue } from "@/utils";

const more_style: React.CSSProperties = {
  minWidth: "216px",
  minHeight: "73px",
  backgroundColor: "#f1f1f1",
  borderRadius: "10px",
  // 'padding': '16px',
  color: "#2a2a2a",
  cursor: "pointer",
};

type DeptCardItemProps = {
  title: string;
  factor: string;
  ratio?: string;
  topRatio?: number;
  fill?: string;
};

type CardInnerType = {
  title: string;
  ratio: string;
  desc: string;
  index?: number;
  clickable?: boolean;
};

const factorNames = factors.map(([name]) => name);

/**
 * 调用方式
 * const Comp = makeCardData(datas);
 *
 * <Comp title="日常" fator="matter" />
 * @param dataset 总数据集
 * @returns
 */
export const makeCardData =
  (dataset: any, showDetail: any, callbacks?: {}) =>
  (options: DeptCardItemProps) => {
    const { title, factor } = options;
    const { list, cost, total, ratio } = getFactorData(factor, dataset);
    const { data, topRatio = 0, more } = _getTopList(list, cost, total);
    const fill = options.fill || colors[factorNames.indexOf(factor)];

    function Wrapper(item: CardInnerType) {
      return (
        <CardStatics
          title={item.title}
          key={title}
          quota={item.ratio}
          desc={item.desc}
          clickable={item.clickable ? true : false}
          onClick={() => {
            // @ts-ignore
            callbacks?.itemClick(item, item.index, factor);
          }}
          fill={fill}
        ></CardStatics>
      );
    }

    const dom = data.map((item: CardInnerType, index: number) => {
      return (
        <Wrapper {...item} key={item.title} index={index} clickable={true} />
      );
    });

    function renderMore() {
      if (more) {
        const otherRatio = getFixedValue(100 - topRatio);
        const otherRatioTotal = getFixedValue((otherRatio * cost) / total);
        return (
          <>
            <Wrapper
              title="其余总计"
              ratio={`${otherRatioTotal}%`}
              desc=""
              clickable={false}
            ></Wrapper>
            <MedCard
              style={more_style}
              onClick={() => showDetail({ list, total, title })}
              key="more"
            >
              查看所有
            </MedCard>
          </>
        );
      }
      return null;
    }

    return dom.length > 0 ? (
      <CardProgress
        title={title}
        progress={ratio}
        bodyStyle={{ display: "flex" }}
      >
        {dom}
        {renderMore()}
      </CardProgress>
    ) : (
      <CardProgress
        title={title}
        progress={ratio}
        bodyStyle={{ display: "flex" }}
      />
    );
  };

const _getTopList = function (list: any[], total: number, all: number) {
  const limit = 3;
  // 统计前3项目的百分比
  let ratio = 0;
  const top = list.slice(0, limit).reduce(function (acc, item) {
    const currentRatio = getFixedValue((item.cost / total) * 100);
    const totalRatio = getFixedValue((item.cost / all) * 100);
    ratio += Number(currentRatio) * 1;
    return [...acc, { title: item.title, ratio: totalRatio + "%", desc: `` }];
  }, []);

  return {
    topRatio: getFixedValue(ratio),
    data: top,
    more: list.length > limit,
  };
};

// export function fixedNumber(value: number, length =1){ return Number(value.toFixed(length)); }

function getFactorData(factor: any, base: any) {
  const key = `${factor}Info`;
  const list = base[key].list;
  const total = Number(base.costInfo.total);
  const cost = Number(base.costInfo[`${factor}Cost`]);
  const ratio = Boolean(cost) ? getFixedValue((cost / total) * 100) : 0;

  return { list, cost, total, ratio };
}
