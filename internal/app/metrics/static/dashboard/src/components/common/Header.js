import React from "react";

// reactstrap components
import { Card, CardBody, CardTitle, Container, Row, Col } from "reactstrap";
import axios from 'axios';

class Header extends React.Component {
  state = {
    loading: false,
    error: false,
    bytesTransferedIn: 0,
    bytesTransferedOut: 0,
    siteNamesUsed: [],
    sessionsUsed: []
  }

  constructor() {
    super();

    this.fetchMetrics = this.fetchMetrics.bind(this);
    this.fetchEvents = this.fetchEvents.bind(this);
  }

  async componentDidMount() {
    await this.fetchMetrics();
    await this.fetchEvents();
  }

  async fetchMetrics() {
    this.setState({
      ...this.state,
      loading: true,
      error: false
    });
    try {
      const response = await axios.get('/api/metrics');
      const metrics = response.data.reduce((prev, curr) => {
        const result = { ...prev };
        if (!result.siteNamesUsed.includes(curr.siteId))
          result.siteNamesUsed = [...result.siteNamesUsed, curr.siteId];
        if (!result.sessionsUsed.includes(curr.sessionId))
          result.sessionsUsed = [...result.sessionsUsed, curr.sessionId];
        if (curr.name === "bytesIn")
          return { ...result, bytesIn: result.bytesIn += curr.value };
        if (curr.name === "bytesOut")
          return { ...result, bytesOut: result.bytesOut += curr.value };
        return result;
      }, { bytesIn: 0, bytesOut: 0, siteNamesUsed: [], sessionsUsed: [] })
      this.setState({
        ...this.state,
        loading: false,
        error: false,
        bytesTransferedIn: metrics.bytesIn,
        bytesTransferedOut: metrics.bytesOut,
        siteNamesUsed: metrics.siteNamesUsed,
        sessionsUsed: metrics.sessionsUsed
      });
    } catch (err) {
      this.setState({
        ...this.state,
        loading: false,
        error: true
      });
    }
  }


  async fetchEvents() {
    this.setState({
      ...this.state,
      loading: true,
      error: false
    });
    try {
      const response = await axios.get('/api/events');
      this.setState({
        ...this.state,
        loading: false,
        error: false,
        eventsCount: response.data.length
      });
    } catch (err) {
      this.setState({
        ...this.state,
        loading: false,
        error: true
      });
    }
  }

  render() {
    return (
      <>
        <div className="header bg-gradient-info pb-8 pt-5 pt-md-8">
          <Container fluid>
            <div className="header-body">
              {/* Card stats */}
              <Row>
                <Col>
                  <Card className="card-stats mb-4 mb-xl-0">
                    <CardBody>
                      <Row>
                        <div className="col">
                          <CardTitle
                            tag="h5"
                            className="text-uppercase text-muted mb-0"
                          >
                            Incoming bytes
                          </CardTitle>
                          <span className="h2 font-weight-bold mb-0">
                            {this.state.bytesTransferedIn}
                          </span>
                        </div>
                        <Col className="col-auto">
                          <div className="icon icon-shape bg-primary text-white rounded-circle shadow">
                            <i className="fas fa-arrow-down" />
                          </div>
                        </Col>
                      </Row>
                    </CardBody>
                  </Card>
                </Col>
                <Col>
                  <Card className="card-stats mb-4 mb-xl-0">
                    <CardBody>
                      <Row>
                        <div className="col">
                          <CardTitle
                            tag="h5"
                            className="text-uppercase text-muted mb-0"
                          >
                            Outgoing bytes
                          </CardTitle>
                          <span className="h2 font-weight-bold mb-0">
                            {this.state.bytesTransferedOut}
                          </span>
                        </div>
                        <Col className="col-auto">
                          <div className="icon icon-shape bg-success text-white rounded-circle shadow">
                            <i className="fas fa-arrow-up" />
                          </div>
                        </Col>
                      </Row>
                      {/* <p className="mt-3 mb-0 text-muted text-sm">
                        <span className="text-danger mr-2">
                          <i className="fas fa-arrow-down" /> 3.48%
                        </span>{" "}
                        <span className="text-nowrap">Since last week</span>
                      </p> */}
                    </CardBody>
                  </Card>
                </Col>
                <Col>
                  <Card className="card-stats mb-4 mb-xl-0">
                    <CardBody>
                      <Row>
                        <div className="col">
                          <CardTitle
                            tag="h5"
                            className="text-uppercase text-muted mb-0"
                          >
                            Used site names
                          </CardTitle>
                          <span className="h2 font-weight-bold mb-0">
                            {this.state.siteNamesUsed.length}
                          </span>
                        </div>
                        <Col className="col-auto">
                          <div className="icon icon-shape bg-info text-white rounded-circle shadow">
                            <i className="fas fa-globe" />
                          </div>
                        </Col>
                      </Row>
                      {/* <p className="mt-3 mb-0 text-muted text-sm">
                        <span className="text-warning mr-2">
                          <i className="fas fa-arrow-down" /> 1.10%
                        </span>{" "}
                        <span className="text-nowrap">Since yesterday</span>
                      </p> */}
                    </CardBody>
                  </Card>
                </Col>
                <Col>
                  <Card className="card-stats mb-4 mb-xl-0">
                    <CardBody>
                      <Row>
                        <div className="col">
                          <CardTitle
                            tag="h5"
                            className="text-uppercase text-muted mb-0"
                          >
                            Used sessions
                          </CardTitle>
                          <span className="h2 font-weight-bold mb-0">
                            {this.state.sessionsUsed.length}
                          </span>
                        </div>
                        <Col className="col-auto">
                          <div className="icon icon-shape bg-yellow text-white rounded-circle shadow">
                            <i className="fas fa-keyboard" />
                          </div>
                        </Col>
                      </Row>
                      {/* <p className="mt-3 mb-0 text-muted text-sm">
                        <span className="text-warning mr-2">
                          <i className="fas fa-arrow-down" /> 1.10%
                        </span>{" "}
                        <span className="text-nowrap">Since yesterday</span>
                      </p> */}
                    </CardBody>
                  </Card>
                </Col>
                <Col>
                  <Card className="card-stats mb-4 mb-xl-0">
                    <CardBody>
                      <Row>
                        <div className="col">
                          <CardTitle
                            tag="h5"
                            className="text-uppercase text-muted mb-0"
                          >
                            Events
                          </CardTitle>
                          <span className="h2 font-weight-bold mb-0">
                            {this.state.eventsCount}
                          </span>
                        </div>
                        <Col className="col-auto">
                          <div className="icon icon-shape bg-danger text-white rounded-circle shadow">
                            <i className="fas fa-exclamation" />
                          </div>
                        </Col>
                      </Row>
                      {/* <p className="mt-3 mb-0 text-muted text-sm">
                        <span className="text-success mr-2">
                          <i className="fas fa-arrow-up" /> 12%
                        </span>{" "}
                        <span className="text-nowrap">Since last month</span>
                      </p> */}
                    </CardBody>
                  </Card>
                </Col>
              </Row>
            </div>
          </Container>
        </div>
      </>
    );
  }
}

export default Header;
