import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Route, Switch } from "react-router-dom";

import "@fortawesome/fontawesome-free/css/all.min.css";
import "./assets/plugins/nucleo/css/nucleo.css";
import "./assets/scss/dashboard.scss";

import Layout from "./components/common/Layout";

ReactDOM.render(
  <BrowserRouter>
    <Switch>
      <Route path="/" render={props => <Layout {...props} />} />
    </Switch>
  </BrowserRouter>,
  document.getElementById("root")
);
