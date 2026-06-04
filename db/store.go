package db

import (
	"context"
	"errors"
	"paul/scorist/models"

	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const guildsCollection = "guilds"

type Store struct {

	client *mongo.Client
	coll *mongo.Collection

}

func Connect(ctx context.Context, uri string) (*Store, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {

		return nil, err

	}

	if err = client.Ping(ctx, nil); err != nil {

		return nil, err

	}

	return &Store{

		client: client,
		coll: client.Database("scorist").Collection(guildsCollection),

	}, nil

}

func (s *Store) Close(ctx context.Context) error {

	return s.client.Disconnect(ctx)

}

func (s *Store) Get(ctx context.Context, guildID snowflake.ID) (*models.Guild, error) {

	var guild models.Guild

	err := s.coll.FindOne(ctx, bson.M{"_id": guildID}).Decode(&guild)

	if errors.Is(err, mongo.ErrNoDocuments) {

		return models.NewGuild(guildID), nil

	}

	if err != nil {

		return nil, err

	}

	return &guild, nil

}

func (s *Store) Save(ctx context.Context, guild *models.Guild) error {

	_, err := s.coll.ReplaceOne(ctx, bson.M{"_id": guild.ID}, guild, options.Replace().SetUpsert(true))

	return err

}

func (s *Store) ListConfigured(ctx context.Context) ([]*models.Guild, error) {

	cursor, err := s.coll.Find(ctx, bson.M{"channels.updates": bson.M{"$ne": 0}})

	if err != nil {

		return nil, err

	}

	defer cursor.Close(ctx)

	var guilds []*models.Guild

	for cursor.Next(ctx) {

		var guild models.Guild

		if err := cursor.Decode(&guild); err != nil {

			return nil, err

		}

		guilds = append(guilds, &guild)

	}

	return guilds, cursor.Err()

}