package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/JakubC0I/permission/src/github.com/JakubC0I/permission/module"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//connect to MongoDb

// const connectionString = "mongodb+srv://jakub_user:ck3YJrce9rtuPRdj@cluster0.vdtkf.mongodb.net/GolangDatabase?retryWrites=true&w=majority"
const connectionString = "mongodb://jakub_user:jakub123@127.0.0.1:27017"
const dbName = "GolangDatabase"
const collName = "users"

//MOST IMPORTANT
var collection *mongo.Collection
var tasks *mongo.Collection

func init() {
	//options
	option := options.Client().ApplyURI(connectionString)
	//connecting to the client
	client, err := mongo.Connect(context.TODO(), option)
	if err != nil {
		log.Fatal(err)
	}
	//connecting to the database
	collection = client.Database(dbName).Collection(collName)
	tasks = client.Database(dbName).Collection("tasks")
}

//Basic CRUD operations
func GetAllPfads() ([]primitive.D, error) {
	var records []primitive.D
	//cursor to duży plik, przesyłany jak buffer. Dlatego należy używać cursor.Next()
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	//ZAPAMIĘTAĆ TEGO LOOPA
	for cursor.Next(context.Background()) {
		var record bson.D
		err := cursor.Decode(&record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func AddUser(user module.User) (*mongo.InsertOneResult, error) {
	user.Date = time.Now()
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func RemovePfad(id string, ids []primitive.M) (*mongo.DeleteResult, error) {
	valid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	// collection.UpdateMany(context.Background(), )
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": valid})
	return result, err
}
func ServeHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Request-Method", "GET")
	// fmt.Println(filepath.Dir("../templates"))
	p := r.URL.Path
	fmt.Print(p)
	re := regexp.MustCompile("^/.*/")
	path := string(re.Find([]byte(p)))
	if path == "" {
		path = r.URL.Path
	}
	if last := len(path) - 1; last >= 0 && path[last] == '/' {
		path = path[:last]
	}
	t, err := template.ParseFiles("/home/lodian100/Informatyka/goLang/src/github.com/JakubC0I/permission/templates/" + path + ".html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, M)
	if err != nil {
		log.Fatal(err)
	}
}

func Statics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	data, err := ioutil.ReadFile("/home/lodian100/Informatyka/goLang/src/github.com/JakubC0I/permission/templates/statics/" + params["id"])
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(data)
}

func fail(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bson.M{"Success": false, "Info": "Could not accept"})
	w.Header().Set("Content-Type", "application/json")
}
