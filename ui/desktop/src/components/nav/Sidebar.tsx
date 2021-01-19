import React from "react";
import { Link, useLocation } from "react-router-dom";
import classNames from "classnames";
import { useAuth } from "../../features/config/authProvider";
import { useSelector } from "react-redux";

const Sidebar = () => {
  const location = useLocation();
  const auth = useAuth();
  const appState = useSelector((store: any) => store.config);

  let userLoggedInAs = "";
  if (appState.user) {
    if (appState.user.nickname) {
      userLoggedInAs = appState.user.nickname;
    } else if (appState.user.name) {
      userLoggedInAs = appState.user.name;
    } else if (appState.user.email) {
      userLoggedInAs = appState.user.email;
    }
  }

  return (
    <aside className="menu is-narrow-mobile is-fullheight">
      <figure className="image">
        <Link to="/">
          <img src="/logo.png" alt="Loophole" />
        </Link>
      </figure>
      {auth.loggedIn ? <p className="menu-label">Tunnels</p> : null}
      {auth.loggedIn ? (
        <ul className="menu-list">
          <li>
            <Link
              to="/tunnels"
              className={classNames({
                "is-active": location.pathname === "/tunnels",
              })}
            >
              <span className="icon">
                <i className="fas fa-list"></i>
              </span>
              List
            </Link>
          </li>
          <li>
            <Link
              to="/tunnels/create"
              className={classNames({
                "is-active": location.pathname === "/tunnels/create",
              })}
            >
              <span className="icon">
                <i className="fas fa-plus"></i>
              </span>
              Create
            </Link>
            <ul>
              <li>
                <Link
                  to="/tunnels/create/http"
                  className={classNames({
                    "is-active": location.pathname === "/tunnels/create/http",
                  })}
                >
                  <span className="icon">
                    <i className="fas fa-server"></i>
                  </span>
                  HTTP
                </Link>
              </li>
              <li>
                <Link
                  to="/tunnels/create/directory"
                  className={classNames({
                    "is-active":
                      location.pathname === "/tunnels/create/directory",
                  })}
                >
                  <span className="icon">
                    <i className="fas fa-folder"></i>
                  </span>
                  Directory
                </Link>
              </li>
              <li>
                <Link
                  to="/tunnels/create/webdav"
                  className={classNames({
                    "is-active": location.pathname === "/tunnels/create/webdav",
                  })}
                >
                  <span className="icon">
                    <i className="fas fa-folder-open"></i>
                  </span>
                  WebDav
                </Link>
              </li>
            </ul>
          </li>
        </ul>
      ) : null}
      <p className="menu-label">Application</p>
      <ul className="menu-list">
        <li>
          <Link
            to="/application/logs"
            className={classNames({
              "is-active": location.pathname === "/application/logs",
            })}
          >
            <span className="icon">
              <i className="fas fa-clipboard-list"></i>
            </span>
            Logs
          </Link>
        </li>
        <li>
          <Link
            to="/application/feedback"
            className={classNames({
              "is-active": location.pathname === "/application/feedback",
            })}
          >
            <span className="icon">
              <i className="fas fa-comment-dots"></i>
            </span>
            Feedback
          </Link>
        </li>
      </ul>
      <p className="menu-label">Account</p>
      <ul className="menu-list">
        {!auth.loggedIn ? (
          <li>
            <Link
              to="/account/login"
              className={classNames({
                "is-active": location.pathname === "/account/login",
              })}
            >
              <span className="icon">
                <i className="fas fa-sign-in-alt"></i>
              </span>
              Login
            </Link>
          </li>
        ) : null}
        {auth.loggedIn && userLoggedInAs ? (
          <li>
            <Link
              to={{
                pathname: "/account/profile",
                state: { from: location },
              }}
              className={classNames({
                "is-active": location.pathname === "/account/profile",
              })}
            >
              <span className="icon">
                <i className="fas fa-user"></i>
              </span>
              <span className="is-uppercase">
                <small>{userLoggedInAs}</small>
              </span>
            </Link>
          </li>
        ) : null}
        {auth.loggedIn ? (
          <li>
            <Link
              to={{
                pathname: "/account/logout",
                state: { from: location },
              }}
              className={classNames({
                "is-active": location.pathname === "/account/logout",
              })}
            >
              <span className="icon">
                <i className="fas fa-sign-out-alt"></i>
              </span>
              Logout
            </Link>
          </li>
        ) : null}
      </ul>
    </aside>
  );
};

export default Sidebar;
