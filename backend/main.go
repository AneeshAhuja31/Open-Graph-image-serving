package main

import (
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	InitDB()
	PublicURL := os.Getenv("PUBLIC_URL")

	if PublicURL == "" {
		PublicURL = "http://localhost:8080"
	}
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
			log.Println(err)
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
		dc := gg.NewContext(1200, 630)

		dc.SetRGB(0.1, 0.1, 0.1)
		dc.Clear()

		if err := dc.LoadFontFace("fonts/InterVariable.ttf", 72); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dc.SetRGB(1, 1, 1)

		dc.DrawStringAnchored(
			user.Username,
			600,
			260,
			0.5,
			0.5,
		)

		if err := dc.LoadFontFace("fonts/InterVariable-Italic.ttf", 48); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dc.DrawStringAnchored(
			fmt.Sprintf("%d points", user.Points),
			600,
			340,
			0.5,
			0.5,
		)
		w.Header().Set("Content-Type","image/png")
		png.Encode(w,dc.Image())
	})

	mux.HandleFunc("GET /{username}",func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")
		collection := Client.Database("og_db").Collection("users")
		var user User

		err :=collection.FindOne(
			r.Context(),
			bson.M{"username":username},
		).Decode(&user)
		
		if err != nil{
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`"error":"User does not exist"`))
			return 
		}
		html := fmt.Sprintf(`
		<!doctype html>
		<html lang="en">
		<head>
			<meta charset="UTF-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1.0" />

			<meta property="og:title" content="%s - %d points">
			<meta property="og:description" content="%s has %d points">
			<meta property="og:image" content="%s/og/%s">

			<title>%s</title>

			<script type="module" crossorigin src="/assets/index-BOpV90vZ.js"></script>
		</head>
		<body>
			<div id="root"></div>
		</body>
		</html>
		`,
			user.Username, 
			user.Points,
			user.Username, 
			user.Points,
			PublicURL,     
			user.Username, 
			user.Username, 
		)
		w.Header().Set("Content-Type","text/html")
		w.Write([]byte(html))
	})
	mux.Handle(
		"/assets/",
		http.StripPrefix(
			"/assets/",
			http.FileServer(http.Dir("../frontend/dist/assets")),
		),
	)
	
	http.ListenAndServe(":8080",mux)
}