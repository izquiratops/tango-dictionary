# Tango 🎋

Tango is a Japanese-English dictionary that provides a web-based interface. It provides a simple way to look up Japanese words and their English translations. While it's still in development, the goal is to create a resource similar to [Jisho.org](https://jisho.org).

This project use the content provided by [JMDict Simplified](https://github.com/scriptin/jmdict-simplified). Without this valuable resource, building this dictionary would have been significantly more time-consuming.

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
└── docker-compose.yml  # Multi-container application setup
```
