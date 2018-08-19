import React, { Component } from 'react';
import io from 'socket.io-client';
import Config from './config';
import Info from './info';
import Display from './display';
import './App.css';

const socket = io('http://localhost:5000', { transports: ['websocket', 'polling', 'flashsocket'] });

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      data: {},
    };
  }

  componentDidMount() {
    socket.on('connect', () => {
    });
    socket.on('data', (data) => {
      this.setState({ data });
    });
  }

  render() {
    return (
      <div className="container-fluid appbg no-padding">
        <div className="row">
          <div className="col-md-1">
            <Config />
          </div>
          <div className="col-md-8">
            <Display />
          </div>
          <div className="col-md-3">
            <Info />
          </div>
        </div>
      </div>
    );
  }
}

export default App;
