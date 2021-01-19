import React, { FormEvent, useState } from "react";

import HostnameSettings from "../components/form/HostnameSettings";
import BasicAuthSettings from "../components/form/BasicAuthSettings";
import LocalWebserverSettings from "../components/form/LocalWebserverSettings";
import { useDispatch } from "react-redux";
import { v4 as uuidv4 } from "uuid";

import { send } from "@giantmachines/redux-websocket";

import ExposeHttpPort from "../interfaces/ExposeHttpPortMessage";
import { useHistory } from "react-router-dom";
import Message from "../interfaces/Message";
import ExposeHttpPortMessage from "../interfaces/ExposeHttpPortMessage";
import { MessageTypeRequestTunnelStartHTTP } from "../constants/websocket";

const HTTP = () => {
  const dispatch = useDispatch();
  const history = useHistory();
  const [port, setPort] = useState("8080");
  const [hostname, setHostname] = useState("127.0.0.1");
  const [isHTTPS, setIsHTTPS] = useState(false);
  const [usingCustomHostname, setUsingCustomHostname] = useState(false);
  const [customHostname, setCustomHostname] = useState("");
  const [usingBasicAuth, setUsingBasicAuth] = useState(false);
  const [basicAuthUsername, setBasicAuthUsername] = useState("");
  const [basicAuthPassword, setBasicAuthPassword] = useState("");
  const [disableProxyErrorPage, setDisableProxyErrorPage] = useState(false);


  const areInputsValid = (): boolean => {
    if (hostname.length === 0) return false;
    if (parseInt(port, 10) <= 0) return false;
    if (
      usingCustomHostname &&
      customHostname.match(/^[a-z][a-z0-9]{0,30}$/) === null
    )
      return false;
    if (
      usingBasicAuth &&
      (basicAuthUsername.length < 3 || basicAuthPassword.length < 3)
    )
      return false;
    return true;
  };

  const startTunnel = (e : FormEvent) => {
    e.preventDefault();
    const options: ExposeHttpPort = {
      local: {
        host: hostname,
        port: parseInt(port, 10),
        https: isHTTPS,
      },
      remote: {
        disableProxyErrorPage: false,
        tunnelId: uuidv4()
      },
    };
    if (usingCustomHostname) {
      options.remote.siteId = customHostname;
    }
    if (usingBasicAuth) {
      options.remote.basicAuthUsername = basicAuthUsername;
      options.remote.basicAuthPassword = basicAuthPassword;
    }

    options.remote.disableProxyErrorPage = disableProxyErrorPage;

    const message: Message<ExposeHttpPortMessage> = {
      type: MessageTypeRequestTunnelStartHTTP,
      payload: options,
    };

    dispatch(send(message));
    history.push("/tunnels");
  };

  return (
    <div className="container">
      <h1 className="subtitle is-4">
        Exposes http server running locally, or on locally available machine to
        the public via loophole tunnel.
      </h1>
      <hr />
      <div className="context-box">
        <form onSubmit={startTunnel}>
          <div className="columns is-multiline">
            <div className="column is-12">
              <h5 className="title is-5">Local endpoint settings</h5>
              <LocalWebserverSettings
                hostnameValue={hostname}
                hostnameChangeCallback={setHostname}
                portValue={port}
                portChangeCallback={setPort}
                httpsValue={isHTTPS}
                httpsChangeCallback={setIsHTTPS}
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
                        setDisableProxyErrorPage(!disableProxyErrorPage);
                      }}
                    />{" "}
                    I want to disable proxy error page and use regular 502 error
                  </label>
                </div>
              </div>
            </div>
            <div className="column is-12">
              <div className="field is-grouped is-pulled-right">
                <div className="control">
                  <button
                    type="submit"
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

export default HTTP;
