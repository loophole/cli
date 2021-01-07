import React, { useState } from "react";

import HostnameSettings from "../components/form/HostnameSettings";
import BasicAuthSettings from "../components/form/BasicAuthSettings";
import LocalWebserverSettings from "../components/form/LocalWebserverSettings";
import { useDispatch } from "react-redux";

import { send } from "@giantmachines/redux-websocket";

import ExposeHttpPort from "../interfaces/ExposeHttpPort";
import { useHistory } from "react-router-dom";
import Message from "../interfaces/Message";

const HTTP = () => {
  const dispatch = useDispatch();
  const history = useHistory();
  const [port, setPort] = useState(8080);
  const [hostname, setHostname] = useState("127.0.0.1");
  const [usingCustomHostname, setUsingCustomHostname] = useState(false);
  const [customHostname, setCustomHostname] = useState("");
  const [usingBasicAuth, setUsingBasicAuth] = useState(false);
  const [basicAuthUsername, setBasicAuthUsername] = useState("");
  const [basicAuthPassword, setBasicAuthPassword] = useState("");
  const [disableProxyErrorPage, setDisableProxyErrorPage] = useState(false);

  const startTunnel = () => {
    const options: ExposeHttpPort = {
      local: {
        host: hostname,
        port,
        https: false,
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

    options.display.disableProxyErrorPage = disableProxyErrorPage;

    const message: Message = {
      messageType: "MT_StartTunnel",
      startTunnelMessage: {
        tunnelType: "Tunnel_HTTP",
        exposeHttpConfig: options,
      },
    };

    dispatch(send(message));
    history.push("/dashboard");
  };

  return (
    <div className="container">
      <h4 className="subtitle is-4">
        Exposes http server running locally, or on locally available machine to
        the public via loophole tunnel.
      </h4>
      <div className="columns">
        <div className="column">
          <div className="box">
            <h5 className="title is-5">Local endpoint settings</h5>
            <LocalWebserverSettings
              hostnameValue={hostname}
              hostnameChangeCallback={setHostname}
              portValue={port}
              portChangeCallback={setPort}
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
            <div className="field">
              <div className="control">
                <label className="checkbox">
                  <input
                    type="checkbox"
                    onChange={(e) => {
                      setDisableProxyErrorPage(!disableProxyErrorPage);
                    }}
                  />{" "}
                  I want to disable proxy error page and use regular 502 error.
                </label>
              </div>
            </div>
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

export default HTTP;
