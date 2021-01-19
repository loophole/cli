import React from "react";
import classNames from "classnames";
import Log from "./interfaces/Log";

const Logs = (props: { logs: Log[] }) => {
  const logs = props.logs.map((log, idx) => {
    return (
      <tr key={idx}>
        <td className="has-text-centered" width="5%">
          <span className={`icon has-text-${log.class ? log.class : "grey"}`}>
            <i
              className={classNames({
                fas: true,
                "fa-info": log.class === "info",
                "fa-check": log.class === "success",
                "fa-exclamation-triangle": log.class === "warning",
                "fa-exclamation": log.class === "danger",
                "fa-bell": log.class === "" || log.class === "grey",
              })}
            />
          </span>
        </td>
        <td>{log.timestamp.toLocaleString()}</td>
        <td>{log.message}</td>
      </tr>
    );
  });

  return logs.length ? (
    <div className="content">
      <div className="table-container">
        <table className="table is-fullwidth is-striped is-narrow">
          <tbody>{logs}</tbody>
        </table>
      </div>
    </div>
  ) : (
    <div className="content has-text-centered">
      <span>There is no logs to be displayed</span>
    </div>
  );
};

export default Logs;
