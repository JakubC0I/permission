package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"permission/src/module"
	"regexp"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Implementing CRUD from helper and serving templates.
var wg sync.WaitGroup

var M module.Variables = module.Variables{
	AddPfad:      "addPfad",
	RemovePfad:   "removePfad",
	AddUser:      "addUser",
	Genehmigen:   "genehmigen",
	AddImage:     "addImage",
	Ticket:       "ticket",
	Login:        "login",
	Register:     "register",
	Secret:       "secret",
	Notification: "notification",
	Deny:         "deny",
}

func ViewAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Request-Method", "GET")
	pfads, err := GetAllPfads()
	if err != nil {
		wg.Add(1)
		go fail(w, r, err)
		wg.Wait()
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
		wg.Add(1)
		go fail(w, r, err)
		wg.Wait()
	}
	//Sprawdź czy taki użytkownik istnieje DODAĆ JESZCZE SPRAWDZANIE PERMISSON LEVEL
	e := collection.FindOne(context.Background(), bson.M{"_id": user.Genehmiger}).Decode(&genehmiger)
	if e != nil {
		wg.Add(1)
		go fail(w, r, e)
		wg.Wait()
	} else {
		result, err := AddUser(user)
		userid := result.InsertedID
		if err != nil {
			wg.Add(1)
			go fail(w, r, err)
			wg.Wait()
		} else {
			wg.Add(2)
			go AddUserSystem(userid.(primitive.ObjectID).Hex())
			go json.NewEncoder(w).Encode(module.Success{Success: true, Message: "Item added successfully", Info: result.InsertedID.(primitive.ObjectID)})
			wg.Done()
			wg.Wait()
		}

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
		wg.Add(1)
		go fail(w, r, err)
		wg.Wait()
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
			wg.Add(1)
			go fail(w, r, err)
			wg.Wait()
		}
		if result.ModifiedCount == 0 {
			wg.Add(2)
			go func() {
				json.NewEncoder(w).Encode(module.Success{Success: false, Message: "Items could not be added or pfad already there"})
				wg.Done()
			}()
			go func() {

				tasks.DeleteOne(context.Background(), bson.M{"_id": task.ID})
				wg.Done()
			}()
			wg.Wait()
		} else {
			//Add user - potem usunac
			// fakeUser, _ := primitive.ObjectIDFromHex("62323c5708234554ed1df393")
			// AddUserSystem(fakeUser)
			//CHANGE GROUP
			wg.Add(3)
			go ChangeGroup(task.Data, task.Betroffene)

			go func() {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(module.Success{Success: true, Message: "Items added successfully", Info: task.ID})
				wg.Done()
			}()
			go tasks.DeleteOne(context.Background(), bson.M{"_id": task.ID})
			wg.Done()

			wg.Wait()
		}
		fmt.Printf("Numer of items with added pfad: %v\n", result.ModifiedCount)
	} else {
		task.Action = r.URL.Path
		task.Created_at = time.Now()
		var user module.User
		cook, _ := r.Cookie("authentication")
		cookieID := IDCookie(cook)
		err := collection.FindOne(context.Background(), bson.M{"_id": cookieID}).Decode(&user)
		if err != nil {
			wg.Add(1)
			go fail(w, r, err)
			wg.Wait()
		} else {
			task.Genehmiger = user.Genehmiger
			task.Besteller = cookieID
			wg.Add(1)
			go tasks.InsertOne(context.Background(), &task)
			wg.Done()
			wg.Wait()
		}
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
		wg.Add(1)
		go fail(w, r, err)
		wg.Wait()
	}
	task.Genehmigt = true
	//https://riptutorial.com/go/example/27703/put-request-of-json-object
	accepted, e := json.Marshal(&task)
	if e != nil {
		wg.Add(1)
		go fail(w, r, err)
		wg.Wait()
	}
	var client = &http.Client{}
	req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:4000"+task.Action, bytes.NewBuffer(accepted))
	if err != nil {
		wg.Add(1)
		go fail(w, r, err)
		wg.Wait()
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	cookie, _ := r.Cookie("authentication")
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		wg.Add(1)
		go fail(w, r, err)
		wg.Wait()
	}
	var response module.Success
	json.NewDecoder(resp.Body).Decode(&response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//Dodawanie zdjęć oraz opisu
func AddImage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var images = make(chan string)
	//TRZEBA BĘDZIE JESZE DODAĆ KTO STWORZYŁ

	m := struct {
		Images      []string `json:"image"`
		Description string   `json:"description"`
		Title       string   `json:"title"`
		Article     int      `json:"article"`
		Status      string   `json:"status"`
		Created_at  time.Time
	}{Status: "Queued", Created_at: time.Now(), Article: 0, Title: "No Informations added please remove this ticket"}
	json.NewDecoder(r.Body).Decode(&m)

	wg.Add(len(m.Images))

	go func() {
		wg.Wait()
		close(images)
	}()
	for i := 0; i < len(m.Images); i++ {
		reg := regexp.MustCompile(",.*")
		str := reg.Find([]byte(m.Images[i]))
		s := str[0:]
		payload := &bytes.Buffer{}
		go addTicketImage(payload, s, images)
	}
	//POWÓD dla komentarza: https://dev.to/sophiedebenedetto/synchronizing-go-routines-with-channels-and-waitgroups-3ke2
	// wg.Wait()
	imgs := struct {
		Images   []string    `json:"images"`
		Incident interface{} `json:"incident"`
	}{}

	for i := range images {
		imgs.Images = append(imgs.Images, i)
	}
	// fmt.Println(imgs)
	wg.Add(1)
	go func() {
		result, _ := incident.InsertOne(context.Background(), bson.D{
			{Key: "description", Value: m.Description},
			{Key: "images", Value: &imgs.Images},
			{Key: "createdAt", Value: m.Created_at},
			{Key: "status", Value: m.Status},
			{Key: "title", Value: m.Title},
			{Key: "article", Value: m.Article},
		})
		iID := result.InsertedID
		imgs.Incident = iID
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&imgs)

		wg.Done()
	}()
	wg.Wait()
}

