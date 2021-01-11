import React from "react";

import { Link, useLocation } from "react-router-dom";
import classNames from "classnames";

import { useAuth } from "../../features/auth/authProvider";

const Nav = () => {
  const location = useLocation();
  const auth = useAuth();
  return (
    <div className="tabs is-centered">
      <ul>
        {/* <li
          className={classNames({
            "is-active": location.pathname === "/",
          })}
        >
          <Link to="/">
            <span className="icon">
              <i className="far fa-question-circle"></i>
            </span>
            Help
          </Link>
        </li> */}
        {auth.loggedIn ? (
          <li
            className={classNames({
              "is-active": location.pathname === "/http",
            })}
          >
            <Link to="/http">
              <span className="icon">
                <i className="fas fa-server"></i>
              </span>
              HTTP
            </Link>
          </li>
        ) : null}
        {auth.loggedIn ? (
          <li
            className={classNames({
              "is-active": location.pathname === "/directory",
            })}
          >
            <Link to="/directory">
              <span className="icon">
                <i className="fas fa-folder"></i>
              </span>
              Directory
            </Link>
          </li>
        ) : null}
        {auth.loggedIn ? (
          <li
            className={classNames({
              "is-active": location.pathname === "/webdav",
            })}
          >
            <Link to="/webdav">
              <span className="icon">
                <i className="fas fa-folder-open"></i>
              </span>
              WebDav
            </Link>
          </li>
        ) : null}
        {auth.loggedIn ? (
          <li
            className={classNames({
              "is-active": location.pathname === "/dashboard",
            })}
          >
            <Link to="/dashboard">
              <span className="icon">
                <i className="fas fa-chart-line"></i>
              </span>
              Dashboard
            </Link>
          </li>
        ) : null}
        {!auth.loggedIn ? (
          <li
            className={classNames({
              "is-active": location.pathname === "/login",
            })}
          >
            <Link to="/login">
              <span className="icon">
                <i className="fas fa-sign-in-alt"></i>
              </span>
              Login
            </Link>
          </li>
        ) : null}
        <li
          className={classNames({
            "is-active": location.pathname === "/logs",
          })}
        >
          <Link to="/logs">
            <span className="icon">
              <i className="fas fa-align-justify"></i>
            </span>
            Logs
          </Link>
        </li>{" "}
        {auth.loggedIn ? (
          <li
            className={classNames({
              "is-active": location.pathname === "/logout",
            })}
          >
            <Link to="/logout">
              <span className="icon">
                <i className="fas fa-sign-out-alt"></i>
              </span>
              Logout
            </Link>
          </li>
        ) : null}
      </ul>
    </div>
  );
};

export default Nav;
