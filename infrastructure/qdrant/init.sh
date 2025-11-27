#!/bin/bash
# Qdrant initialization script
# This script creates the necessary collections in Qdrant

echo "Initializing Qdrant collections..."

# Wait for Qdrant to be ready
until curl -f http://localhost:6333/healthz; do
  echo "Waiting for Qdrant to be ready..."
  sleep 2
done

echo "Qdrant is ready. Creating collections..."

# Create api-knowledge collection
curl -X PUT "http://localhost:6333/collections/api-knowledge" \
  -H "Content-Type: application/json" \
  -d '{
    "vectors": {
      "size": 384,
      "distance": "Cosine"
    }
  }'

echo "Created api-knowledge collection"

# Create learned-patterns collection
curl -X PUT "http://localhost:6333/collections/learned-patterns" \
  -H "Content-Type: application/json" \
  -d '{
    "vectors": {
      "size": 384,
      "distance": "Cosine"
    }
  }'

echo "Created learned-patterns collection"

echo "Qdrant initialization complete!"

