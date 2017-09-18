/*
 * Copyright (c) 2016-2017 Snowplow Analytics Ltd. All rights reserved.
 *
 * This program is licensed to you under the Apache License Version 2.0,
 * and you may not use this file except in compliance with the Apache License Version 2.0.
 * You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the Apache License Version 2.0 is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
 */

/// <reference path="../../../typings/node/node.d.ts" />
/// <reference path="../../../typings/react/react.d.ts" />
/// <reference path="../../../typings/react/react-dom.d.ts" />
/// <reference path="../.././Interfaces.d.ts"/>

import React = require('react');
import ReactDOM = require("react-dom");
import AlertContainer from 'react-alert';
import alertOptions from './AlertOptions'
import axios from 'axios';

var alertContainer = new AlertContainer();

export default React.createClass({
  getInitialState () {
    return {
      iglu_server_uri: '',
      iglu_server_apikey: '',
      disabled: false
    };
  },

  handleChange(evt) {
    if (evt.target.name == 'iglu_server_uri'){
      this.setState({
        iglu_server_uri: evt.target.value
      });
    }
    else if (evt.target.name == 'iglu_server_apikey'){
      this.setState({
        iglu_server_apikey: evt.target.value
      });
    }
  },

  sendFormData()  {
    var alertShow = alertContainer.show
    var _this = this

    var igluServerUri = this.state.iglu_server_uri
    var igluServerApikey = this.state.iglu_server_apikey

    function setInitState() {
      _this.setState({
        iglu_server_uri: "",
        iglu_server_apikey: "",
        disabled: false
      });
    }

    _this.setState({
      disabled: true
    });
    var params = new URLSearchParams();
    params.append('iglu_server_uri', _this.state.iglu_server_uri)
    params.append('iglu_server_apikey', _this.state.iglu_server_apikey)

    axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'
    axios.post('/control-plane/add-external-iglu-server', params, {})
    .then(function (response) {
      setInitState()
      alertShow('Uploaded successfully', {
        time: 2000,
        type: 'success'
      });
    })
    .catch(function (error) {
      setInitState()
      alertShow('Error: ' + error.response.data, {
        time: 2000,
        type: 'error'
      });
    });
  },

  handleSubmit(event) {
    var alertShow = alertContainer.show
    alertShow('Please wait...', {
      time: 2000,
      type: 'info'
    });
    event.preventDefault();
    this.sendFormData();
  },

  render() {
    return (
      <div className="tab-content">
        <h4>Add external Iglu Server: </h4>
        <form action="" onSubmit={this.handleSubmit}>
          <div className="form-group">
            <label htmlFor="iglu_server_uri">Iglu Server URI: </label>
            <input className="form-control" name="iglu_server_uri" ref="iglu_server_uri" required type="text" onChange={this.handleChange} value={this.state.iglu_server_uri} />
          </div>
          <div className="form-group">
            <label htmlFor="iglu_server_apikey">Iglu Server ApiKey: </label>
            <input className="form-control" name="iglu_server_apikey" ref="iglu_server_apikey" required type="text" onChange={this.handleChange} value={this.state.iglu_server_apikey}/>
          </div>
          <div className="form-group">
            <button className="btn btn-primary" type="submit" disabled={this.state.disabled}>Add External Iglu Server</button>
          </div>
        </form>
        <AlertContainer ref={a => alertContainer = a} {...alertOptions} />
      </div>
    );
  }
});
