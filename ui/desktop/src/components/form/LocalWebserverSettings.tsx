import React from "react";
import classNames from "classnames";

interface LocalWebserverSettingsProps {
  hostnameValue: string;
  hostnameChangeCallback: Function;
  portValue: number;
  portChangeCallback: Function;
}
const LocalWebserverSettings = (
  props: LocalWebserverSettingsProps
): JSX.Element => {
  const hostname = props.hostnameValue;
  const setHostname = props.hostnameChangeCallback;
  const port = props.portValue;
  const setPort = props.portChangeCallback;

  const isHostValid = (): boolean => {
    return hostname.length >= 1;
  };

  const isPortValid = (): boolean => {
    return port > 0;
  };

  return (
    <div>
      <div className="field">
        <label className="label">Host</label>
        <div className="control has-icons-left has-icons-right">
          <input
            className={classNames({
              input: true,
              "is-success": isHostValid(),
              "is-danger": !isHostValid(),
            })}
            type="text"
            placeholder="Host on which the server is running"
            value={hostname}
            onChange={(e) => setHostname(e.target.value)}
          />
          <span className="icon is-small is-left">
            <i className="fas fa-signature"></i>
          </span>
          <span className="icon is-small is-right">
            <i 
                    className={classNames({
                      fas: true,
                      "fa-check": isHostValid(),
                      "fa-exclamation-triangle": !isHostValid(),
                    })}></i>
          </span>
        </div>
              {isHostValid() ? (
                <p className="help is-success">Host is valid</p>
              ) : (
                <p className="help is-danger">Host is invalid</p>
              )}
      </div>
      <div className="field">
        <label className="label">Port</label>
        <div className="control has-icons-left has-icons-right">
          <input
            className={classNames({
              input: true,
              "is-success": isPortValid(),
              "is-danger": !isPortValid(),
            })}
            type="number"
            placeholder="Port on which the server is running"
            value={port}
            onChange={(e) => setPort(e.target.value)}
          />
          <span className="icon is-small is-left">
            <i className="fas fa-plug"></i>
          </span>
          <span className="icon is-small is-right">
            <i 
                    className={classNames({
                      fas: true,
                      "fa-check": isPortValid(),
                      "fa-exclamation-triangle": !isPortValid(),
                    })}></i>
          </span>
              {isPortValid() ? (
                <p className="help is-success">Port is valid</p>
              ) : (
                <p className="help is-danger">Port is invalid</p>
              )}
        </div>
      </div>
    </div>
  );
};

export default LocalWebserverSettings;
