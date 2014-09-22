package main

import (
	"os"
	_ "fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	_ "reflect"
	"text/template"
)

func main() {
  // url := "http://tethys:2375/containers/265926ce9f05/json"
	url :- "unix:/"
	
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	var f interface{}
	
	json.Unmarshal(body, &f)
	
	tmpl,_ := template.New("test").Parse("value={{(index (index .NetworkSettings.Ports \"1521/tcp\") 0).HostPort}}")
	tmpl.Execute(os.Stdout, f)
	
	// fmt.Printf("===> " + reflect.TypeOf(f).Name())
}
