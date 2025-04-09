# DepGit

DepGit is a modern Git repository management system that provides a web interface and Git server for managing repositories. It's designed to be a self-hosted alternative to services like GitHub or GitLab, with a focus on simplicity and performance.

## Features

- Git repository hosting via SSH
- Web-based repository management
- User management and authentication
- Role-based access control for repositories
- Support for multiple storage backends (local filesystem and MinIO)
- Modern web interface built with Svelte and Tailwind CSS

## Project Structure

```
.
├── api/                  # OpenAPI specifications
├── cmd/                  # Application entry points
│   └── app/              # Main application
├── docs/                 # Documentation
├── internal/             # Internal packages
│   ├── config/           # Configuration handling
│   ├── database/         # Database access and models
│   │   └── migrations/   # Database migration tools
│   ├── gen/              # Generated code
│   │   └── api/          # Generated API code from OpenAPI specs
│   ├── git/              # Git server implementation
│   ├── models/           # Data models
│   ├── share/            # Shared utilities
│   │   ├── errors/       # Error definitions
│   │   └── logger/       # Logging utilities
│   ├── stroage/          # Storage backends (file, MinIO)
│   └── web/              # Web server and API handlers
├── migrations/           # SQL migration files
├── scripts/              # Utility scripts
└── web/                  # Web frontend (Svelte)
    ├── public/           # Static assets
    └── src/              # Frontend source code
        ├── components/   # Reusable UI components
        └── pages/        # Page components
```

## Technology Stack

### Backend
- **Go**: Main programming language
- **SQLite**: Database
- **Echo**: HTTP server framework
- **go-git**: Git implementation in Go
- **gliderlabs/ssh**: SSH server implementation
- **MinIO**: Object storage (optional)

### Frontend
- **Svelte**: UI framework
- **TypeScript**: Programming language
- **Tailwind CSS**: Utility-first CSS framework
- **Vite**: Build tool

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Node.js 16 or higher
- npm or yarn
- SQLite

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/GoldenDeals/DepGit.git
   cd DepGit
   ```

2. Install backend dependencies:
   ```bash
   go mod download
   ```

3. Install frontend dependencies:
   ```bash
   cd web
   npm install
   cd ..
   ```

4. Create a `.env` file in the project root with the following content:
   ```
   APP_ENV=dev
   DB_PATH=./depgit.db
   DB_MIGRATIONS_PATH=./migrations
   WEB_ADDRESS=:8080
   WEB_STATIC_DIR=./web/dist
   GIT_SSH_ADDRESS=:2222
   DEPGIT_SSH_GIT_HOSTKEY=./path/to/ssh/host/key
   ```

5. Generate an SSH host key:
   ```bash
   ssh-keygen -t rsa -f ./ssh_host_key -N ""
   ```

6. Apply database migrations:
   ```bash
   make db-migrate
   ```

7. Build the application:
   ```bash
   make build-full
   ```

### Running the Application

Start the application:
```bash
./build/depgit
```

For development, you can use:
```bash
make live
```

This will start the application with hot reloading using Air.

## Development

### Building

- `make build`: Build the application
- `make build-debug`: Build with debug information
- `make build-release`: Build optimized release binary
- `make build-full`: Build backend and frontend

### Web Development

- `make web-install`: Install web dependencies
- `make web-dev`: Start web development server
- `make web-build`: Build web application

### Testing

- `make test`: Run tests
- `make test-verbose`: Run tests with verbose output
- `make test-coverage`: Run tests with coverage report
- `make test-html`: Generate HTML coverage report
- `make test-git`: Run git test script

### API Generation

- `make gen-api`: Generate API code from OpenAPI spec

## Database Migrations

Database migrations are automatically applied when the application starts. See [migrations/README.md](internal/database/migrations/README.md) for more information on how to create and manage migrations.

## Storage Backends

DepGit supports multiple storage backends:

1. **Local Filesystem**: Stores Git objects on the local filesystem
2. **MinIO**: Stores Git objects in a MinIO object storage server

The storage backend can be configured in the `.env` file.

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
