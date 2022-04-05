package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"permission/src/module"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey string = "AS24ias32NDp!@$#913[oLASDL~|!#@/ el/ij;sa>LPd0-;ohuahd8923huasdh9109283ihasdj10293ujojskabdi1gu2h3oi1o1ibdoasdj12pi312ndlas;do12p3h12bakdb3iukv4b23rir3j"

func Login(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email      string
		Password   string
		Genehmiger primitive.ObjectID
	}{}
	json.NewDecoder(r.Body).Decode(&data)
	var user module.User

	result := collection.FindOne(context.Background(), bson.M{"email": data.Email})
	err := result.Decode(&user)
	if err != nil {
		wg.Add(1)
		fail(w, r, err)
		wg.Wait()
	}
	wg.Add(1)
	go func() {
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    user.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		})
		e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
		if e == nil {
			token, err := claims.SignedString([]byte(SecretKey))
			if e != nil {
				panic(err)
			}

			cookie := http.Cookie{Name: "authentication", Value: token, HttpOnly: true}
			http.SetCookie(w, &cookie)

		} else {
			panic(e)
		}
		wg.Done()
	}()
	wg.Wait()
}
func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:4000/login", http.StatusSeeOther)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Location", "http://localhost:4000/login")
}

//checking user middleware

func IsLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utoken, err := r.Cookie("authentication")
		if err != nil {
			redirect(w, r)
		} else {
			token, err := jwt.Parse(utoken.Value, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(SecretKey), nil
			})
			if err != nil {
				redirect(w, r)
				// panic(err)
			}
			fmt.Println(token.Valid)
			if token.Valid {
				next(w, r)
			} else {
				redirect(w, r)
			}
		}
	}
}

//Robienie jakiś rzeczy na systemie po autoryzacji
