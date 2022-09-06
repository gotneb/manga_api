package db

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gotneb/manga_api/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const database = "manga_api"
const keyURI = "mongodb+srv://gotneb:D96jF02VEo5dyQK7@mangahoot-storage-512mb.qtc73bn.mongodb.net/?retryWrites=true&w=majority"

func AddManga(manga *web.Manga) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(keyURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("manga_api").Collection("meus_mangas")
	_, err = coll.InsertOne(
		context.TODO(),
		bson.D{
			{"title", manga.Title},
			{"author", manga.Author},
			{"thumbnail", manga.Thumbnail},
			{"tags", manga.Tags},
			{"description", manga.Description},
			{"situation", manga.Situation},
			{"total_chapters", manga.TotalChapters},
			{"chapters", manga.Chapters},
		},
	)

	if err != nil {
		panic(err)
	}
	log.Println("OK: Added with sucess:", manga.Title)
}

// Returns a specified manga with the given title
func SeachManga(title string) (manga web.Manga, err error) {
	log.Println("Searching into the database!")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(keyURI))
	if err != nil {
		panic(err)
	}
	log.Println("Connected to the database!")
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Disconnected!")
			panic(err)
		}
	}()

	coll := client.Database(database).Collection("meus_mangas")
	log.Println("Got connection!")

	model := mongo.IndexModel{Keys: bson.D{{"title", "text"}}}
	_, err = coll.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"$text", bson.D{{"$search", title}}}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(results[0], "", "    ")
	log.Println("Got JSON!")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(jsonData, &manga)
	return manga, nil
}

func AddChapter(ch *web.Chapter) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(keyURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(database).Collection("meus_mangas_chapters")
	_, err = coll.InsertOne(
		context.TODO(),
		bson.D{
			{"title", ch.Title},
			{"value", ch.Value},
			{"pages", ch.Pages},
		},
	)

	if err != nil {
		panic(err)
	}
	log.Println("OK: Added pages with sucess:", ch.Title)
}

func SearchChapter(name, val string) (web.Chapter, error) {
	manga, err := SeachManga(name)
	if err != nil {
		return web.Chapter{}, err
	}
	c, _ := web.FetchImagesByName(manga.Title, val)
	return c, nil
}
