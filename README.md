# kafka-rest-proxy

```docker-compose up -d```

```
curl -X POST http://localhost:8080/topics/test-topic \
  -H "Content-Type: text/plain" \
  -d "Hello from REST API"
```