# go-grpc-crud

Shared Go library providing reusable building blocks for gRPC microservices in the NextJS-GO stack.

---

## Overview

Used as a git submodule in each service (`services/lib`). Contains no business logic — only generic infrastructure.

## Packages

### `crud`

Schema definition types for the admin UI. Services expose a `GetSchema` gRPC method returning a `crud.Schema` so the frontend admin panel can render forms dynamically.

```go
type Schema struct {
    Name        string
    DisplayName string
    Fields      []Field
    Searchable  []string
}
```

### `server`

Generic list handler for gRPC services with pagination, search, filtering, and sorting via GORM.

```go
server.DefaultList[T](ctx, db, req, cfg)
```

### `proxy`

Fiber middleware helpers for proxying gRPC list/CRUD calls to REST endpoints in `backend-core`. Handles protojson marshalling, pagination, and error mapping.

### `errors`

Shared error types (`Problem`, `FieldError`) that map gRPC status codes to RFC 7807 Problem Details responses. Used consistently across all services and the backend gateway.

## Usage

Add as a git submodule in your service:

```bash
git submodule add https://github.com/Mognus/go-grpc-crud services/lib
```

Then import:

```go
import (
    "github.com/Mognus/go-grpc-crud/crud"
    "github.com/Mognus/go-grpc-crud/server"
    "github.com/Mognus/go-grpc-crud/proxy"
    "github.com/Mognus/go-grpc-crud/errors"
)
```
