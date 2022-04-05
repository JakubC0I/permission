package router

import (
	"permission/src/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/statics/{id}", controller.Statics).Methods("GET")
	router.HandleFunc("/pfads", controller.ViewAll).Methods("GET")
	//Wszędzie gdzie "serve", jest to po prostu do wysłania html
	router.HandleFunc("/"+controller.M.AddUser, controller.ServeHTML).Methods("GET")
	router.HandleFunc("/"+controller.M.AddUser, controller.AddOneUser).Methods("POST")

	router.HandleFunc("/"+controller.M.AddPfad, controller.IsLoggedIn(controller.ServeHTML)).Methods("GET")
	router.HandleFunc("/"+controller.M.AddPfad, controller.IsLoggedIn(controller.AddOnePfad)).Methods("PUT")

	router.HandleFunc("/"+controller.M.RemovePfad, controller.ServeHTML).Methods("GET")
	router.HandleFunc("/"+controller.M.RemovePfad, controller.RemoveOnePfad).Methods("PUT")

	router.HandleFunc("/"+controller.M.Genehmigen+"/{id}", controller.IsLoggedIn(controller.ServeHTML)).Methods("GET")
	router.HandleFunc("/"+controller.M.Genehmigen+"/{id}", controller.IsLoggedIn(controller.Genehmigen)).Methods("PUT")

	router.HandleFunc("/"+controller.M.Deny+"/{id}", controller.IsLoggedIn(controller.Deny)).Methods("DELETE")

	router.HandleFunc("/"+controller.M.AddImage, controller.ServeHTML).Methods("GET")
	router.HandleFunc("/"+controller.M.AddImage, controller.AddImage).Methods("POST")

	router.HandleFunc("/"+controller.M.Ticket+"/{id}", controller.IsLoggedIn(controller.ServeHTMLid)).Methods("GET")
	router.HandleFunc("/search", controller.LiveSearch).Methods("POST")

	router.HandleFunc("/addComment", controller.AddCommentTicket).Methods("POST")

	router.HandleFunc("/"+controller.M.Login, controller.ServeHTML).Methods("GET")
	router.HandleFunc("/"+controller.M.Login, controller.Login).Methods("POST")

	router.HandleFunc("/"+controller.M.Secret, controller.IsLoggedIn(controller.ServeHTML)).Methods("GET")

	router.HandleFunc("/"+controller.M.Notification, controller.IsLoggedIn(controller.Notification)).Methods("GET")
	// router.HandleFunc("/createUserSystem", controller.AddUserSystem).Methods("POST")
	return router
}
