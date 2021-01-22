import React from "react";
import classNames from "classnames";
import { isLoopholeHostnameValid } from "../../features/validator/validators";

interface BasicAuthSettingsProps {
  usingValue: boolean;
  usingChangeCallback: Function;
  hostnameValue: string;
  hostnameChangeCallback: Function;
}

const HostnameSettings = (props: BasicAuthSettingsProps): JSX.Element => {
  const usingCustomHostname = props.usingValue;
  const setUsingCustomHostname = props.usingChangeCallback;
  const customHostname = props.hostnameValue;
  const setCustomHostname = props.hostnameChangeCallback;

  const isHostnameValid = (): boolean => {
    return isLoopholeHostnameValid(customHostname);
  };

  return (
    <div>
      <div className="field">
        <div className="control">
          <label className="checkbox">
            <input
              type="checkbox"
              onChange={(e) => {
                setUsingCustomHostname(!usingCustomHostname);
              }}
            />{" "}
            I want to use custom hostname
          </label>
        </div>
      </div>
      {usingCustomHostname ? (
        <div className="field">
          <label className="label">Custom hostname</label>
          <div className="control has-icons-left has-icons-right">
            <input
              className={classNames({
                input: true,
                "is-success": isHostnameValid(),
                "is-danger": !isHostnameValid(),
              })}
              type="text"
              placeholder="Hostname to expose tunnel on"
              value={customHostname}
              onChange={(e) => setCustomHostname(e.target.value.toLowerCase())}
            />
            <span className="icon is-small is-left">
              <i className="fas fa-signature"></i>
            </span>
            <span className="icon is-small is-right">
              <i
                className={classNames({
                  fas: true,
                  "fa-check": isHostnameValid(),
                  "fa-exclamation-triangle": !isHostnameValid(),
                })}
              ></i>
            </span>
          </div>
          {isHostnameValid() ? (
            <p className="help is-success">Hostname is valid</p>
          ) : (
            <p className="help is-danger">Hostname is invalid</p>
          )}
        </div>
      ) : null}
    </div>
  );
};

export default HostnameSettings;
