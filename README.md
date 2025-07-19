# Toll Calculator

A distributed microservices system for calculating toll invoices based on vehicle movements. The system consists of three services — Receiver, Distance Calculator, and Aggregator — that communicate asynchronously using Apache Kafka. It supports both gRPC and HTTP clients to send and retrieve data.

---

## Overview

This project simulates a toll collection architecture implemented in Go. It processes vehicle passage events, calculates distances between entry and exit points, and maintains a real-time invoice per vehicle (OBU ID).

---

## Services

### 1. Receiver Service

- Accepts raw toll passage events.
- Converts data into JSON `OBUData` messages.
- Produces them to the Kafka topic `obudata`.

### 2. Distance Calculator Service

- Consumes `OBUData` messages from the `obudata` topic.
- Matches entry and exit records per vehicle.
- Calculates the distance between them.
- Uses `client.HTTPClient` or `client.GRPCClient` to forward data
- Converts data into `types.AggregatorReq` and sends it

### 3. Aggregator Service

- Implements both **gRPC** and **HTTP servers**
- Stores data in in-memory store (`MemoryStore`)
- Exposes two HTTP endpoints:
  - `POST /aggregate` → Accepts distance data
  - `GET /invoice?obu=ID` → Returns total invoice for OBU
- Exposes `Aggregate` gRPC method for receiving distances
- Separates business logic using interfaces (`Aggregator`) and decorators (`LogMiddleware`)

---

## Technologies

- **Language**: Go
- **Messaging**: Apache Kafka
- **Serialization**: Protocol Buffers (Protobuf) and JSON
- **Transport**: gRPC and HTTP
- **Storage**: In-memory via `Storer` interface
- **Build Tool**: Go Modules

---

## Proto Messages

### AggregatorReq

```protobuf
message AggregatorReq {
  int64 obu_id = 1;
  float value = 2;
  int64 unix = 3;
}
```

### AggregatorResp

```protobuf
message AggregatorResp {}
```

---

## HTTP API

### POST /aggregate

Sends distance data.

Request Body:

```json
{
  "obu_id": 123,
  "value": 15.7,
  "unix": 1623948471
}
```

### GET /invoice?obu=123

Returns current invoice for the specified OBU.

Response:

```json
{
  "obuId": 123,
  "totalDistance": 78.4,
  "totalAmount": 25.50
}
```

---

## gRPC API

Aggregator gRPC service:

```protobuf
service Aggregator {
  rpc Aggregate(AggregatorReq) returns (AggregatorResp);
}
```

---

## Clients

Implemented clients for interacting with the Aggregator.

### gRPC Client
- `NewGRPCClient(endpoint)` sets up gRPC connection.
- `Aggregate(context, *AggregatorReq)` sends request.

### HTTP Client
- `NewHTTPClient(endpoint)` creates client.
- `Aggregate(context, *AggregatorReq)` marshals to JSON and POSTs to `/aggregate`.

---

## Setup Instructions

1. **Run Kafka and Zookeeper**  
   Using Docker or locally.

2. **Build and Start Services**

```bash
make agg
make calc
make rec
make obu
```

---

## Author

Created by [AyanokojiKiyotaka8](https://github.com/AyanokojiKiyotaka8) — open for feedback and collaboration!

---

## License

MIT License
