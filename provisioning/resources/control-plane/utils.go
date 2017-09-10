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
  "regexp"
  "errors"
  "time"
  "io/ioutil"
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

// check if given URL is reachable
func isUrlReachable(url string) bool {
  _, err := http.Get("http://" + url)
  if err != nil {
    return false
  }
  return true
}

// check whether given UUID follows the correct format
func isValidUuid(uuid string) bool {
    r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
    return r.MatchString(uuid)
}

// returns public IP if the host machine is EC2 instance
func getPublicEC2IP() (string, error) {
  // URL of the instance meta service of AWS EC2
  var urlForCheckingPubIP = "http://169.254.169.254/latest/meta-data/public-ipv4"
  var netClient = &http.Client{
    Timeout: time.Second * 5,
  }

  resp, err := netClient.Get(urlForCheckingPubIP)
  if err != nil {
    return "", err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  return string(body), nil
}

// get IP addresses of the given domain name
func getDomainNameIP(domainName string) ([]string, error) {
  var (
    ipAddresses []string
    ctx context.Context
    cancel context.CancelFunc
  )

  ctx, cancel = context.WithCancel(context.Background())
  defer cancel()

  addrs, err := net.DefaultResolver.LookupIPAddr(ctx, domainName)
  if err != nil {
    return nil, err
  }

  for _, ipnet := range addrs {
    if ipnet.IP.To4() != nil {
      ipAddresses = append(ipAddresses, ipnet.IP.String())
    }
  }

  return ipAddresses, nil
}

// return whether given domain name resolves to the host IP or not
func checkHostDomainName(domainName string) error{
  // if host machine is ec2 instance,
  // public IP must be got from instance meta service of AWS EC2
  hostIPAddress, err := getPublicEC2IP()
  if err != nil {
    return errors.New("Error while getting host ip addresses")
  }

  domainIPAdresses, err := getDomainNameIP(domainName)
  if err != nil {
    return errors.New("Error while getting ip addresses of domain")
  }

  for _, domainIP := range domainIPAdresses {
    if domainIP == hostIPAddress {
      return nil
    }
  }

  return errors.New("Given domain name does not redirect to host")
}
