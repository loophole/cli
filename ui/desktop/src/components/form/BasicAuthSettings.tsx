import React from "react";
import classNames from "classnames";

interface BasicAuthSettingsProps {
  usingValue: boolean;
  usingChangeCallback: Function;
  usernameValue: string;
  usernameChangeCallback: Function;
  passwordValue: string;
  passwordChangeCallback: Function;
}

const BasicAuthSettings = (props: BasicAuthSettingsProps): JSX.Element => {
  const usingBasicAuth = props.usingValue;
  const setUsingBasicAuth = props.usingChangeCallback;
  const basicAuthUsername = props.usernameValue;
  const setBasicAuthUsername = props.usernameChangeCallback;
  const basicAuthPassword = props.passwordValue;
  const setBasicAuthPassword = props.passwordChangeCallback;

  const isUsernameValid = (): boolean => {
    return basicAuthUsername.length >= 3;
  };
  const isPasswordValid = (): boolean => {
    return basicAuthPassword.length >= 3;
  };

  return (
    <div>
      <div className="field">
        <div className="control">
          <label className="checkbox">
            <input
              type="checkbox"
              onChange={(e) => {
                setUsingBasicAuth(!usingBasicAuth);
              }}
            />{" "}
            I want to use basic auth
          </label>
        </div>
      </div>
      {usingBasicAuth
        ? [
            <div className="field" key="username">
              <label className="label">Username</label>
              <div className="control has-icons-left has-icons-right">
                <input
                  className={classNames({
                    input: true,
                    "is-success": isUsernameValid(),
                    "is-danger": !isUsernameValid(),
                  })}
                  type="text"
                  placeholder="Basic auth username"
                  value={basicAuthUsername}
                  onChange={(e) => setBasicAuthUsername(e.target.value)}
                />
                <span className="icon is-small is-left">
                  <i className="fas fa-user"></i>
                </span>
                <span className="icon is-small is-right">
                  <i
                    className={classNames({
                      fas: true,
                      "fa-check": isUsernameValid(),
                      "fa-exclamation-triangle": !isUsernameValid(),
                    })}
                  ></i>
                </span>
              </div>
              {isUsernameValid() ? (
                <p className="help is-success">Username is valid</p>
              ) : (
                <p className="help is-danger">Username is invalid</p>
              )}
            </div>,
            <div className="field" key="password">
              <label className="label">Password</label>
              <div className="control has-icons-left has-icons-right">
                <input
                  className={classNames({
                    input: true,
                    "is-success": isPasswordValid(),
                    "is-danger": !isPasswordValid(),
                  })}
                  type="password"
                  placeholder="Basic auth password"
                  value={basicAuthPassword}
                  onChange={(e) => setBasicAuthPassword(e.target.value)}
                />
                <span className="icon is-small is-left">
                  <i className="fas fa-key"></i>
                </span>
                <span className="icon is-small is-right">
                  <i
                    className={classNames({
                      fas: true,
                      "fa-check": isPasswordValid(),
                      "fa-exclamation-triangle": !isPasswordValid(),
                    })}
                  ></i>
                </span>
              </div>
              {isPasswordValid() ? (
                <p className="help is-success">Password is valid</p>
              ) : (
                <p className="help is-danger">Password is invalid</p>
              )}
            </div>,
          ]
        : []}
    </div>
  );
};

export default BasicAuthSettings;
