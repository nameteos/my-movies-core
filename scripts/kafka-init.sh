#!/bin/bash
set -e

# Wait for Kafka to be ready
echo "Waiting for Kafka to be ready..."
kafka-topics --bootstrap-server kafka:9092 --list

# Create topics for each domain event
echo "Creating Kafka topics..."

# Watchlist domain topics
echo "Creating watchlist domain topics..."
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic watchlist.movie_added --partitions 3 --replication-factor 1
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic watchlist.movie_removed --partitions 3 --replication-factor 1

# Library domain topics
echo "Creating library domain topics..."
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic library.movie_watched --partitions 3 --replication-factor 1

# Rating domain topics
echo "Creating rating domain topics..."
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic rating.movie_rated --partitions 3 --replication-factor 1
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic rating.movie_unrated --partitions 3 --replication-factor 1

# User domain topics
echo "Creating user domain topics..."
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic user.user_registered --partitions 3 --replication-factor 1
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic user.user_updated --partitions 3 --replication-factor 1
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic user.user_deleted --partitions 3 --replication-factor 1

# List all created topics
echo "Listing all topics:"
kafka-topics --bootstrap-server kafka:9092 --list

echo "Kafka topics created successfully"