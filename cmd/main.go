package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
)

const ID_BTC = "1"
const ID_ETH = "1027"

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

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderMenu(w)
	})

	http.HandleFunc("/cryptodata/", func(w http.ResponseWriter, r *http.Request) {
	    cryptoId := strings.SplitN(r.URL.Path, "/", 3)[2]
	    log.Println("Query " + r.URL.Path);
	    data, err := query(cryptoId)
	    if isError(err, w) {
			return
		}
		dataJson, err := json.Marshal(data)
		if isError(err, w) {
			return
		}
		renderResult(w, dataJson)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
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

func renderMenu(w http.ResponseWriter) {
	w.Write([]byte("<a href=\"/cryptodata/"+ID_BTC+"\">BTC price</a><br/>"))
	w.Write([]byte("<a href=\"/cryptodata/"+ID_ETH+"\">ETH price</a><br/>"))
}

func renderResult(w http.ResponseWriter, dataString []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(dataString)
	w.Write([]byte("<br/><br/><a href=\"/\">Back</a><br/>"))

}

func isError(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Panic(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
