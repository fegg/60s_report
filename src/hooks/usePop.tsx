import React from "react";
import { Popup } from "antd-mobile";
/**
 * 弹出层 hooks
 */
function usePop(
  title: string,
  forceRender = false
): [React.FC<React.HTMLAttributes<HTMLDivElement>>, () => void, any] {
  const [show, setShow] = React.useState(false);
  const [showText, setShowText] = React.useState(title);

  const Wrapper: React.FC<{ title?: string }> = function (
    props
  ): React.ReactElement {
    return (
      <>
        <div
          className="filter-item-title"
          onClick={() => {
            setShow(true);
          }}
        >
          {showText || title}
        </div>
        <Popup
          visible={show}
          onMaskClick={() => setShow(false)}
          position="bottom"
        >
          {props.children}
        </Popup>
      </>
    );
  };

  return [
    React.memo(Wrapper, () => forceRender),
    () => setShow(false),
    setShowText,
  ];
}

export const usePopPicker = function (title: string, forceRender = false) {
  const [show, setShow] = React.useState(false);
  const [showText, setShowText] = React.useState(title);

  const Wrapper: React.FC<{
    renderPicker: () => React.ReactNode;
    title?: string;
  }> = function (props): React.ReactElement {
    return (
      <>
        <div
          className="filter-item-title"
          onClick={() => {
            setShow(true);
          }}
        >
          {props.title || showText || title}
        </div>
        <Popup
          visible={show}
          onMaskClick={() => setShow(false)}
          position="bottom"
        >
          {props.renderPicker()}
        </Popup>
      </>
    );
  };

  return [
    React.memo(Wrapper, () => forceRender),
    () => setShow(false),
    setShowText,
  ];
};

export default usePop;
