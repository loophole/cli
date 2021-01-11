import React from "react";
import Logs from "../features/logs/Logs";

import { useSelector } from "react-redux";

const LogsPage = (): JSX.Element => {
  const logsState = useSelector((store: any) => store.logs);

  return (
    <div className="container">
      <Logs logs={logsState.globalLogs}/>
    </div>
  );
};

export default LogsPage;
