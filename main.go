package main

import (
	"os"
	"net"
	"fmt"
	"strings"
	"errors"
	"time"
	"net/http"
	"net/http/httputil"
	"io/ioutil"
	"encoding/json"
	"text/template"
)

// value={{(index (index .NetworkSettings.Ports \"4001/tcp\") 0).HostPort}}

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
	format := ""
	etcd := false
	wait := false

	n := 1
	for {
		if os.Args[n] == "-f" {
			format = os.Args[n+1]
			n += 2;
		} else if os.Args[n] == "-p" {
			format = fmt.Sprintf("{{(index (index .NetworkSettings.Ports \"%s/tcp\") 0).HostPort}}", os.Args[n + 1])
			n += 2
		} else if os.Args[n] == "--etcd" {
			etcd = true
			n++
		} else if os.Args[n] == "-w" {
			wait = true
			n++
		} else {
			break
		}
	}
	if format != "" && etcd {
		format = "value=" + format
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/containers/%s/json", os.Args[n]), nil)

	conn, err := Dial()
  if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := httputil.NewClientConn(conn, nil)

	var f interface{}
	for {
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)


		json.Unmarshal(body, &f)

		if !wait {
			break
		}

		if f != nil {
			state := f.(map[string]interface{})["State"].(map[string]interface{})
			running := state["Running"].(bool)
			if running {
				break
			}
		}

		time.Sleep(50 * time.Millisecond)
	}

	if format != "" {
		tmpl,_ := template.New("test").Parse(format + "\n")
		tmpl.Execute(os.Stdout, f)
	} else {
		b, _ := json.MarshalIndent(f, "", "  ")
		fmt.Printf("%s\n", b)
	}

}
