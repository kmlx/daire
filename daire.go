/*
  Copyright 2015 Adrian Stanescu. All rights reserved.
  Use of this source code is governed by the MIT License (MIT) that can be found in the LICENSE file.
 */

/*
  daire
  Go program that acts as a single host reverse proxy

  Usage:
  go run ./daire.go [listen host and port] [to host and port]
 */

package main

import (
  "log"
  "net/http"
  "net/http/httputil"
  "net/url"
  "os"
)

type SingleRequestTransport struct {
  // bottleneck chan bool
}

func (s *SingleRequestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
  // s.bottleneck <- true
  response, err := http.DefaultTransport.RoundTrip(req)
  // <- s.bottleneck
  return response, err;
}


func main() {
  if len(os.Args) != 3 {
    log.Fatal("Usage: go run ./daire.go [listen host and port] [to host and port]")
  }

  hostAndPort   := os.Args[1]
  proxyHost     := os.Args[2]

  proxy         := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: proxyHost, Path: "/"})
  proxy.Transport = &SingleRequestTransport{/*bottleneck: make(chan bool, 1)*/}

  server        := &http.Server{
    Addr:           hostAndPort,
    Handler:        proxy,
  }

  log.Print("Daire proxy listening on: ", hostAndPort, ", redirecting to: ", proxyHost)
  err := server.ListenAndServe()
  if (err != nil) {
    log.Fatal(err)
  }
}
