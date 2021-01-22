import React from "react";
import Tunnel from "./interfaces/Tunnel";

const TunnelStatusIcons = (props: { tunnel: Tunnel }) => {
  return (
    <nav className="level is-mobile">
      <div className="level-left">
        {props.tunnel.type === "HTTP" ? (
          <span className="level-item">
            <span
              className="icon is-small has-tooltip-right has-tooltip-arrow"
              title={"It's a tunnel to local webserver"}
            >
              <i className="fas fa-server" />
            </span>
          </span>
        ) : null}
        {props.tunnel.type === "Directory" ? (
          <span className="level-item">
            <span
              className="icon is-small has-tooltip-right has-tooltip-arrow"
              title={"It's a HTTPS tunnel to local directory"}
            >
              <i className="fas fa-folder" />
            </span>
          </span>
        ) : null}
        {props.tunnel.type === "WebDav" ? (
          <span className="level-item">
            <span
              className="icon is-small has-tooltip-right has-tooltip-arrow"
              title={"It's a WebDav tunnel to local directory"}
            >
              <i className="fas fa-folder-open" />
            </span>
          </span>
        ) : null}
        {props.tunnel.usingBasicAuth ? (
          <span className="level-item">
            <span
              className="icon is-small has-tooltip-right has-tooltip-arrow"
              title={"Tunnel is protected by basic auth"}
            >
              <i className="fas fa-user-lock" />
            </span>
          </span>
        ) : null}
        {props.tunnel.proxyErrorDisabled ? (
          <span className="level-item">
            <span
              className="icon is-small has-tooltip-right has-tooltip-arrow"
              title={"Tunnel is not using customized 502 error page"}
            >
              <i className="fas fa-store-slash" />
            </span>
          </span>
        ) : null}
      </div>
    </nav>
  );
};

export default TunnelStatusIcons;
