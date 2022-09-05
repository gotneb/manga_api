package db

import (
	"context"
	"log"

	"github.com/gotneb/manga_api/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
