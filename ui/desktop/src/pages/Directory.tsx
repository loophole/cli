import React, { FormEvent, useState } from "react";

import HostnameSettings from "../components/form/HostnameSettings";
import BasicAuthSettings from "../components/form/BasicAuthSettings";
import LocalDirectorySettings from "../components/form/LocalDirectorySettings";
import { useDispatch, useSelector } from "react-redux";
import { v4 as uuidv4 } from "uuid";

import { send } from "@giantmachines/redux-websocket";

import ExposeDirectory from "../interfaces/ExposeDirectoryMessage";
import { useHistory } from "react-router-dom";
import Message from "../interfaces/Message";
import ExposeDirectoryMessage from "../interfaces/ExposeDirectoryMessage";
import { MessageTypeRequestTunnelStartDirectory } from "../constants/websocket";
import {
  isBasicAuthPasswordValid,
  isBasicAuthUsernameValid,
  isLocalPathValid,
  isLoopholeHostnameValid,
} from "../features/validator/validators";

const DirectoryPage = () => {
  const dispatch = useDispatch();
  const history = useHistory();
  const appState = useSelector((state: any) => state.config);
  const [path, setPath] = useState(appState.homeDirectory);
  const [usingCustomHostname, setUsingCustomHostname] = useState(false);
  const [customHostname, setCustomHostname] = useState("");
  const [usingBasicAuth, setUsingBasicAuth] = useState(false);
  const [basicAuthUsername, setBasicAuthUsername] = useState("");
  const [basicAuthPassword, setBasicAuthPassword] = useState("");
  const [disableOldCiphers, setDisableOldCiphers] = useState(false);

  const areInputsValid = (): boolean => {
    if (!isLocalPathValid(path)) return false;
    if (usingCustomHostname && !isLoopholeHostnameValid(customHostname))
      return false;
    if (
      usingBasicAuth &&
      (!isBasicAuthUsernameValid(basicAuthUsername) ||
        !isBasicAuthPasswordValid(basicAuthPassword))
    )
      return false;
    return true;
  };
  const startTunnel = (e: FormEvent) => {
    e.preventDefault();
    const options: ExposeDirectory = {
      local: {
        path,
      },
      remote: {
        disableProxyErrorPage: false,
        disableOldCiphers: false,
        tunnelId: uuidv4(),
      },
    };
    if (usingCustomHostname) {
      options.remote.siteId = customHostname;
    }
    if (usingBasicAuth) {
      options.remote.basicAuthUsername = basicAuthUsername;
      options.remote.basicAuthPassword = basicAuthPassword;
    }

    options.remote.disableOldCiphers = disableOldCiphers;

    const message: Message<ExposeDirectoryMessage> = {
      type: MessageTypeRequestTunnelStartDirectory,
      payload: options,
    };

    dispatch(send(message));
    history.push("/tunnels");
  };

  return (
    <div className="container">
      <h1 className="subtitle is-4">
        Exposes local directory to the public via loophole tunnel (download only
        mode) available through HTTPS.
      </h1>
      <hr />
      <div className="context-box">
        <form onSubmit={startTunnel}>
          <div className="columns is-multiline">
            <div className="column is-12">
              <h5 className="title is-5">Local directory settings</h5>
              <LocalDirectorySettings
                pathValue={path}
                pathChangeCallback={setPath}
              />
            </div>
            <div className="column is-12">
              <h5 className="title is-5">Remote endpoint settings</h5>
              <HostnameSettings
                usingValue={usingCustomHostname}
                usingChangeCallback={setUsingCustomHostname}
                hostnameValue={customHostname}
                hostnameChangeCallback={setCustomHostname}
              />
              <BasicAuthSettings
                usingValue={usingBasicAuth}
                usingChangeCallback={setUsingBasicAuth}
                usernameValue={basicAuthUsername}
                usernameChangeCallback={setBasicAuthUsername}
                passwordValue={basicAuthPassword}
                passwordChangeCallback={setBasicAuthPassword}
              />
              <div className="field">
                <div className="control">
                  <label className="checkbox">
                    <input
                      type="checkbox"
                      onChange={(e) => {
                        setDisableOldCiphers(!disableOldCiphers);
                      }}
                    />{" "}
                    I want to disable old TLS ciphers (older than TLS 1.2)
                  </label>
                </div>
              </div>
            </div>
            <div className="column is-12">
              <div className="field is-grouped is-pulled-right">
                <div className="control">
                  <button
                    className="button is-link"
                    disabled={!areInputsValid()}
                  >
                    Submit
                  </button>
                </div>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default DirectoryPage;
