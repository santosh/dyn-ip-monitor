package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IPInfo struct {
	IP string `json:"ip"`
}

func main() {
	// Attempt to load .env file
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Database connection parameters
	uri := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Println("Error disconnecting from the database:", err)
		}
	}()

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return
	}

	collection := client.Database("dynip").Collection("ip_log")

	// Default sleep time is 900 seconds (15 minutes) if not set
	interval := 900
	if val, ok := os.LookupEnv("INTERVAL"); ok {
		if parsedInterval, err := strconv.Atoi(val); err == nil {
			interval = parsedInterval
		}
	}

	for {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		ip := getPublicIP()
		if ip != "" {
			fmt.Printf("%s - IP: %s\n", timestamp, ip)
			saveToDatabase(collection, timestamp, ip)
		} else {
			fmt.Println("Failed to retrieve IP address")
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func getPublicIP() string {
	resp, err := http.Get("http://ipinfo.io")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return ""
	}

	return ipInfo.IP
}

func saveToDatabase(collection *mongo.Collection, timestamp, ip string) {
	doc := bson.M{"timestamp": timestamp, "ip": ip}
	_, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Println("Error inserting into database:", err)
	}
}
