import React, { Component} from 'react';
import "./field.css"
import $ from 'jquery';

class Config extends Component {
  constructor(props){
    super(props);
    this.state = {
      InFreq:400,
      OutFreq:40,
      InMag:100,
      OutMag:100
    }
    this.inputval = this.inputval.bind(this)
    this.configPressed = this.configPressed.bind(this)
  }
  configPressed(event){
    $.ajax({
      type: 'POST',
      url: 'http://localhost:5000/config',
      data:JSON.stringify({
        "InFreq":this.state.InFreq,
        "OutFreq":this.state.OutFreq,
        "InMag":this.state.InMag,
        "OutMag":this.state.OutMag
        })
    }
    )
  }
  inputval(event){
    var set = Number(event.target.value);
    this.setState({
      [event.target.name]:set
    });
  }
  render() {
    const configItems = 
    <div className="ConfigFields">
      <label>
        InFreq:<input name="InFreq" 
                      type="text" 
                      value={this.state.InFreq} 
                      onChange={this.inputval}/>
      </label>
      <label>
        OutFreq:<input name="OutFreq" 
                      type="text" 
                      value={this.state.OutFreq} 
                      onChange={this.inputval}/>
      </label>
      <label>
        InMag:<input name="InMag" 
                      type="text" 
                      value={this.state.InMag} 
                      onChange={this.inputval}/>
      </label>
      <label>
        OutMag:<input name="OutMag" 
                      type="text" 
                      value={this.state.OutMag} 
                      onChange={this.inputval}/>
      </label>
      </div>
      ;
    return (
      <div className="container-fluid configbg">
        <header className = "header">
            Configuration
        </header>
        <div>
          <ul>{configItems}</ul>
        </div>
        <button
        onClick={this.configPressed}
        title="Write Config"
        color="#841584"
        >Write</button>
        </div>
    );
  }
}

export default Config;