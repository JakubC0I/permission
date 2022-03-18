package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"

	"github.com/JakubC0I/permission/src/github.com/JakubC0I/permission/module"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Implementing CRUD from helper and serving templates.

var M struct {
	AddPfad    string
	RemovePfad string
	AddUser    string
	Genehmigen string
} = module.Variables{
	AddPfad:    "addPfad",
	RemovePfad: "removePfad",
	AddUser:    "addUser",
	Genehmigen: "genehmigen",
}

func ViewAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Request-Method", "GET")
	pfads, err := GetAllPfads()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(pfads)
}

func AddOneUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Request-Method", "POST")
	var user module.User
	var genehmiger module.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	//Sprawdź czy taki użytkownik istnieje DODAĆ JESZCZE SPRAWDZANIE PERMISSON LEVEL
	e := collection.FindOne(context.Background(), bson.M{"_id": user.Genehmiger}).Decode(&genehmiger)
	if e != nil {
		log.Fatal(e)
	} else {
		result, err := AddUser(user)
		if err != nil {
			log.Fatal(err)
			json.NewEncoder(w).Encode(module.Success{Success: false, Message: "Item could not be added"})
		}
		json.NewEncoder(w).Encode(module.Success{Success: true, Message: "Item added successfully", Info: result.InsertedID.(primitive.ObjectID)})

	}

}

//PRZY REMOVE I ADD PFAD NALEZY JESZCZE DODAC SYSTEM GENEHMIGOWANIA

func RemoveOnePfad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Request-Method", "PUT")
	var user module.UserArray
	json.NewDecoder(r.Body).Decode(&user)
	// cursor, err := collection.Find(context.Background(), bson.M{"_id": bson.M{"$in": user.IDS}})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var finds []primitive.M
	// for cursor.Next(context.Background()) {
	// 	var find bson.M
	// 	err := cursor.Decode(&find)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	finds := append(finds, find)

	// }
	result, err := collection.UpdateMany(context.Background(), bson.M{"_id": bson.M{"$in": user.IDS}}, bson.M{"$pull": bson.M{"pfads": bson.M{"$in": user.Pfads}}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of items with deleted pfad: %v\n", result.ModifiedCount)
}

func AddOnePfad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Request-Method", "PUT")
	var task module.Task
	json.NewDecoder(r.Body).Decode(&task)
	if task.Genehmigt {
		result, err := collection.UpdateMany(context.Background(), bson.M{"_id": bson.M{"$in": task.Betroffene}}, bson.M{"$addToSet": bson.M{"pfads": bson.M{"$each": task.Data}}})
		if err != nil {
			fail(w, r)
			panic(err)
		}
		if result.ModifiedCount == 0 {
			json.NewEncoder(w).Encode(module.Success{Success: false, Message: "Items could not be added"})
		} else {
			//CHANGE GROUP
			ChangeGroup(task.Data, task.Betroffene)

			json.NewEncoder(w).Encode(module.Success{Success: true, Message: "Items added successfully"})
			tasks.DeleteOne(context.Background(), bson.M{"_id": task.ID})
		}
		fmt.Printf("Numer of items with added pfad: %v\n", result.ModifiedCount)
	} else {
		task.Action = r.URL.Path
		task.Created_at = time.Now()
		var user module.User
		collection.FindOne(context.Background(), bson.M{"_id": task.Besteller}).Decode(&user)
		// fmt.Println(user)
		task.Genehmiger = user.Genehmiger
		tasks.InsertOne(context.Background(), &task)
	}

}

func Genehmigen(w http.ResponseWriter, r *http.Request) {
	//Find w tasks zmienic na true, wyslac to na tego samego pfada co w action
	//powinno się zmienić w user i usunąć taska
	params := mux.Vars(r)
	var task module.Task
	id, _ := primitive.ObjectIDFromHex(params["id"])
	err := tasks.FindOne(context.Background(), bson.M{"_id": id}).Decode(&task)
	if err != nil {
		fail(w, r)
		panic(err)
	}
	task.Genehmigt = true
	//https://riptutorial.com/go/example/27703/put-request-of-json-object
	accepted, e := json.Marshal(task)
	if e != nil {
		fail(w, r)
		log.Fatal(e)
	}
	var client = &http.Client{}
	req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:4000"+task.Action, bytes.NewBuffer(accepted))
	if err != nil {
		fail(w, r)
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		fail(w, r)
		panic(err)
	}
	fmt.Print(resp.StatusCode)
}

func ChangeGroup(pfads []string, ids []primitive.ObjectID) {
	// for v, k = range pfads {
	exec.Command("setfacl")
	// }
}

func AddUserSystem(id primitive.ObjectID, email string) {
	comp, _ := regexp.Compile(".*@")
	s := []byte(email)
	str := string(comp.Find(s))
	if last := len(str) - 1; last >= 0 && str[last] == '@' {
		str = str[:last]
	}
	exec.Command("useradd")

	fmt.Print(str)
}

//Crypto and creating jwt (authentication)
