import React from "react";
// node.js library that concatenates classes (strings)
import classnames from "classnames";
// javascipt plugin for creating charts
import Chart from "chart.js";
// react plugin used to create charts
import { Line } from "react-chartjs-2";

import { QRCode } from "react-qr-svg";
import axios from 'axios';
import moment from 'moment';
// reactstrap components
import {
  Card,
  CardHeader,
  CardBody,
  NavItem,
  NavLink,
  Nav,
  Container,
  Row,
  Col
} from "reactstrap";

// core components
import {
  chartOptions,
  parseOptions,
  formatData,
  linearChart
} from "../variables/charts.js";

class Index extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      activeNav: 1,
      dataToPresent: "bytesIn",
      data: {
        bytesIn: {
          labels: [],
          values: []
        },
        bytesOut: {
          labels: [],
          values: []
        }
      },
      url: "",
      upSince: ""
    };
    if (window.Chart) {
      parseOptions(Chart, chartOptions());
    }
  }

  async componentDidMount() {
    await this.fetchCurrent();
    await this.fetchMetrics();
  }

  async fetchCurrent() {
    this.setState({
      ...this.state,
      loading: true,
      error: false
    });
    try {
      const response = await axios.get('/api/current');
      this.setState({
        ...this.state,
        url: response.data.url,
        upSince: moment(response.data.startedAt).fromNow()
      });
    } catch (err) {
      console.log(err);
      this.setState({
        ...this.state,
        loading: false,
        error: true
      });
    }
  }

  async fetchMetrics() {
    this.setState({
      ...this.state,
      loading: true,
      error: false
    });
    try {
      const response = await axios.get('/api/metrics');
      const grouped = response.data
        .map((metric) => {
          return {
            ...metric,
            timestamp: moment(metric.timestamp).startOf('second').format('DD/MM/YY HH:mm:ss')
          }
        })
        .reduce((prev, curr) => {
          if (prev[curr.name] && prev[curr.name][curr.timestamp])
            return {
              ...prev,
              [curr.name]: {
                ...prev[curr.name],
                [curr.timestamp]: prev[curr.name][curr.timestamp] + curr.value
              }
            };
          return {
            ...prev,
            [curr.name]: {
              ...prev[curr.name],
              [curr.timestamp]: curr.value
            }
          };
        }, {});
      const emptyBytesIn = { labels: [], values: [] };
      const bytesIn = grouped.bytesIn ? Object.keys(grouped.bytesIn).reduce((agg, curr) => {
        return {
          ...agg,
          labels: [...agg.labels, curr],
          values: [...agg.values, grouped.bytesIn[curr]]
        }
      }, emptyBytesIn) : emptyBytesIn;
      const emptyBytesOut = { labels: [], values: [] };
      const bytesOut = grouped.bytesOut ? Object.keys(grouped.bytesOut).reduce((agg, curr) => {
        return {
          ...agg,
          labels: [...agg.labels, curr],
          values: [...agg.values, grouped.bytesOut[curr]]
        }
      }, emptyBytesOut) : emptyBytesOut;

      this.setState({
        ...this.state,
        loading: false,
        error: false,
        data: {
          bytesIn,
          bytesOut
        }
      });
    } catch (err) {
      console.log(err);
      this.setState({
        ...this.state,
        loading: false,
        error: true
      });
    }
  }


  toggleNavs = (e, index) => {
    e.preventDefault();
    this.setState({
      ...this.state,
      activeNav: index,
      dataToPresent: index === 1 ? "bytesIn" : "bytesOut"
    });
  };
  render() {
    return (
      <>
        {/* Page content */}
        <Container className="mt--7" fluid>
          <Row>
            <Col className="mb-5 mb-xl-0" xl="8">
              <Card className="bg-gradient-default shadow">
                <CardHeader className="bg-transparent">
                  <Row className="align-items-center">
                    <div className="col">
                      <h6 className="text-uppercase text-light ls-1 mb-1">
                        Overview
                      </h6>
                      <h2 className="text-white mb-0">Bytes transferred</h2>
                    </div>
                    <div className="col">
                      <Nav className="justify-content-end" pills>
                        <NavItem>
                          <NavLink
                            className={classnames("py-2 px-3", {
                              active: this.state.activeNav === 1
                            })}
                            href="#pablo"
                            onClick={e => this.toggleNavs(e, 1)}
                          >
                            <span className="d-none d-md-block">In</span>
                            <span className="d-md-none">I</span>
                          </NavLink>
                        </NavItem>
                        <NavItem>
                          <NavLink
                            className={classnames("py-2 px-3", {
                              active: this.state.activeNav === 2
                            })}
                            data-toggle="tab"
                            href="#pablo"
                            onClick={e => this.toggleNavs(e, 2)}
                          >
                            <span className="d-none d-md-block">Out</span>
                            <span className="d-md-none">O</span>
                          </NavLink>
                        </NavItem>
                      </Nav>
                    </div>
                  </Row>
                </CardHeader>
                <CardBody>
                  {/* Chart */}
                  <div className="chart">
                    <Line
                      data={formatData("Bytes transferred", this.state.data[this.state.dataToPresent].labels, this.state.data[this.state.dataToPresent].values)}
                      options={linearChart.options}
                      getDatasetAtEvent={e => console.log(e)}
                    />
                  </div>
                </CardBody>
              </Card>
            </Col>
            <Col xl="4">
              <Card className="shadow">
                <CardHeader className="bg-transparent">
                  <Row className="align-items-center">
                    <div className="col">
                      <h6 className="text-uppercase text-muted ls-1 mb-1">
                        Share your website
                      </h6>
                      <h2 className="mb-0">
                        {this.state.url}
                      </h2>
                    </div>
                  </Row>
                </CardHeader>
                <CardBody>
                  <div className="chart">
                    <QRCode
                      bgColor="#FFFFFF"
                      fgColor="#000000"
                      level="Q"
                      style={{ width: 300 }}
                      value={this.state.url} />
                  </div>
                </CardBody>
              </Card>
            </Col>
          </Row>
        </Container>
      </>
    );
  }
}

export default Index;
