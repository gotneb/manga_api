package db

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/gotneb/manga_api/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var keyURI = os.Getenv("MONGODB_URI")
var database = os.Getenv("DATABASE")
var collection = os.Getenv("COLLECTION")

/*
I know it's a bad practice "repeat yourself", but I was too tired so,
I literally just copied and paste to get a *mango.Client in beginning of in every function.
*/

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

	coll := client.Database(database).Collection(collection)
	_, err = coll.InsertOne(
		context.TODO(),
		bson.D{
			{"title", manga.Title},
			{"author", manga.Author},
			{"thumbnail", manga.Thumbnail},
			{"tags", manga.Tags},
			{"description", manga.Description},
			{"status", manga.Status},
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
func GetManga(title string) (manga web.Manga, err error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(keyURI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Disconnected!")
			panic(err)
		}
	}()

	coll := client.Database(database).Collection(collection)

	model := mongo.IndexModel{Keys: bson.D{{"title", "text"}}}
	_, err = coll.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		return
	}
	filter := bson.D{{"$text", bson.D{{"$search", title}}}}
	var result bson.M
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return
	}
	json.Unmarshal(jsonData, &manga)
	return
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
	manga, err := GetManga(name)
	if err != nil {
		return web.Chapter{}, err
	}
	c, _ := web.FetchImagesByName(manga.Title, val)
	return c, nil
}

// Search all mangas with the given title
func SearchManga(title string) (mangas []web.Manga, err error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(keyURI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Disconnected!")
			panic(err)
		}
	}()

	coll := client.Database(database).Collection(collection)

	model := mongo.IndexModel{Keys: bson.D{{"title", "text"}}}
	_, err = coll.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"$text", bson.D{{"$search", title}}}}
	cursor, err := coll.Find(context.TODO(), filter)
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for i := range results {
		var manga web.Manga
		jsonData, err := json.MarshalIndent(results[i], "", "    ")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(jsonData, &manga)
		mangas = append(mangas, manga)
	}
	return
}
