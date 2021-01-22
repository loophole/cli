import { useSelector, useDispatch } from "react-redux";
import { useHistory, useLocation } from "react-router-dom";

import { send } from "@giantmachines/redux-websocket";

import { useAuth } from "../features/config/authProvider";
import { MessageTypeRequestLogout } from "../constants/websocket";

const LoginPage = () => {
  const history = useHistory();
  const location = useLocation();
  const appState = useSelector((state) => state.config);
  const auth = useAuth();
  const dispatch = useDispatch();

  const logout = () => {
    const message = {
      type: MessageTypeRequestLogout,
    };

    dispatch(send(message));
  };

  const goBack = () => {
    const { from } = location.state || { from: { pathname: "/" } };
    history.replace(from);
  };

  if (!appState.loggedIn) {
    auth.logout(() => {
      history.replace("/account/login");
    });
    return null;
  }

  return (
    <div className="container">
      <h1 className="subtitle is-4">Logout</h1>
      <hr />
      <div className="context-box">
        <div className="content has-text-centered">
          <h2 className="subtitle is-4">Are you sure you want to do this?</h2>
        </div>
        <div className="buttons is-centered">
          <button
            className="button is-danger"
            onClick={logout}
            onKeyDown={logout}
          >
            Yes, log me out
          </button>
          <button className="button" onClick={goBack} onKeyDown={goBack}>
            No, take me back
          </button>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
