package migration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Migration handles the process of importing JSON data files into MongoDB collections.
// It maps file names to collection names, manages the MongoDB client connection,
// specifies the directory containing data files, and holds the target database name.
//
// Fields:
//   - FileToCollection: A map where the key is the JSON file name and the value is the corresponding MongoDB collection name.
//   - Client:           The MongoDB client used to connect and perform operations on the database.
//   - DataDir:          The directory path where the JSON data files are stored.
//   - DatabaseName:     The name of the MongoDB database where collections reside.
type Migration struct {
	FileToCollection map[string]string
	Client           *mongo.Client
	DataDir          string
	DatabaseName     string
}

// NewMigration creates a new Migration instance.
func NewMigration(client *mongo.Client, fileToCollection map[string]string, dir string) *Migration {
	return &Migration{
		FileToCollection: fileToCollection,
		Client:           client,
		DataDir:          dir,
	}
}

// closeFileWithLog safely closes a file and logs any errors using zerolog.
func closeFileWithLog(file *os.File, filePath string) {
	if err := file.Close(); err != nil {
		log.Warn().Err(err).Str("file", filePath).Msg("Failed to close file")
	}
}

// Up imports JSON data files into their corresponding MongoDB collections.
// For each file-to-collection mapping, it checks if the collection already contains data.
// If the collection is empty, it reads the JSON file and inserts its contents into the collection.
// Skips collections that already have documents to avoid duplicate imports.
func (m *Migration) Up(ctx context.Context) error {
	dataDir := m.DataDir
	database := m.Client.Database(m.DatabaseName)

	for fileName, collectionName := range m.FileToCollection {
		collection := database.Collection(collectionName)

		// Check if collection already contains documents
		count, err := collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("failed to count documents in collection %s: %w", collectionName, err)
		}
		if count > 0 {
			log.Info().Str("collection", collectionName).Msg("skipping collection, data already exists in collections")
			continue // Skip if documents already exist
		}

		filePath := filepath.Join(dataDir, fileName)

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", filePath, err)
		}

		content, err := io.ReadAll(file)
		closeFileWithLog(file, filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		// Decode JSON into slice of maps
		var documents []any
		if err := json.Unmarshal(content, &documents); err != nil {
			return fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
		}

		// Insert documents if any
		if len(documents) > 0 {
			_, err = collection.InsertMany(ctx, documents, options.InsertMany().SetOrdered(false))
			if err != nil {
				return fmt.Errorf("failed to insert documents into collection %s: %w", collectionName, err)
			}
		}
	}

	return nil
}
