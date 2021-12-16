import React from "react";
import { createPortal } from "react-dom";
import Report from "../modules/Report";
import DashBoard from "../modules/Dashboard";
import Card from "@/components/Card";
import "./style.less";
export class FramePortal extends React.PureComponent {
  containerEl: HTMLElement;
  iframe?: HTMLIFrameElement | null;

  constructor(props: {}) {
    super(props);
    this.containerEl = document.createElement("div");
  }

  render() {
    return (
      <iframe title="iframe" ref={(el) => (this.iframe = el)}>
        {createPortal(this.props.children, this.containerEl)}
      </iframe>
    );
  }

  componentDidMount() {
    this.iframe!.contentDocument!.body.appendChild(this.containerEl);
  }
}

const Home: React.FC<{}> = (props) => {
  const ref = React.createRef<HTMLIFrameElement>();
  const [loaded, setLoaded] = React.useState(false);

  React.useLayoutEffect(() => {
    const mountNode = ref.current?.contentWindow?.document?.body;
  }, []);

  return (
    <div className="home-page">
      <Card title="一分钟日报初体验" fill="#f0f2f5"></Card>
      <div className="simulator-box">
        <Report />
      </div>
      <div className="dashboard-box">
        <DashBoard />
      </div>
    </div>
  );
};

export default Home;
