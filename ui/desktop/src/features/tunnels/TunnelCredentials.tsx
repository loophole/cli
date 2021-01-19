import classNames from "classnames";
import React, { useState } from "react";
import CopyToClipboard from "react-copy-to-clipboard";
import Tunnel from "./interfaces/Tunnel";

const TunnelCredentials = (props: { tunnel: Tunnel }) => {
  const [usernameCopied, setUsernameCopied] = useState(false);
  const [passwordCopied, setPasswordCopied] = useState(false);
  const [passwordVisible, setPasswordVisible] = useState(false);

  const onUsernameCopiedEffect = () => {
    setUsernameCopied(true);
    setTimeout(() => {
      setUsernameCopied(false);
    }, 1000);
  };
  const onPasswordCopiedEffect = () => {
    setPasswordCopied(true);
    setTimeout(() => {
      setPasswordCopied(false);
    }, 1000);
  };

  return (
    <div className="content">
      <p>Your basic auth credentials for this tunnel are shown below</p>
      <div className="field is-grouped">
        <p className="control has-icons-left has-icons-right is-expanded">
          <input
            className="input"
            type="text"
            placeholder="Username"
            value={props.tunnel.basicAuthUsername || ""}
            readOnly={true}
          />
          <span className="icon is-small is-left">
            <i className="fas fa-user" />
          </span>
        </p>
        <p className="control">
          <CopyToClipboard
            text={props.tunnel.basicAuthUsername || ""}
            onCopy={onUsernameCopiedEffect}
          >
            <button className="button is-light">
              {usernameCopied ? "Copied!" : "Copy"}
            </button>
          </CopyToClipboard>
        </p>
      </div>
      <div className="field is-grouped">
        <p className="control has-icons-left has-icons-right is-expanded">
          <input
            className="input"
            type={passwordVisible ? "text" : "password"}
            placeholder="Password"
            value={props.tunnel.basicAuthPassword || ""}
            readOnly={true}
          />
          <span className="icon is-small is-left">
            <i className="fas fa-lock" />
          </span>
          <span
            className="icon is-small is-right is-clickable"
            style={{
              pointerEvents: "all",
            }}
            onClick={() => setPasswordVisible(!passwordVisible)}
          >
            <i
              className={classNames({
                far: true,
                "fa-eye": !passwordVisible,
                "fa-eye-slash": passwordVisible,
              })}
            />
          </span>
        </p>
        <p className="control">
          <CopyToClipboard
            text={props.tunnel.basicAuthPassword || ""}
            onCopy={onPasswordCopiedEffect}
          >
            <button className="button is-light">
              {passwordCopied ? "Copied!" : "Copy"}
            </button>
          </CopyToClipboard>
        </p>
      </div>
    </div>
  );
};

export default TunnelCredentials;
