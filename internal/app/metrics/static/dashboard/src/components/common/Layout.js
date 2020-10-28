import React from "react";
import { Route, Switch, Redirect } from "react-router-dom";
// reactstrap components
import { Container } from "reactstrap";
// core components
import Navbar from "./Navbar";
import Header from "./Header";
import Footer from "./Footer";
import Sidebar from "./Sidebar";

import routes from "../../routes";

class Layout extends React.Component {
  componentDidUpdate(e) {
    document.documentElement.scrollTop = 0;
    document.scrollingElement.scrollTop = 0;
    this.refs.mainContent.scrollTop = 0;
  }
  getRoutes = routes => {
    return routes.map((prop, key) => {
      return (
        <Route
          path={prop.path}
          component={prop.component}
          key={key}
        />
      );
    });
  };
  getBrandText = path => {
    for (let i = 0; i < routes.length; i++) {
      if (
        this.props.location.pathname.indexOf(
          routes[i].path
        ) !== -1
      ) {
        return routes[i].name;
      }
    }
    return "Loophole";
  };
  render() {
    return (
      <>
        <Sidebar
          {...this.props}
          routes={routes}
          logo={{
            innerLink: "/index",
            imgSrc: require("../../assets/img/brand/logo.png"),
            imgAlt: "..."
          }}
        />
        <div className="main-content" ref="mainContent">
          <Header />
          <Navbar
            {...this.props}
            brandText={this.getBrandText(this.props.location.pathname)}
          />
          <Switch>
            {this.getRoutes(routes)}
            <Redirect from="*" to="/index" />
          </Switch>
          <Container fluid>
            <Footer />
          </Container>
        </div>
      </>
    );
  }
}

export default Layout;
