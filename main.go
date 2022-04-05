package main

import (
	"log"
	"net/http"
	"permission/src/router"
)

//Ten projekt też będzie na MongoDB aby ogarnąć samemu o co chodzi w tworzeniu stron na GOLANG, wykorzystać GOROUTINES
//Następnie trzeba będzie stworzyć już manager z MySQL

func main() {
	//jako handler należy podać router aby można było wykorzystywać w ten sposób ścieżki
	err := http.ListenAndServe(":4000", router.Router())
	if err != nil {
		log.Fatal(err)
	}
}

//NAJPRAWDOPODOBNIEJ MUSZĘ:
//	Zainsanolać mongoDB na Laptopie (serwerze)
//	Zobaczyć jak dodawać użytkowników za pomocą golang
//  Stworzyć stronę z powiadomieniami DONE
