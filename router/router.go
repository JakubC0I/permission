package router

import (
	"github.com/JakubC0I/permission/src/github.com/JakubC0I/permission/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/statics/{id}", controller.Statics).Methods("GET")
	router.HandleFunc("/pfads", controller.ViewAll).Methods("GET")
	//Wszędzie gdzie "serve", jest to po prostu do wysłania html
	router.HandleFunc("/"+controller.M.AddUser, controller.ServeHTML).Methods("GET")
	router.HandleFunc("/"+controller.M.AddUser, controller.AddOneUser).Methods("POST")

	router.HandleFunc("/"+controller.M.AddPfad, controller.ServeHTML).Methods("GET")
	router.HandleFunc("/"+controller.M.AddPfad, controller.AddOnePfad).Methods("PUT")

	router.HandleFunc("/"+controller.M.RemovePfad, controller.ServeHTML).Methods("GET")
	router.HandleFunc("/"+controller.M.RemovePfad, controller.RemoveOnePfad).Methods("PUT")

	router.HandleFunc("/"+controller.M.Genehmigen+"/{id}", controller.Genehmigen).Methods("PUT")
	router.HandleFunc("/"+controller.M.Genehmigen+"/{id}", controller.ServeHTML).Methods("GET")

	return router
}
