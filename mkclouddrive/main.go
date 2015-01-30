// Modified from https://raw.githubusercontent.com/kelseyhightower/isod/master/main.go
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"bytes"
  "flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
  "text/template"
)

var outfile string
var publickeyfile string

func init() {
  flag.StringVar(&outfile, "outfile", "configdrive.iso", "file name to output the config drive")
  flag.StringVar(&publickeyfile, "publickeyfile", "", "public key file path")
}

func genisoimageHandler(w http.ResponseWriter,user_data string) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	configPath := path.Join(tmpdir, "/openstack/latest")
	userData := path.Join(configPath, "user_data")
	err = os.MkdirAll(configPath, 755)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpdir)
  ioutil.WriteFile(userData, []byte(user_data), 0x777)

	cmd := exec.Command("/usr/local/bin/mkisofs", "-R", "-V", "config-2", tmpdir)
	var iso bytes.Buffer
	cmd.Stdout = &iso
	err = cmd.Run()
	if err != nil {
		log.Println(err.Error())
		return
	}
  ioutil.WriteFile(outfile, iso.Bytes(), 0x777)
}

type CloudConfig struct {
  Ssh_authorized_keys string
}

func main() {
  flag.Parse()
  r,err := ioutil.ReadFile(publickeyfile)
  if err != nil { panic(err) }
  cloudconfig := CloudConfig{string(r)}
  cloudconfigTmpl,err := template.New("node").Parse(`#cloud-config
ssh_authorized_keys:
    - {{.Ssh_authorized_keys}}`)
  if err != nil { panic(err) }
  var buf_cloudconfig bytes.Buffer
  err = cloudconfigTmpl.Execute(&buf_cloudconfig, cloudconfig)
  if err != nil { panic(err) }
	genisoimageHandler(nil,buf_cloudconfig.String())
}
