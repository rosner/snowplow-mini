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
      iglu_server_super_uuid: '',
      disabled: false
    };
  },

  handleChange(evt) {
    this.setState({
      iglu_server_super_uuid: evt.target.value
    });
  },

  sendFormData()  {
    var _this = this
    var alertShow = alertContainer.show
    var igluServerSuperUUID = this.state.iglu_server_super_uuid

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
    params.append('iglu_server_super_uuid', _this.state.iglu_server_super_uuid)
    axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'
    axios.post('/control-plane/add-iglu-server-super-uuid', params, {})
    .then(function (response) {
      setInitState()
      alertShow('Uploaded successfully', {
        time: 2000,
        type: 'success'
      });
    })
    .catch(function (error) {
      setInitState()
      alertShow('Error:' + error.response.data, {
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
        <h4>Add super UUID for local Iglu Server: </h4>
        <form action="" onSubmit={this.handleSubmit}>
          <div className="form-group">
            <label htmlFor="iglu_server_super_uuid">Iglu Server Super UUID: </label>
            <input className="form-control" name="iglu_server_super_uuid" ref="iglu_server_super_uuid" required type="text" onChange={this.handleChange} value={this.state.iglu_server_super_uuid} />
          </div>
          <div className="form-group">
            <button className="btn btn-primary" type="submit" disabled={this.state.disabled}>Add UUID</button>
          </div>
        </form>
        <AlertContainer ref={a => alertContainer = a} {...alertOptions} />
      </div>
    );
  }
});
