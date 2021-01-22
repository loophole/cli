import React from "react";
import Logs from "../features/logs/Logs";

import { useSelector } from "react-redux";

const LogsPage = (): JSX.Element => {
  const logsState = useSelector((store: any) => store.logs);

  return (
    <div className="container">
      <h1 className="subtitle is-4">Logs</h1>
      <hr />
      <div className="context-box">
        <div className="columns is-multiline">
          <div className="column is-12">
            <Logs logs={logsState.communicationLogs} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default LogsPage;
