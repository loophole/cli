import React from "react";

import axios from 'axios';
import moment from 'moment';

// reactstrap components
import {
  Card,
  CardHeader,
  Table,
  Container,
  Row
} from "reactstrap";
// core components

class Tables extends React.Component {
  state = {
    loading: false,
    error: false,
    events: []
  }

  constructor() {
    super();

    this.fetchEvents = this.fetchEvents.bind(this);
  }

  async componentDidMount() {
    await this.fetchEvents();
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
        events: response.data
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
    const eventRows = this.state.events
      .sort((a, b) => moment(b.timestamp).diff(a.timestamp))
      .map((event) => {
        return (
          <tr key={event.timestamp}>
            <td>{moment(event.timestamp).format()}</td>
            <td>{event.message}</td>
            <td>{event.siteId}</td>
            <td>{event.sessionId}</td>
          </tr>
        )
      })

    return (
      <>
        {/* Page content */}
        <Container className="mt--7" fluid>
          {/* Table */}
          <Row>
            <div className="col">
              <Card className="shadow">
                <CardHeader className="border-0">
                  <h3 className="mb-0">Card tables</h3>
                </CardHeader>
                <Table className="align-items-center table-flush" responsive>
                  <thead className="thead-light">
                    <tr>
                      <th scope="col">Time</th>
                      <th scope="col">Message</th>
                      <th scope="col">Site</th>
                      <th scope="col">Session</th>
                    </tr>
                  </thead>
                  <tbody>
                    {eventRows}
                  </tbody>
                </Table>
              </Card>
            </div>
          </Row>
        </Container>
      </>
    );
  }
}

export default Tables;
