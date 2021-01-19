import { useSelector, useDispatch } from "react-redux";
import { useHistory, useLocation } from "react-router-dom";

import QRCode from "qrcode.react";
import { send } from "@giantmachines/redux-websocket";

import { useAuth } from "../features/config/authProvider";

import OpenCopyBlock from "../components/routing/OpenCopyBlock";
import { MessageTypeRequestLogin } from "../constants/websocket";

const LoginPage = () => {
  const history = useHistory();
  const location = useLocation();
  const appState = useSelector((state) => state.config);
  const auth = useAuth();
  const dispatch = useDispatch();

  const startAuthProcess = () => {
    const message = {
      type: MessageTypeRequestLogin,
    };

    dispatch(send(message));
  };

  if (appState.loggedIn) {
    const { from } = location.state || { from: { pathname: "/" } };
    auth.login(() => {
      history.replace(from);
    });
    return null;
  }

  if (!appState.authInstructions) {
    if (appState.syncedWithBackend) startAuthProcess();

    return (
      <div className="container has-text-centered">
        <h1 className="title">Login</h1>
        <p>Obtaining authentication instructions...</p>
      </div>
    );
  }

  return (
    <div className="container">
      <h1 className="subtitle is-4">Login</h1>
      <hr />
      <div className="context-box">
        <div className="columns is-multiline">
          <div className="column is-12">
            <p className="title is-3">QR Code</p>
            <p className="content has-text-centered">
              <QRCode
                renderAs="svg"
                size={256}
                width="80%"
                level="H"
                value={appState.authInstructions.verificationUriComplete}
              />
            </p>
          </div>
          <div className="column is-12">
            <p className="title is-3">Direct link</p>
            <p className="content">
              <pre>
                <OpenCopyBlock
                  target={appState.authInstructions.verificationUriComplete}
                />
              </pre>
            </p>
          </div>
          <div className="column is-12">
            <p className="title is-3">Manual steps</p>
            <p className="content">
              <ol>
                <li>
                  Open below URL:
                  <OpenCopyBlock
                    target={appState.authInstructions.verificationUri}
                  />
                </li>
                <li>
                  Provide the following code:
                  <OpenCopyBlock
                    target={appState.authInstructions.userCode}
                    open={false}
                  ></OpenCopyBlock>
                </li>
              </ol>
            </p>
          </div>
          <div className="column is-12">
            <p className="title is-3">Restart the process</p>
            <p className="content">
              Sometimes it is required to restart the login process (e.g. when
              the code expires). Please use the below button is such case.
            </p>
            <button
              className="button is-danger"
              onClick={startAuthProcess}
              onKeyDown={startAuthProcess}
            >
              Restart login process
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
