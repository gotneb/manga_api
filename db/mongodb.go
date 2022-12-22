package db

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/gotneb/manga_api/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// It represents a site where manga data will be scraped, and also collections avaliable
const (
	SEEMANGAS = iota
	MANGAINN
	READM
	MANGAS_CHAN
)

// Options
const (
	ALL = iota
	HIGHLIGHTS
	POPULARS
	RECENT_UPLOADS
)

var keyURI = os.Getenv("MONGODB_URI")

var databases = map[int]string{
	SEEMANGAS:   os.Getenv("SEEMANGAS_DATABASE"),
	MANGAS_CHAN: os.Getenv("MANGAS_CHAN_DATABASE"),
}

var collections = map[int]string{
	ALL:            os.Getenv("ALL_COLL"),
	HIGHLIGHTS:     os.Getenv("HIGHLIGHTS_COLL"),
	POPULARS:       os.Getenv("POPULARS_COLL"),
	RECENT_UPLOADS: os.Getenv("RECENT_UPLOADS_COLL"),
}

// This key will be used by me when uploading manga data using another language or tool (such as Python[Selenium])
var AuthUpload = os.Getenv("AUTH_KEY")

// Errors
var ErrMangaNotFound = errors.New("manga not found")
var ErrUnauthorized = errors.New("you're not unauthorized to use this service")

/*
I know it's a bad practice "repeat yourself", but I was too tired so,
I literally just copied and paste to get a *mango.Client in beginning of every function.
*/

func AddManga(server int, manga *web.Manga) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(keyURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(databases[server]).Collection(collections[ALL])
	_, err = coll.InsertOne(
		context.TODO(),
		bson.D{
			{"title", manga.Title},
			{"author", manga.Author},
			{"thumbnail", manga.Thumbnail},
			{"genres", manga.Genres},
			{"summary", manga.Summary},
			{"status", manga.Status},
			{"total_chapters", manga.TotalChapters},
			{"chapters", manga.Chapters},
		},
	)

	if err != nil {
		panic(err)
	}
	log.Println("OK | Added with sucess:", manga.Title)
}

// Returns a specified manga with the given title. E.g: "vinland saga"
func GetManga(server int, title string) (manga web.Manga, err error) {
	mangas, err := SearchManga(server, title)
	if len(mangas) >= 1 {
		manga = mangas[0]
	} else {
		err = ErrMangaNotFound
	}
	return
}

// Search all mangas with the given title
func SearchManga(server int, title string) (mangas []web.Manga, err error) {
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

	coll := client.Database(databases[server]).Collection(collections[ALL])

	model := mongo.IndexModel{Keys: bson.D{{"title", "text"}}}
	_, err = coll.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"$text", bson.D{{"$search", title}}}}
	sort := bson.D{{"score", bson.D{{"$meta", "textScore"}}}}
	opts := options.Find().SetSort(sort)
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		panic(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	if len(results) == 0 {
		err = ErrMangaNotFound
		return
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

// Update manga data on database
func UpdateManga(server int, manga *web.Manga) (err error) {
	// If do not exists
	_, err = FindManga(server, manga.Title)
	if err != nil {
		return
	}

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
	coll := client.Database(databases[server]).Collection(collections[ALL])

	// Update total_chapters
	filter := bson.D{{"title", manga.Title}}
	update := bson.D{{"$set", bson.D{{"total_chapters", manga.TotalChapters}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return
	}
	// Update array chapters
	_, err = coll.UpdateOne(
		context.TODO(),
		bson.D{{"title", manga.Title}},
		bson.D{{"$set", bson.D{{"chapters", manga.Chapters}}}},
	)
	if err != nil {
		return
	}

	log.Printf("%s | Updated!", manga.Title)
	return nil
}

// Finds by the exactly given manga title on the database
func FindManga(server int, title string) (manga web.Manga, err error) {
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

	coll := client.Database(databases[server]).Collection(collections[ALL])
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
