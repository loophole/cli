import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from "react-router-dom";

import PrivateRoute from "./components/routing/PrivateRoute";
import AuthProvider from "./features/config/authProvider";

import HTTPPage from "./pages/HTTP";
import DirectoryPage from "./pages/Directory";
import WebDavPage from "./pages/WebDav";
import DashboardPage from "./pages/Dashboard";
import LoginPage from "./pages/Login";
import ProfilePage from "./pages/Profile";
import LogoutPage from "./pages/Logout";
import LogsPage from "./pages/Logs";
import FeedbackPage from "./pages/Feedback";

import Layout from "./components/layout/Layout";
import TunnelsPage from "./pages/Tunnels";

const App = () => {
  return (
    <AuthProvider>
      <Router>
        <Layout title="Loophole" description="Instant hosting, right from your local machine">
          <Switch>
            <Route exact path="/">
              <Redirect to="/tunnels" />
            </Route>
            <PrivateRoute exact path="/tunnels">
              <TunnelsPage />
            </PrivateRoute>
            <PrivateRoute exact path="/tunnels/create">
              <Redirect to="/tunnels/create/http" />
            </PrivateRoute>
            <PrivateRoute path="/tunnels/create/http">
              <HTTPPage />
            </PrivateRoute>
            <PrivateRoute path="/tunnels/create/directory">
              <DirectoryPage />
            </PrivateRoute>
            <PrivateRoute path="/tunnels/create/webdav">
              <WebDavPage />
            </PrivateRoute>
            <PrivateRoute path="/tunnels/:tunnelId">
              {/* Nested routes inside */}
              <DashboardPage />
            </PrivateRoute>
            <Route path="/application/logs">
              <LogsPage />
            </Route>
            <Route path="/application/feedback">
              <FeedbackPage />
            </Route>
            <Route path="/account/login">
              <LoginPage />
            </Route>
            <Route path="/account/profile">
              <ProfilePage />
            </Route>
            <PrivateRoute path="/account/logout">
              <LogoutPage />
            </PrivateRoute>
          </Switch>
        </Layout>
      </Router>
    </AuthProvider>
  );
};

export default App;
