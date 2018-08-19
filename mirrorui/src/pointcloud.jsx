
import React, { Component } from 'react';
import DeckGL, {PointCloudLayer,COORDINATE_SYSTEM} from 'deck.gl';
import "./field.css"

const INITIAL_VIEW_STATE = {
  latitude: 40,
  longitude: -100,
  zoom: 7,
  bearing: 30,
  pitch: 30
};

const data = [
  {position:[0,0,0],normal:[0,0,0],color:[255,0,0]},
  {position:[10000,0,0],normal:[0,0,0],color:[255,255,0]},
  {position:[0,10000,0],normal:[0,0,0],color:[255,255,0]},
  {position:[0,0,10000],normal:[0,0,0],color:[255,255,0]},
  {position:[0,10000,10000],normal:[0,0,0],color:[255,255,255]},
  {position:[10000,10000,0],normal:[0,0,0],color:[255,255,255]},
  {position:[10000,0,10000],normal:[0,0,0],color:[255,255,255]},
  {position:[10000,10000,10000],normal:[0,0,0],color:[255,0,0]},
]


class PointCloud extends Component {
  render() {
    const layers=[
      new PointCloudLayer({
        id: 'point-cloud-layer',
        data:data,
        pickable: false,
        coordinateSystem: COORDINATE_SYSTEM.METER_OFFSETS,
        coordinateOrigin: [-99.8, 40, 0],
        radiusPixels: 9,
        getPosition: d => d.position,
        getNormal: d => d.normal,
        getColor: d => d.color,
        lightSettings: {}
      })
    ];
    return (
      <div>
        <DeckGL
          width='100%'
          height='100%'
          initialViewState={INITIAL_VIEW_STATE}
          controller={true}
          layers={layers}
        >
        </DeckGL>
      </div>
    );
  }
}

export default PointCloud;
