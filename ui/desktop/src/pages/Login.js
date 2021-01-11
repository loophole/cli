import { useSelector, useDispatch } from "react-redux";
import { useHistory, useLocation } from "react-router-dom";

import QRCode from "qrcode.react";
import { send } from "@giantmachines/redux-websocket";

import { useAuth } from "../features/auth/authProvider";

import ExternalLink from "../components/routing/ExternalLink";

const LoginPage = () => {
  const history = useHistory();
  const location = useLocation();
  const authState = useSelector((state) => state.auth);
  const auth = useAuth();
  const dispatch = useDispatch();

  const startAuthProcess = () => {
    const message = {
      messageType: "MT_Authorize",
    };

    dispatch(send(message));
  };

  if (authState.loggedIn) {
    const { from } = location.state || { from: { pathname: "/" } };
    auth.login(() => {
      history.replace(from);
    });
    return null;
  }

  if (!authState.authInstructions) {
    if (authState.syncedWithBackend) startAuthProcess();

    return (
      <div className="container has-text-centered">
        <h1 class="title">Login</h1>
        <p>Obtaining authentication instructions...</p>
      </div>
    );
  }

  return (
    <div className="container">
      <div className="content has-text-centered">
        <h1 className="title">Login</h1>
      </div>
      <div className="tile is-ancestor">
        <div className="tile is-6 is-parent">
          <div className="tile is-child box has-text-centered">
            <p className="title is-3">QR Code</p>
            <p className="content">
              <QRCode
                renderAs="svg"
                size={512}
                width="80%"
                level="H"
                value={authState.authInstructions.verificationUriComplete}
              />
            </p>
          </div>
        </div>
        <div className="tile is-vertical is-parent">
          <div className="tile is-child box">
            <p className="title is-3">Direct link</p>
            <p className="content">
              <pre>
                <ExternalLink
                  url={authState.authInstructions.verificationUriComplete}
                />
              </pre>
            </p>
          </div>
          <div className="tile is-child box">
            <p className="title is-3">Manual steps</p>
            <p className="content">
              <ol>
                <li>
                  Visit
                  <ExternalLink
                    url={authState.authInstructions.verificationUri}
                  />
                </li>
                <li>
                  Provide
                  <code>{authState.authInstructions.userCode}</code>
                  as a code.
                </li>
              </ol>
            </p>
          </div>
          <div className="tile is-child box">
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