//Searchbar
func addTicketImage(payload *bytes.Buffer, s []byte, images chan<- string) {
	client := &http.Client{}

	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("image", string(s))
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.imgur.com/3/image", payload)
	if err != nil {
		fmt.Println("Forming request")
		panic(err)
	}
	req.Header.Add("Authorization", "Client-ID 5dc2c8848f66984")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Sending to imgur")
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	type Body struct {
		Data struct {
			Link string `json:"link"`
		} `json:"data"`
	}
	var data Body
	json.Unmarshal(body, &data)
	images <- data.Data.Link
	// fmt.Println(string(body))
	wg.Done()
}
func ServeHTMLid(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	t := noExecute(w, r, p)
	params := mux.Vars(r)
	oID, _ := primitive.ObjectIDFromHex(params["id"])
	result := incident.FindOne(context.Background(), bson.M{"_id": oID})
	type MID struct {
		module.Variables
		ID          string
		Images      []string `json:"image"`
		Description string   `json:"description"`
		Title       string   `json:"title"`
		Article     int      `json:"article"`
		Comments    []string `json:"comments"`
	}
	mid := MID{M, params["id"], nil, "No records", "No Title", 0, nil}
	result.Decode(&mid)

	t.Execute(w, &mid)
	fmt.Println(&mid)
}

func LiveSearch(w http.ResponseWriter, r *http.Request) {
	search := struct {
		Searchbar string `json:"searchbar"`
	}{}
	json.NewDecoder(r.Body).Decode(&search)
	fmt.Println(&search)
	if len(search.Searchbar) >= 3 {
		cursor, err := incident.Find(context.Background(), bson.M{"$or": bson.A{
			bson.M{"title": bson.M{"$regex": ".*" + search.Searchbar + ".*", "$options": "i"}},
			bson.M{"description": bson.M{"$regex": ".*" + search.Searchbar + ".*", "$options": "i"}},
		}})
		if err != nil {
			panic(err)
		}
		var incs [][]string
		for cursor.Next(context.Background()) {
			data, err := cursor.Current.Elements()
			incs = append(incs, []string{data[1].Value().StringValue(), data[5].Value().StringValue(), data[0].Value().ObjectID().Hex()})
			if err != nil {
				panic(err)
			}
		}
		inc := struct {
			Data [][]string `json:"data"`
		}{incs}
		defer cursor.Close(context.Background())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inc)
	} else {
		success := module.Success{Success: false, Message: "At least 3 characters needed to look for a result"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&success)
	}

}

func AddCommentTicket(w http.ResponseWriter, r *http.Request) {
	x := r.Referer()
	id := []byte(x)
	bID, _ := primitive.ObjectIDFromHex(string(id[29:53]))
	com := struct {
		Status  string
		Comment string
	}{Status: "Check replay"}
	json.NewDecoder(r.Body).Decode(&com)
	res, err := incident.UpdateOne(context.Background(), bson.M{"_id": bID}, bson.M{"$set": bson.M{"status": &com.Status}, "$addToSet": bson.M{"comments": &com.Comment}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.ModifiedCount, &com)
}

//Do strony na której będą wyświetlane wszystkie możliwe taski do zaakceptowania
func Notification(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	t := noExecute(w, r, p)
	decode := struct {
		module.Variables
		Tasks [][]interface{}
	}{M, [][]interface{}{}}
	// var decode ToSite
	// decode.Variable = M
	// decode.Tasks = append(decode.Tasks, []string{"string"})
	cook, _ := r.Cookie("authentication")
	id := IDCookie(cook)
	fmt.Println(id)
	tasksArr, err := tasks.Find(context.Background(), bson.M{"genehmiger": id})
	i := 0
	for tasksArr.Next(context.Background()) {
		ele, err := tasksArr.Current.Elements()
		if err != nil {
			panic(err)
		}
		// var find bson.M
		// e := tasksArr.Decode(&find)
		// if e != nil {
		// 	fmt.Errorf("Error %v", e)
		// }
		// finds = append(finds, find)
		oID := ele[0].Value().ObjectID().Hex()
		oBetro := ele[1].Value().ObjectID().Hex()
		oAction := ele[2].Value().String()
		oData, _ := ele[6].Value().Array().Elements()
		decode.Tasks = append(decode.Tasks, []interface{}{oID, oAction, oData, oBetro})
		i++
	}
	if err != nil {
		fmt.Errorf("error occured %v", err)
	}
	t.Execute(w, decode)
}

func Deny(w http.ResponseWriter, r *http.Request) {
	//zaprogramować usuwanie taska z bazy danych
	t := struct {
		TaskID string
	}{}
	var task module.Task
	cook, err := r.Cookie("authentication")
	if err != nil {
		panic(err)
	}
	genehmiger := IDCookie(cook)
	json.NewDecoder(r.Body).Decode(&t)
	params := mux.Vars(r)
	oID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		panic(err)
	}
	result := tasks.FindOne(context.Background(), bson.M{"_id": oID})

	e := result.Decode(&task)
	if task.Genehmiger.Hex() == genehmiger.Hex() {
		_, err = tasks.DeleteOne(context.Background(), bson.M{"_id": oID})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bson.M{"success": true, "info": "Task denied", "_id": task.ID})
		}
	} else {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}

}

//Crypto and creating jwt (authentication)
