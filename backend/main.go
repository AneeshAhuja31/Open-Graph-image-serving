package main

import (
	"encoding/json"
	"image/png"

	"fmt"
	"net/http"

	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {

	InitDB()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{username}",func(w http.ResponseWriter, r *http.Request) {
		username:= r.PathValue("username")
		collection := Client.Database("og_db").Collection("users")
		ctx := r.Context()
		var user User

		err := collection.FindOne(
			ctx,
			bson.M{"username":username},
		).Decode(&user)
		w.Header().Set("Content-Type","application/json")
		if err != nil{
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"User does not exist"}`))
			return
		}
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("GET /og/{username}",func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")
		collection := Client.Database("og_db").Collection("users")
		ctx := r.Context()
		var user User

		err := collection.FindOne(
			ctx,
			bson.M{"username":username},
		).Decode(&user)
		if err!=nil{
			user.Username = username
			user.Points=-1
		}
		dc := gg.NewContext(1200,630)
		dc.SetRGB(0.1,0.1,0.1)
		dc.Clear()
		if err := dc.LoadFontFace("fonts/InterVariable.ttf", 72); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dc.SetRGB(1, 1, 1)
		dc.DrawString(user.Username, 80, 180)
		dc.LoadFontFace("fonts/InterVariable-Italic.ttf", 48)
		dc.DrawString(
			fmt.Sprintf("%d points", user.Points),
			80,
			280,
		)
		w.Header().Set("Content-Type","image/png")
		png.Encode(w,dc.Image())
	})
	http.ListenAndServe(":8080",mux)
}