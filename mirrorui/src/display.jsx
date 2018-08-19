import React, { Component } from 'react';
import "./field.css"
import PointCloud from './pointcloud.jsx';

// source: Natural Earth http://www.naturalearthdata.com/ via geojson.xyz

class Display extends Component {
  render() {
    return (
      <div className="container-fluid displaybg">
        <div className="row">
          <header className = "header">
            Display
          </header>
        </div>
        <div className="row pcloud">
          <PointCloud></PointCloud>
        </div>
      </div>
    );
  }
}

export default Display;
