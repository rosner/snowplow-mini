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
var addExternalIgluServerScript = "add_external_iglu_server.sh"
var addIgluSuperUUIDScript = "add_iglu_server_super_uuid.sh"

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
  http.HandleFunc("/add-external-iglu-server", addExternalIgluServer)
  http.HandleFunc("/add-iglu-server-super-uuid", addIgluServerSuperUUID)
  http.HandleFunc("/check-url", checkUrl)
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

func addExternalIgluServer(resp http.ResponseWriter, req *http.Request) {
  if req.Method == "POST" {
    req.ParseForm()
    if len(req.Form["iglu_server_uri"]) == 0 {
      http.Error(resp, "parameter iglu_server_uri is not given", 400)
      return
    }
    if len(req.Form["iglu_server_apikey"]) == 0 {
      http.Error(resp, "parameter iglu_server_apikey is not given", 400)
      return
    }
    igluServerUri := req.Form["iglu_server_uri"][0]
    igluServerApikey := req.Form["iglu_server_apikey"][0]

    if !isUrlReachable(igluServerUri) {
      http.Error(resp, "Given URL is not reachable", 400)
      return
    }
    if !isValidUuid(igluServerApikey) {
      http.Error(resp, "Given apikey is not valid UUID.", 400)
      return
    }

    shellScriptCommand := []string{scriptsPath + "/" + addExternalIgluServerScript,
                                    igluServerUri,
                                    igluServerApikey,
                                    configPath,
                                    scriptsPath}
    cmd := exec.Command("/bin/bash", shellScriptCommand...)
    err := cmd.Run()
    if err != nil {
      http.Error(resp, err.Error(), 400)
      return
    }
    //restart SP services to get action the external iglu server
    _, err = callRestartSPServicesScript()
    if err != nil {
      http.Error(resp, err.Error(), 400)
      return
    }
    resp.WriteHeader(http.StatusOK)
    io.WriteString(resp, "added successfully")
  }
}

func addIgluServerSuperUUID(resp http.ResponseWriter, req *http.Request) {
  if req.Method == "POST" {
    req.ParseForm()
    if len(req.Form["iglu_server_super_uuid"]) == 0 {
      http.Error(resp, "parameter iglu_server_super_uuid is not given", 400)
      return
    }
    igluServerSuperUUID := req.Form["iglu_server_super_uuid"][0]
    if !isValidUuid(igluServerSuperUUID) {
      http.Error(resp, "Given apikey is not valid UUID", 400)
      return
    }
    shellScriptCommand := []string{scriptsPath + "/" + addIgluSuperUUIDScript,
                                   igluServerSuperUUID,
                                   configPath}
    cmd := exec.Command("/bin/bash", shellScriptCommand...)
    err := cmd.Run()
    if err != nil {
      http.Error(resp, err.Error(), 400)
      return
    }
    //restart SP services to get action the added Iglu apikey
    _, err = callRestartSPServicesScript()
    if err != nil {
      http.Error(resp, err.Error(), 400)
      return
    }
    resp.WriteHeader(http.StatusOK)
    io.WriteString(resp, "added successfully")
  }
}
