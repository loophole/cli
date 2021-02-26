import React from "react";
import classNames from "classnames";
import {
  isLocalHostValid,
  isLocalPortValid,
  isUrlPathValid,
} from "../../features/validator/validators";

interface LocalWebserverSettingsProps {
  hostnameValue: string;
  hostnameChangeCallback: Function;
  portValue: string;
  portChangeCallback: Function;
  httpsChangeCallback: Function;
  httpsValue: boolean;
  usingPathValue: boolean;
  usingPathChangeCallback: Function;
  urlPathValue: string;
  urlPathChangeCallback: Function;
}

const LocalWebserverSettings = (
  props: LocalWebserverSettingsProps
): JSX.Element => {
  const setHostname = props.hostnameChangeCallback;
  const setPort = props.portChangeCallback;
  const port = parseInt(props.portValue, 10);
  const setHTTPS = props.httpsChangeCallback;
  const setUsingPath = props.usingPathChangeCallback;
  const setUrlPath = props.urlPathChangeCallback;

  const isHostValid = (): boolean => {
    return isLocalHostValid(props.hostnameValue);
  };

  const isPortValid = (): boolean => {
    return isLocalPortValid(port);
  };
  const isPathValid = (): boolean => {
    return isUrlPathValid(props.urlPathValue);
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
            value={props.hostnameValue}
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
              })}
            ></i>
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
              })}
            ></i>
          </span>
          {isPortValid() ? (
            <p className="help is-success">Port is valid</p>
          ) : (
            <p className="help is-danger">Port is invalid</p>
          )}
        </div>
      </div>
      <div className="field">
        <div className="control">
          <label className="checkbox">
            <input
              type="checkbox"
              onChange={(e) => {
                setUsingPath(!props.usingPathValue);
              }}
            />{" "}
            I want to serve subpath only.
          </label>
        </div>
      </div>
      {props.usingPathValue ? (
        <div className="field">
          <label className="label">Subpath</label>
          <div className="control has-icons-left has-icons-right">
            <input
              className={classNames({
                input: true,
                "is-success": isPathValid(),
                "is-danger": !isPathValid(),
              })}
              type="text"
              placeholder="Subpath which you want to expose"
              value={props.urlPathValue}
              onChange={(e) => setUrlPath(e.target.value)}
            />
            <span className="icon is-small is-left">
              <i className="fas fa-signature"></i>
            </span>
            <span className="icon is-small is-right">
              <i
                className={classNames({
                  fas: true,
                  "fa-check": isPathValid(),
                  "fa-exclamation-triangle": !isPathValid(),
                })}
              ></i>
            </span>
          </div>
          {isPathValid() ? (
            <p className="help is-success">Path is valid</p>
          ) : (
            <p className="help is-danger">Path is invalid</p>
          )}
        </div>
      ) : null}
      <div className="field">
        <div className="control">
          <label className="checkbox">
            <input
              type="checkbox"
              onChange={(e) => {
                setHTTPS(!props.httpsValue);
              }}
            />{" "}
            The server is already using HTTPS.
          </label>
        </div>
      </div>
    </div>
  );
};

export default LocalWebserverSettings;
