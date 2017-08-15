package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/liimaorg/liimactl/client/util"
)

//MockClient used for mocking http client
type MockClient struct{}

var ts httptest.Server

//CloseMockClient will close the started test server
func CloseMockClient() {
	defer ts.Close()
}

// NewMockClient creates a new liima client from the config
func NewMockClient(config *Config) (*Client, error) {

	//Set Server with handlers
	mux := serverMuxHandler()
	ts := httptest.NewServer(mux)
	//set localhost in config
	config.Host = ts.URL + "/"

	return &Client{
		client: http.DefaultClient,
		config: config,
		url:    ts.URL,
	}, nil
}

//Mux Http handlers
func serverMuxHandler() *http.ServeMux {
	r := http.NewServeMux()

	//ToDo: Deployment test handler
	//r.HandleFunc("/deplyoments", createDeploymentHandler).Methods("POST")

	//Get Deployment test handler
	r.HandleFunc("/resources/deployments", listGetDeploymentHandler)

	//Hostname test handler
	r.HandleFunc("/resources/hostNames", listHostnameHandler)

	return r
}

//Hostname test handler
func listHostnameHandler(w http.ResponseWriter, r *http.Request) {

	//Create hostname response
	respHostname := Hostnames{{}}

	//Split the requested Url
	cuttedURL := strings.Split(r.URL.String(), "?")[1]
	//Split the subcommands
	commands := strings.Split(cuttedURL, "&")
	//Set the requested command-value as respond if a tag exits in the hostename
	for _, command := range commands {
		key := strings.Split(command, "=")[0]
		value := strings.Split(command, "=")[1]
		util.SetValueIfTagExists(&respHostname[0], key, value)
	}
	//Send hostname as response
	hostname, err := json.Marshal(respHostname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(hostname)
}

//Get deployment test handler
func listGetDeploymentHandler(w http.ResponseWriter, r *http.Request) {

	//Create  response
	response := Deployments{{}}

	//Split the requested Url
	cuttedURL := strings.Split(r.URL.String(), "?")[1]
	//Split the subcommands
	commands := strings.Split(cuttedURL, "&")
	//Set the requested command-value as respond if a tag exits in the hostename
	for _, command := range commands {
		key := strings.Split(command, "=")[0]
		value := strings.Split(command, "=")[1]
		util.SetValueIfTagExists(&response[0], key, value)
	}
	//Send response
	deployment, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(deployment)
}

//DoRequest set up a json for the given url and calls the llima client.
//Method: http.MethodX
//URL: Resturl
//The bodyType will be marshaled to the rest body, depending the method
//The result will be unmarshaled to the responseType
func (c *MockClient) DoRequest(method string, url string, bodyType interface{}, responseType interface{}) error {

	//Setup body if MethodPost
	bData := []byte{}
	if method == http.MethodPost {
		bData, err := json.Marshal(bodyType)
		if err != nil {
			log.Fatal(err)
		}
		_ = bData
	}

	var bodydata = bytes.NewBuffer(bData)

	// Dump response
	data, err := ioutil.ReadAll(bodydata)
	if err != nil {
		return err
	}

	//Unmarshal json respond to responseType
	return json.Unmarshal(data, responseType)

}
