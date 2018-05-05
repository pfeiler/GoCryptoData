package main

import (
	"net/http"
	"strings"
	"encoding/json"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderMenu(w)
	})

	http.HandleFunc("/cryptodata/", func(w http.ResponseWriter, r *http.Request) {
	    cryptoId := strings.SplitN(r.URL.Path, "/", 3)[2]

	    data, err := query(cryptoId)
	    if err != nil {
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        return
	    }

	    w.Header().Set("Content-Type", "application/json; charset=utf-8")
	    json.NewEncoder(w).Encode(data)
		w.Write([]byte("<br/><br/><a href=\"/\">Back</a><br/>"))
	})
	
    http.ListenAndServe(":8080", nil)

}

func renderMenu(w http.ResponseWriter) {
	w.Write([]byte("<a href=\"/cryptodata/1\">BTC price</a><br/>"))
	w.Write([]byte("<a href=\"/cryptodata/1027\">ETH price</a><br/>"))
}


type cryptoData struct {
    Data struct {
		Name string `json:"name"`
	    Data struct {
		    Data struct {
		        Price float64 `json:"price"`
			} `json:"USD"`
	    } `json:"quotes"`
 	} `json:"data"`
}

func query(cryptoId string) (cryptoData, error) {
    resp, err := http.Get("https://api.coinmarketcap.com/v2/ticker/" + cryptoId)
    if err != nil {
        return cryptoData{}, err
    }

    defer resp.Body.Close()

    var data cryptoData

    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return cryptoData{}, err
    }

    return data, nil
}
