# nats-sandbox

### Run NATS server via Docker
---
```
docker run -p 4222:4222 -ti nats:latest
```

### Run NATS server with JetStream
```
docker run -p 4222:4222 -ti nats:latest -js
```