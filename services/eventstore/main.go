package eventstore

// eventstore is a NATS to MongoDB event store
// It listens for events on a NATS subject and writes them to a MongoDB collection
// It uses a batch writer to write events in batches for efficiency
// The batch size and timeout can be configured via command line flags or config file
// The config file can be specified using the --config flag
// The default config file is config.yaml in the current directory
// The NATS subject to listen to can be specified using the --subject flag
// The default subject is "event.>"
// The MongoDB URI and database can be specified in the config file
// The default MongoDB URI is "mongodb://localhost:27017"
// The default MongoDB database is "events"
// The batch size and timeout can be specified using the --batch-size and --batch-timeout flags
// The default batch size is 1000
// The default batch timeout is 5 seconds
// The program will listen for events on the specified NATS subject and write them to the specified MongoDB collection
// It will exit gracefully on SIGINT or SIGTERM

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julianshen/nebula/data"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rootCmd = &cobra.Command{
	Use:   "eventstore",
	Short: "NATS to MongoDB event store",
	Run:   run,
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().Int("batch-size", 1000, "batch size for MongoDB writes")
	rootCmd.PersistentFlags().Duration("batch-timeout", 5*time.Second, "max time to wait before flushing batch")
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) {
	// Initialize MongoDB client
	mongoURI := viper.GetString("mongo.uri")
	if mongoURI == "" {
		log.Fatal("MongoDB URI not configured")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Initialize NATS connection
	natsURL := viper.GetString("nats.url")
	if natsURL == "" {
		log.Fatal("NATS URL not configured")
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Create buffered channel for events
	batchSize := viper.GetInt("batch-size")
	batchTimeout := viper.GetDuration("batch-timeout")
	events := make(chan data.Event, batchSize*2)

	// Start batch writer
	go batchWriter(client, events, batchSize, batchTimeout)

	// Subscribe to NATS events matching pattern
	subject := viper.GetString("nats.subject")
	if subject == "" {
		subject = "event.>"
	}

	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		var event data.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			return
		}
		events <- event
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening for events on:", subject)
	waitForInterrupt()
}

func batchWriter(client *mongo.Client, events <-chan data.Event, batchSize int, batchTimeout time.Duration) {
	coll := client.Database(viper.GetString("mongo.database")).Collection("events")
	ticker := time.NewTicker(batchTimeout)
	defer ticker.Stop()

	var batch []interface{}
	for {
		select {
		case event, ok := <-events:
			if !ok {
				if len(batch) > 0 {
					flushBatch(coll, batch)
				}
				return
			}
			batch = append(batch, event)
			if len(batch) >= batchSize {
				flushBatch(coll, batch)
				batch = nil
				ticker.Reset(batchTimeout)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				flushBatch(coll, batch)
				batch = nil
			}
		}
	}
}

func flushBatch(coll *mongo.Collection, batch []interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := coll.InsertMany(ctx, batch)
	if err != nil {
		log.Printf("Error writing batch: %v", err)
		return
	}
	log.Printf("Wrote batch of %d events", len(batch))
}

func waitForInterrupt() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
