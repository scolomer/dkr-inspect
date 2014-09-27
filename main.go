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
)

func main() {
	req, err := http.NewRequest("GET", "/containers/2a03159ba971/json", nil)

	// conn, err := net.Dial("unix", "/var/run/docker.sock")
	conn, err := net.Dial("tcp", "tethys:2375")
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
