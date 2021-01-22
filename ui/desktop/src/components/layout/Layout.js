import React from "react";
import Sidebar from "../nav/Sidebar";

import { WebSocket } from "../../features/websocket/WebSocket";
import AboutApplication from "../../features/config/AboutApplication";

const Layout = (props) => {
  return (
    <section className="section">
      <WebSocket />
      <div className="columns is-multiline">
        <div className="column is-3 is-narrow-mobile is-fullheight">
          <Sidebar />
        </div>
        <div className="column is-9 pb-0 mb-5">{props.children}</div>
        <div className="column is-12 mb-0 mt-o pb-0 pt-0"><AboutApplication /></div>
      </div>
    </section>
  );
};

export default Layout;
