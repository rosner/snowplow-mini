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
  "io"
  "net/http"
  "log"
  "os/exec"
  "os"
  "flag"
)

// script file names
var restartServicesScript = "restart_SP_services.sh"

// global variables for paths from flags
var scriptsPath string
var enrichmentsPath string
var configPath string

func main() {
  scriptsPathFlag := flag.String("scriptsPath", "", "path for control-plane-api scripts")
  enrichmentsPathFlag := flag.String("enrichmentsPath", "", "path for enrichment files")
  configPathFlag := flag.String("configPath", "", "path for config files")
  flag.Parse()
  scriptsPath = *scriptsPathFlag
  enrichmentsPath = *enrichmentsPathFlag
  configPath = *configPathFlag

  http.HandleFunc("/restart-services", restartSPServices)
  http.HandleFunc("/upload-enrichments", uploadEnrichments)
  log.Fatal(http.ListenAndServe(":10000", nil))
}

func restartSPServices(resp http.ResponseWriter, req *http.Request) {
  if (req.Method == "PUT") {
    _, err := callRestartSPServicesScript()
    if err != nil {
      http.Error(resp, err.Error(), 400)
      return
    } else {
      resp.WriteHeader(http.StatusOK)
      io.WriteString(resp, "OK")
    }
  }
}

func uploadEnrichments(resp http.ResponseWriter, req *http.Request) {
  if req.Method == "POST" {
    req.ParseMultipartForm(32 << 20)
    file, handler, err := req.FormFile("enrichmentjson")
    if err != nil {
      http.Error(resp, err.Error(), 400)
      return
    }
    defer file.Close()
    f, err := os.OpenFile(enrichmentsPath + "/" + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
      http.Error(resp, err.Error(), 400)
      return
    }
    defer f.Close()
    fileContentBytes, err := ioutil.ReadAll(file)
    fileContent := string(fileContentBytes)

    if !isJSON(fileContent) {
      http.Error(resp, "JSON is not valid", 400)
      return
    }

    io.WriteString(f, fileContent)
    //restart SP services to get action the enrichments
    _, err = callRestartSPServicesScript()
    resp.WriteHeader(http.StatusOK)
    if err != nil {
      http.Error(resp, err.Error(), 400)
      return
    }
    resp.WriteHeader(http.StatusOK)
    io.WriteString(resp, "uploaded successfully")
    return
  }
}
