package main

import (
	"os"
	"net"
	_ "fmt"
	"net/http"
	"net/http/httputil"
	"io/ioutil"
	"encoding/json"
	_ "reflect"
	"text/template"
	"strings"
	"errors"
)

func Dial()(net.Conn, error) {
	url := os.Getenv("DOCKER_HOST")
	if (url == "") {
		return net.Dial("unix", "/var/run/docker.sock")
	}

	if (strings.HasPrefix(url, "tcp://")) {
		return net.Dial("tcp", url[6:len(url)])
	}

	return nil, errors.New("Unhandled form : " + url)
}

func main() {
	req, err := http.NewRequest("GET", "/containers/etcd/json", nil)

	conn, err := Dial()
  if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := httputil.NewClientConn(conn, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var f interface{}

	json.Unmarshal(body, &f)

	tmpl,_ := template.New("test").Parse("value={{(index (index .NetworkSettings.Ports \"4001/tcp\") 0).HostPort}}\n")
	tmpl.Execute(os.Stdout, f)
}
