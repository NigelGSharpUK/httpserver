// This program forwards requests from js in the browser
// to the core json-rpc server on port 8332
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler2) // each request calls handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the requested URL.
/*func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}*/

func handler2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received, handler2 called")

	// Get the json request string
	json2 := r.Body

	// Get the request origin
	origin := r.Header.Get("Origin")
	fmt.Println(origin)
	fmt.Println(json2)

	// Talk to core
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8332/", json2)
	req.SetBasicAuth("", "notmypassword")
	req.Header.Add("content-type", "application/JSON;")

	res, e := http.DefaultClient.Do(req)
	if e != nil {
		fmt.Println("ERROR: Could not talk to core on local machine.")
		fmt.Println("Perhaps core is not running?")
		fmt.Println("Perhaps core is still downloading the blockchain?")
		fmt.Println("Perhaps the firewall is interfering?")
		fmt.Println(e)
	} else {
		defer res.Body.Close()
		// Respond with core's response
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
		// The following tells the browser to allow requests from 127.0.0.1
		// This helps with CORS restrictions - Cross-Origin Resource Sharing
		if origin == "http://127.0.0.1" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
		w.Write(body)
		//fmt.Fprintf(w, "%q\n", string(body))
	}

}
