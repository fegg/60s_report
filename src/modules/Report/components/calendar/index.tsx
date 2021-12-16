import React, { forwardRef, useEffect, useState } from "react";
import "./index.scss";

interface Props {
  timeScope?: string;
  onChange?: any;
}
function PopCalendar(props: Props, ref: any) {
  const { timeScope = "", onChange } = props;

  return <span>{timeScope}</span>;
}

export default forwardRef(PopCalendar);
