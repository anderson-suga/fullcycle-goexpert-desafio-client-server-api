# Full Cycle - Go Expert - Challenge - Client-Server-API

This project consists of two systems written in Go: a `client` and a `server`. The goal is to demonstrate the use of HTTP web servers, context for timeout management, database manipulation (SQLite), and file handling.

## Requirements

- Go 1.20+
- GCC (required for `go-sqlite3` CGO compilation)

## Project Structure

```text
.
├── client/
│   └── client.go
├── server/
│   └── server.go
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Architecture

### Server (`server.go`)

- Runs on port `:8080`.
- Endpoint: `/cotacao`.
- Fetches USD-BRL exchange rate from `awesomeapi` (Timeout: 200ms).
- Persists the bid value in a SQLite database `cotacoes.db` (Timeout: 10ms).

### Client (`client.go`)

- Requests the exchange rate from the local server (Timeout: 300ms).
- Receives the JSON response.
- Saves the data to a file `cotacao.txt` in the format: `Dólar: {value}`.

## How to Run

### 1. Setup

Initialize the module and download dependencies:

```bash
go mod tidy
```

### 2. Run the Server

Open a terminal in the project root and execute:

```bash
go run server/server.go
```

You should see: "Server started on port 8080"

Note: The server will create a `cotacoes.db` SQLite file in the root directory.

### 3. Run the Client

Open another terminal window (also in the project root) and execute:

```bash
go run client/client.go
```

You should see: "Process finished successfully. Bid saved."

### 4. Verify Results

Check if the file `cotacao.txt` was created in the root directory:

```bash
cat cotacao.txt
# Output example: Dólar: 6.0521
```

Check the database records (optional):

```bash
sqlite3 cotacoes.db "select * from cotacoes;"
```

## </br>

</br>

> This is a postgraduate degree challenge from [Full Cycle - Go Expert](https://goexpert.fullcycle.com.br/pos-goexpert/)
