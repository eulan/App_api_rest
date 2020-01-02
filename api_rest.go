package main

import (
	"./datapack"
	//"fmt"
	"net/http"
    "encoding/json"
    "github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)


func main() {
	
	router := chi.NewRouter() 
	router.Use(middleware.Recoverer)
	router.Use(cors.Default().Handler)


	router.Get("/", func (w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Bienvenido este es una API que hace cosas, usa \"/={dominio}\""))	})


	router.Get("/consults", func (w http.ResponseWriter, r *http.Request){
	
		var items_data datapack.Items	
		datapack.Builder_hist_db(&items_data)
		json.NewEncoder(w).Encode(items_data)
	})


	router.Get("/={domain}", func (w http.ResponseWriter, r *http.Request){
		domain := chi.URLParam(r, "domain")
		url := "https://api.ssllabs.com/api/v3/analyze?host=" + domain
		var server_data datapack.Server_data
		datapack.Data_obtain(domain, url, &server_data)
		json.NewEncoder(w).Encode(server_data)
	})
	
    
   	http.ListenAndServe(":8000", router)

	

    


}




