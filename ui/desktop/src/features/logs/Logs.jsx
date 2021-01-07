import React from "react";
import classNames from 'classnames';

const Logs = (props) => {
  const logs = props.logs.map((log, idx) => {
    return (
      <tr key={idx}>
        <td className="has-text-centered" width="5%">
          <span className={`icon has-text-${log.class}`}>
            <i className={classNames({
                fas: true,
                'fa-info': log.class === 'info',
                'fa-check': log.class === 'success',
                'fa-exclamation-triangle': log.class === 'warning',
                'fa-exclamation': log.class === 'danger',
                'fa-bell': log.class === '',
            })} />
          </span>
        </td>
        <td>{log.message}</td>
      </tr>
    );
  });

  return (
    <div className="card">
      <header className="card-header">
        <p className="card-header-title">Logs</p>
      </header>
      <div className={`card-content`}>
        {logs.length ? (
          <div className="content">
            <div className="table-container">
              <table className="table is-fullwidth is-striped is-narrow">
                <tbody>{logs}</tbody>
              </table>
            </div>
          </div>
        ) : (
          <div className="content has-text-centered">
            <span>There is no logs to display</span>
          </div>
        )}
      </div>
    </div>
  );
};

export default Logs;
