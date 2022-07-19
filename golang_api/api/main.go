package main

import (
	"net/http"
	"fmt"
)

func main() {
	
	http.HandleFunc("/", HandleGet)

	http.ListenAndServe(":8080", nil)
}


func HandleGet(w http.ResponseWriter, r *http.Request){
	
	// videos := getVideos()

	// videoBytes, err  := json.Marshal(videos)

	// if err != nil {
  	// panic(err)
	// }
	fmt.Println("Endpoint Hit: homePage")
	//w.Write("all good!")
}