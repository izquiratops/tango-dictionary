# Tango 🎋

Tango is a Japanese-English dictionary that provides a web-based interface. It provides a simple way to look up Japanese words and their English translations. While it's still in development, the goal is to create a resource similar to [Jisho.org](https://jisho.org).

This project use the content provided by JMDict Simplified. Without this valuable resource, building this dictionary would have been significantly more time-consuming.

## Repository Structure

```
tango/
├── client/             # Web server implementation
│   ├── main.go         # Entry point for the client application
│   ├── server/         # Server implementation (routes, handlers)
│   ├── static/         # Static assets (mounted volume)
│   └── template/       # HTML templates (mounted volume)
│
├── common/             # Shared code between client and import tool
│   ├── database/       # Database connection and operations
│   ├── types/          # Common data structures
│   └── utils/          # Utility functions including config loading
│
├── import/             # Dictionary import tool
│   └── main.go         # Entry point for the import process
│
├── jmdict_source/      # Dictionary data and search index (mounted volume)
│
├── Dockerfile          # Container build instructions
├── docker-compose.yml  # Multi-container application setup
└── Caddyfile           # Caddy server configuration
```

## Deployment

The application is deployed using Docker Compose with the following services:

- `client`: The main Tango application
- `mongo`: MongoDB database service
- `caddy`: Reverse proxy for handling web requests
- `watchtower`: Automatic container updater (checks daily)

## Getting Started

### Prerequisites

- Docker and Docker Compose
- A MongoDB instance installed somewhere
- [JMDict Simplified](https://github.com/scriptin/jmdict-simplified) source files

### Setup

1. Clone the repository
2. Create a `.env` file with necessary configuration
	- TANGO_VERSION
	- TANGO_MONGO_RUNS_LOCAL
	- MONGO_INITDB_ROOT_USERNAME (optional)
	- MONGO_INITDB_ROOT_PASSWORD (optional)
3. Place JMdict source files in the `jmdict_source` directory
4. Run the import tool to prepare the database:
   ```
   cd import
   go run main.go
   ```
5. Start the application using Docker Compose:
   ```
   docker-compose up -d
   ```
