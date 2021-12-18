import React from "react";
import ReactDOM from "react-dom";
import Home from "./pages/index";

export const mobileInit = () => {
  (function () {
    function pxTorem() {
      //html:50px: 15rem 750px
      const base = 50; //rem
      let wid = window.innerWidth || document.body.clientWidth;
      wid > 750 && (wid = 750);
      const size = wid / (750 / base);
      // @ts-ignore
      document.querySelector("html").style.fontSize = `${48}px`;
    }
    pxTorem();
    window.addEventListener("resize", function () {
      pxTorem();
    });
  })();
};

export const App: React.FC<{}> = () => {
  React.useEffect(() => {
    mobileInit();
  }, []);
  return <Home></Home>;
};

ReactDOM.render(<App />, document.getElementById("app"));
