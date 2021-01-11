import React from "react";
import Logs from "../features/logs/Logs";
import TunnelInfo from "../features/tunnels/TunnelInfo";

import StopTunnel from "../features/tunnels/StopTunnel";
import { useSelector } from "react-redux";

const DashboardPage = (): JSX.Element => {
  const tunnelsState = useSelector((store: any) => store.tunnels);
  const logsState = useSelector((store: any) => store.logs);

  if (!tunnelsState.tunnel || !tunnelsState.tunnel.siteId) {
    return (
      <div className="container has-text-centered">
        <p>Tunnel is not running!</p>
      </div>
    );
  }
  return (
    <div className="container">
      <div className="tile is-ancestor">
        <div className="tile is-4 is-parent">
          <div className="tile is-child">
            <TunnelInfo tunnel={tunnelsState.tunnel} />
          </div>
        </div>
        <div className="tile is-vertical is-parent">
          <div className="tile is-child">
            <Logs logs={logsState.tunnelLogs} />
          </div>
          <div className="tile is-child">
            <StopTunnel tunnel={tunnelsState.tunnel} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
