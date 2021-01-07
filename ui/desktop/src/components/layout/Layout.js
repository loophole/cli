import React from "react";
import { Link } from "react-router-dom";
import Nav from "../nav/Nav";

import { WebSocket } from "../../features/websocket/WebSocket";

const Layout = (props) => {
  return [
    <section className="hero is-small mt-6" key="logo">
      <div className="container">
        <div className="column is-4 is-offset-4">
          <figure className="image">
            <Link to="/">
              <img src="/logo.png" alt="Loophole" />
            </Link>
          </figure>
        </div>
      </div>
    </section>,
    <section className="section" key="nav">
      <Nav />
      <WebSocket />
    </section>,
    <section className="section" key="content">
      {props.children}
    </section>,
  ];
};

export default Layout;
