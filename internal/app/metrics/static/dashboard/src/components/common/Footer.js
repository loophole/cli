import React from "react";

import { Row, Col } from "reactstrap";

class Footer extends React.Component {
  render() {
    return (
      <footer className="footer">
        <Row className="align-items-center justify-content-xl-between">
          <Col xl="6">
            <div className="copyright text-center text-xl-left text-muted">
              © 2020{" "}
              <a
                className="font-weight-bold ml-1"
                href="https://loophole.cloud"
                rel="noopener noreferrer"
                target="_blank"
              >
                Loophole
              </a>
            </div>
          </Col>
        </Row>
      </footer>
    );
  }
}

export default Footer;
