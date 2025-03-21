# 🐾 Tindermals Backend

🔒 **Tindermals** is a **Go (Gin)** backend API designed to manage **animals** and **users**, featuring **JWT authentication**, **PostgreSQL storage**, and a **DDD (Domain-Driven Design)** architecture.

---
## 📌 Table of Contents
- [🚀 Installation](#-installation)
- [📂 Architecture](#-architecture)
- [✅ Testing and code quality](#-testing-and-code-quality)
- [📚 Swagger – API Documentation](#-swagger--api-Documentation)
---

## 🚀 Installation

### 1️⃣ **Clone the project**
```sh
git clone https://github.com/votre-repo/tindermals-backend.git
cd tindermals-backend
```
### 2️⃣ Setup environment variables
```sh
cp .env.example .env
```
### 3️⃣ Install dependencies
```sh
go mod tidy
```
### 4️⃣ Start the database (Postgres) and the server (compiled) with Docker
#### Only db (for local development)
```sh
docker-compose up -d db
```
#### db and server (for production) if you do that you need to create a .env.production file
```sh
docker-compose up -d
```
### 4️⃣ Start the server
```sh
go run cmd/api/main.go
```
## 📂 Architecture
```
📦 tindermals-backend
├── 📁 cmd/api                 # Main entry point (main.go)
├── 📁 internal
│   ├── 📁 adapter/http        # API Handlers (Routes)
│   ├── 📁 infra/db            # Database connection & migrations
│   ├── 📁 modules             # DDD - Modules
│   │   ├── 📁 ...
│   │   │   ├── 📁 domain      # Entities & business rules
│   │   │   ├── 📁 repository  # Database access (PostgreSQL)
│   │   │   ├── 📁 usecase     # Use cases
│   │   ├── 📁 ...
│   │   │   ├── 📁 domain
│   │   │   ├── 📁 repository
│   │   │   ├── 📁 usecase
│   ├── 📁 shared              # Security, errors, middleware
├── .env.sample                # Configuration example
├── docker-compose.yml         # DB configuration with Docker
├── lefthook.yml               # Pre-commit & pre-push hooks
└── go.mod / go.sum            # Go dependencies
```
## ✅ Testing and code quality
We use [Lefthook](https://github.com/evilmartians/lefthook) to run automated checks before every commit and push.

### 🔁 What happens before a commit?

Every time you commit, the following steps are run automatically:

- ✅ Code is formatted with `gofmt`
- ✅ Code is statically analyzed with `go vet`
- ✅ Code is linted using `golangci-lint`
- ✅ All unit tests are executed with `go test`

If any step fails, the commit will be blocked.

### 🧪 Run manually
#### 1️⃣ **Run tests**
```sh
go test ./...
```
#### 2️⃣ **Check formatting**
```sh
gofmt -s -w .
```
#### 3️⃣ **Check linting**
```sh
golangci-lint run ./...
```

### 🚀 Pre-push hook
Before pushing your code, Lefthook will again:
- Run all tests
- Check formatting
- Check linting

## 📚 Swagger – API Documentation

- The API documentation is available at:  
  👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

---

### 🛠️ How to document an endpoint

To auto-generate Swagger documentation, add special comments **above your route handlers**.  
Example:

```
// @Summary      Ping
// @Description  Check if the API is alive
// @Tags         Health
// @Success      200  {string}  string  "pong"
// @Router       /ping [get]
func Ping(c *gin.Context) {
    c.String(200, "pong")
}
```
#### ⚙️ Generate or update the documentation
Every time you add or change a route, run:
```sh
swag init -g cmd/api/main.go
```
Then, restart the server:
```sh
go run cmd/api/main.go
```
#### 🛡️ Best practices
- Always include: @Summary, @Description, @Tags, @Success, @Router
- Run swag init every time you change your routes
- Do not expose /swagger in production — or secure it with auth