import React from "react";
import Logs from "../features/logs/Logs";
import TunnelInfo from "../features/tunnels/TunnelInfo";

import TunnelActions from "../features/tunnels/TunnelActions";
import TunnelCredentials from "../features/tunnels/TunnelCredentials";
import { useSelector } from "react-redux";
import {
  Link,
  Route,
  Switch,
  useLocation,
  useParams,
  useRouteMatch,
} from "react-router-dom";
import TunnelsState from "../features/tunnels/interfaces/TunnelsState";
import classNames from "classnames";
import TunnelStatusIcons from "../features/tunnels/TunnelStatusIcon";

interface DashboardParams {
  tunnelId: string;
}

const DashboardPage = (): JSX.Element => {
  const { path, url } = useRouteMatch();
  const location = useLocation();

  const params: DashboardParams = useParams();
  const tunnelsState: TunnelsState = useSelector((store: any) => store.tunnels);
  const logsState = useSelector((store: any) => store.logs);

  const tunnel = tunnelsState.tunnels.find(
    (tunnel) => tunnel.tunnelId === params.tunnelId
  );
  const tunnelLogs = logsState.tunnelLogs[params.tunnelId] || [];
  if (!tunnel) {
    return (
      <div className="container has-text-centered">
        <p>Tunnel is not running!</p>
      </div>
    );
  }
  return (
    <div className="container">
      <h1 className="subtitle is-4">Tunnel: {tunnel.siteId? tunnel.siteId : params.tunnelId}</h1>
      <hr />
      <div className="context-box">
        <div className="columns is-multiline">
          <div className="column is-12">
            <div className="tabs">
              <ul>
                <li
                  className={classNames({
                    "is-active": location.pathname === url,
                  })}
                >
                  <Link to={url}>
                    <span className="icon is-small">
                      <i className="fas fa-info" />
                    </span>
                    Summary
                  </Link>
                </li>
                <li
                  className={classNames({
                    "is-active": location.pathname === `${url}/logs`,
                  })}
                >
                  <Link to={`${url}/logs`}>
                    <span className="icon is-small">
                      <i className="fas fa-clipboard-list" />
                    </span>
                    Logs
                  </Link>
                </li>
                {tunnel.usingBasicAuth ? (
                  <li
                    className={classNames({
                      "is-active": location.pathname === `${url}/credentials`,
                    })}
                  >
                    <Link to={`${url}/credentials`}>
                      <span className="icon is-small">
                        <i className="fas fa-user" />
                      </span>
                      Credentials
                    </Link>
                  </li>
                ) : null}
              </ul>
            </div>
          </div>
          <div className="column is-12">
            <Switch>
              <Route exact path={path}>
                <TunnelStatusIcons tunnel={tunnel} />
                <TunnelInfo tunnel={tunnel} />
                <TunnelActions tunnel={tunnel} />
              </Route>
              <Route path={`${path}/logs`}>
                <Logs logs={tunnelLogs} />
              </Route>
              <Route path={`${path}/credentials`}>
                <TunnelCredentials tunnel={tunnel} />
              </Route>
            </Switch>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
