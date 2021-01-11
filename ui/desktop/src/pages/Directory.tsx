import React, { useState } from "react";

import HostnameSettings from "../components/form/HostnameSettings";
import BasicAuthSettings from "../components/form/BasicAuthSettings";
import LocalDirectorySettings from "../components/form/LocalDirectorySettings";
import { useDispatch } from "react-redux";

import { send } from "@giantmachines/redux-websocket";

import ExposeDirectory from "../interfaces/ExposeDirectory";
import { useHistory } from "react-router-dom";
import Message from "../interfaces/Message";

const DirectoryPage = () => {
  const dispatch = useDispatch();
  const history = useHistory();
  const [path, setPath] = useState("");
  const [usingCustomHostname, setUsingCustomHostname] = useState(false);
  const [customHostname, setCustomHostname] = useState("");
  const [usingBasicAuth, setUsingBasicAuth] = useState(false);
  const [basicAuthUsername, setBasicAuthUsername] = useState("");
  const [basicAuthPassword, setBasicAuthPassword] = useState("");

  const startTunnel = () => {
    const options: ExposeDirectory = {
      local: {
        path,
      },
      remote: {},
      display: {},
    };
    if (usingCustomHostname) {
      options.remote.siteID = customHostname;
    }
    if (usingBasicAuth) {
      options.remote.basicAuthUsername = basicAuthUsername;
      options.remote.basicAuthPassword = basicAuthPassword;
    }

    const message: Message = {
      messageType: "MT_StartTunnel",
      startTunnelMessage: {
        tunnelType: "Tunnel_Directory",
        exposeDirectoryConfig: options,
      },
    };

    dispatch(send(message));
    history.push("/dashboard");
  };

  return (
    <div className="container">
      <h4 className="subtitle is-4">
        Exposes local directory to the public via loophole tunnel (download only mode)
        available through HTTPS.
      </h4>
      <div className="columns">
        <div className="column">
          <div className="box">
            <h5 className="title is-5">Local directory settings</h5>
            <LocalDirectorySettings
              pathValue={path}
              pathChangeCallback={setPath}
            />
          </div>
        </div>
        <div className="column">
          <div className="box">
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
          </div>
        </div>
      </div>
      <div>
        <div className="field is-grouped is-pulled-right">
          <div className="control">
            <button
              className="button is-link"
              onClick={() => startTunnel()}
              onKeyDown={() => startTunnel()}
            >
              Submit
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DirectoryPage;
