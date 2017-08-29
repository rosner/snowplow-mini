/**
 * Copyright (c) 2014-2017 Snowplow Analytics Ltd.
 * All rights reserved.
 *
 * This program is licensed to you under the Apache License Version 2.0,
 * and you may not use this file except in compliance with the Apache
 * License Version 2.0.
 * You may obtain a copy of the Apache License Version 2.0 at
 * http://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the Apache License Version 2.0 is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied.
 *
 * See the Apache License Version 2.0 for the specific language
 * governing permissions and limitations there under.
 */

package main

import (
  "encoding/json"
  "net/http"
  "net"
  "os/exec"
  "context"
)

// restarts services
func callRestartSPServicesScript() (string, error){
  shellScriptCommand := []string{scriptsPath + "/" +  restartServicesScript}
  cmd := exec.Command("/bin/bash", shellScriptCommand...)
  err := cmd.Run()
  if err != nil {
    return "ERR", err
  }
  return "OK", err
}

// check if JSON string is valid or not
func isJSON(jsonString string) bool {
  var js map[string]interface{}
  return json.Unmarshal([]byte(jsonString), &js) == nil
}
