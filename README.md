# 🐾 Tindermals Backend

🔒 **Tindermals** is a **Go (Gin)** backend API designed to manage **animals** and **users**, featuring **JWT authentication**, **PostgreSQL storage**, and a **DDD (Domain-Driven Design)** architecture.

---
## 📌 Table of Contents
- [🚀 Installation](#-installation)
- [📂 Architecture](#-architecture)
- [✅ Testing and code quality](#-testing-and-code-quality)
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
### 4️⃣ Start the database (Postgres) with Docker
```sh
docker-compose up -d
```
### 5️⃣ Start server
```sh
go run main.go
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