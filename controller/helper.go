package controller

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"permission/src/module"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

//connect to MongoDb

const connectionString = "mongodb+srv://jakub_user:ck3YJrce9rtuPRdj@cluster0.vdtkf.mongodb.net/GolangDatabase?retryWrites=true&w=majority"

// const connectionString = "mongodb://jakub_user:jakub123@127.0.0.1:27017"
const dbName = "GolangDatabase"
const collName = "users"

//MOST IMPORTANT
var collection *mongo.Collection
var tasks *mongo.Collection
var incident *mongo.Collection

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
	incident = client.Database(dbName).Collection("incidents")
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
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		panic(err)
	}
	user.Password = string(hash)
	user.Date = time.Now()
	user.Pfads = []string{}
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		panic(err)
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
	p := r.URL.Path
	t := noExecute(w, r, p)
	t.Execute(w, M)
}

func noExecute(w http.ResponseWriter, r *http.Request, p string) *template.Template {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Request-Method", "GET")
	// fmt.Println(filepath.Dir("../templates"))
	dirname, _ := os.Getwd()
	fmt.Print(p)
	re := regexp.MustCompile("^/.*/")
	path := string(re.Find([]byte(p)))
	if path == "" {
		path = r.URL.Path
	}
	if last := len(path) - 1; last >= 0 && path[last] == '/' {
		path = path[:last]
	}
	t, err := template.ParseFiles(dirname + "/templates" + path + ".html")
	if err != nil {
		log.Fatal(err)
	}
	// val, _ := ioutil.ReadFile("./templates/partials/footer.html")
	// M.Templates.Footer = string(val)
	cook, err := r.Cookie("authentication")
	if err != nil {
		M.UserID = ""
	} else {
		cookieID := IDCookie(cook).Hex()
		M.UserID = cookieID
	}
	t.ParseGlob("./templates/partials/*")
	return t
}
func Statics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dirname, _ := os.Getwd()
	data, err := ioutil.ReadFile(dirname + "/templates/statics/" + params["id"])
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(data)
}

func fail(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	m := struct {
		Success bool
		Info    string
	}{
		Success: false, Info: err.Error(),
	}
	json.NewEncoder(w).Encode(m)
	fmt.Println(err)
	wg.Done()
}

func ChangeGroup(pfads []string, ids []primitive.ObjectID) {
	for _, key := range ids {
		user := key.Hex()
		for _, k := range pfads {
			res := exec.Command("setfacl", "-Rm", fmt.Sprintf("u:%v:rw", user), k)
			res.Output()
		}
	}
	wg.Done()
}

func AddUserSystem(id string) {
	// id := struct {
	// 	UserID string `json:"userid" bson:"userid"`
	// }{}
	// json.NewDecoder(r.Body).Decode(&id)
	//Jednak nie będę tworzyć po mailu
	// comp, _ := regexp.Compile(".*@")
	// s := []byte(email)
	// id := string(comp.Find(s))
	// if last := len(id) - 1; last >= 0 && id[last] == '@' {
	// 	id = id[:last]
	// }
	agrsUser := []string{"-m", id}
	re := exec.Command("useradd", agrsUser...)
	res, err := re.Output()
	//Po dodaniu comendy trzeba jeszcze Run albo Output!!
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
	wg.Done()
}
func IDCookie(cook *http.Cookie) primitive.ObjectID {
	parts := strings.Split(cook.Value, ".")
	base, _ := base64.RawStdEncoding.DecodeString(parts[1])
	reader := bytes.NewReader(base)
	cookie := struct {
		Exp int64  `json:"exp"`
		Iss string `json:"iss"`
	}{}
	json.NewDecoder(reader).Decode(&cookie)

	cookieID, _ := primitive.ObjectIDFromHex(cookie.Iss)
	return cookieID
}
