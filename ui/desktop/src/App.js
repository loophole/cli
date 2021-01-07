import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect
} from "react-router-dom";

import PrivateRoute from './components/routing/PrivateRoute';
import AuthProvider from './features/auth/authProvider';

import HTTPPage from "./pages/HTTP";
import DirectoryPage from "./pages/Directory";
import WebDavPage from "./pages/WebDav";
import DashboardPage from "./pages/Dashboard";
import LoginPage from './pages/Login';
import LogoutPage from './pages/Logout';
import LogsPage from './pages/Logs';

import Layout from "./components/layout/Layout";

const App = () => {
  return (
    <AuthProvider>
      <Router>
        <Layout title="Loophole" description="Test">
          <Switch>
            <PrivateRoute exact path="/">
              <Redirect to="/http" />
            </PrivateRoute>
            <PrivateRoute path="/http">
              <HTTPPage />
            </PrivateRoute>
            <PrivateRoute path="/directory">
              <DirectoryPage />
            </PrivateRoute>
            <PrivateRoute path="/webdav">
              <WebDavPage />
            </PrivateRoute>
            <PrivateRoute path="/dashboard">
              <DashboardPage />
            </PrivateRoute>
            <Route path="/login">
              <LoginPage />
            </Route>
            <Route path="/logout">
              <LogoutPage />
            </Route>
            <Route path="/logs">
              <LogsPage />
            </Route>
          </Switch>
        </Layout>
      </Router>
    </AuthProvider>
  );
};

export default App;
