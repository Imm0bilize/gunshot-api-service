# gunshot-api-service


### Environment
```dotenv
# HTTP server
HTTP_PORT=8080

# Tracing
OTEL_HOST=localhost
OTEL_PORT=4317

# Database
DB_HOST=localhost
DB_PORT=27017
DB_NAME=GunShotService
DB_PASSWORD

# Kafka 
KAFKA_PEERS=localhost:9092
KAFKA_TOPIC=ApiServiceOutput
```

### TODO:
1. [x] use mongo
2. [ ] impl grpc and grpc stream
3. [ ] use auth (jwt token)
4. [ ] add swagger docs 
5. [ ] golangci-lint (configure CI/CD pipeline)
6. [ ] create pipeline (configure CI/CD pipeline)

