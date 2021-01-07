import { useSelector, useDispatch } from "react-redux";
import { useHistory, useLocation } from "react-router-dom";

import { send } from "@giantmachines/redux-websocket";

import { useAuth } from "../features/auth/authProvider";

const LoginPage = () => {
  const history = useHistory();
  const location = useLocation();
  const authState = useSelector((state) => state.auth);
  const auth = useAuth();
  const dispatch = useDispatch();

  const logout = () => {
    const message = {
      messageType: "MT_Logout",
    };

    dispatch(send(message));
  };

  const goBack = () => {
    console.log(location);
    const { from } = location.state || { from: { pathname: "/" } };
    history.replace(from);
  };

  if (!authState.loggedIn) {
    auth.logout(() => {
      history.replace("/login");
    });
    return null;
  }

  return (
    <div className="container">
      <div className="content has-text-centered">
        <h1 className="title">Logout</h1>

        <h2 className="subtitle">Are you sure you want to log out?</h2>
      </div>
      <div className="buttons is-centered">
        <button
          className="button is-danger"
          onClick={logout}
          onKeyDown={logout}
        >
          Yes, log me out
        </button>
        <button
          className="button is-success"
          onClick={goBack}
          onKeyDown={goBack}
        >
          No, take me back
        </button>
      </div>
    </div>
  );
};

export default LoginPage;
