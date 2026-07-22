# Task Management API

REST API สำหรับจัดการงานและโปรเจกต์ พัฒนาด้วย Go, Gin, GORM และ PostgreSQL เพื่อเรียนรู้แนวคิดการออกแบบตามแนวคิด Layered Architecture และ Swagger

## Tech Stack

- Go
- Gin
- GORM
- PostgreSQL
- JWT (`golang-jwt/jwt/v5`)
- Bcrypt Password Hashing (`golang.org/x/crypto/bcrypt`)
- Docker Compose
- Swagger / Swaggo

## โครงสร้างโปรเจกต์

```text
task-management-api/
├── cmd/server          # จุดเริ่มต้นของ application
├── docs                # เอกสารออกแบบ API, database และ Swagger spec
├── internal/apperrors  # error กลางของระบบ
├── internal/bootstrap  # dependency injection
├── internal/config     # โหลด environment config
├── internal/database   # database connection, migration และ seed
├── internal/docs       # Swagger runtime docs.go สำหรับ gin-swagger
├── internal/dto        # request/response DTO
├── internal/handlers   # HTTP handler layer
├── internal/mappers    # แปลง model เป็น response DTO
├── internal/middleware # Gin middlewares
├── internal/models     # GORM models
├── internal/repositories # database access layer
├── internal/routes     # register routes
├── internal/services   # business logic layer
└── internal/util       # JWT token utilities
```

## ความสามารถปัจจุบัน

- เชื่อมต่อ PostgreSQL ผ่าน GORM
- Auto migration สำหรับ User, Project และ Task
- Seed test user ตอนเริ่ม server
- Project CRUD API
- User profile API พื้นฐาน
- Swagger UI
- แยกโครงสร้างแบบ Repository, Service, Handler

## การเตรียมฐานข้อมูล

โปรเจกต์นี้ใช้ PostgreSQL ผ่าน Docker Compose

```powershell
docker compose up -d
```

ตรวจสอบ container:

```powershell
docker ps
```

ค่าฐานข้อมูลเริ่มต้น:

```text
host: localhost
port: 5432
database: task_management
user: task_user
password: task_password
```

## การรันโปรเจกต์

ติดตั้ง dependencies:

```powershell
go mod tidy
```

รัน server:

```powershell
go run ./cmd/server
```

server จะเริ่มที่:

```text
http://localhost:8080
```

## Environment Config

ถ้าไม่ได้กำหนด environment variable ระบบจะใช้ค่า default สำหรับ local development

```text
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=task_user
DB_PASSWORD=task_password
DB_NAME=task_management
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Bangkok
JWT_SECRET=super-secret-jwt-key
```

## Swagger

เปิด Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

ถ้าแก้ Swagger annotation แล้วต้อง generate docs ใหม่:

```powershell
swag init -g cmd/server/main.go -o internal/docs --ot go
swag init -g cmd/server/main.go -o docs/swagger --ot json,yaml
```

## API Endpoints

Base URL:

```text
http://localhost:8080/api/v1
```

### 1. Auth (Public)

```text
POST /auth/register
POST /auth/login
```

ตัวอย่าง request สำหรับสมัครสมาชิก:

```json
{
  "name": "Developer Test",
  "email": "dev@example.com",
  "password": "password123"
}
```

ตัวอย่าง request สำหรับเข้าสู่ระบบ:

```json
{
  "email": "dev@example.com",
  "password": "password123"
}
```

### 2. Projects (Protected - Require `Authorization: Bearer <token>`)

```text
GET    /projects
POST   /projects
GET    /projects/:id
PATCH  /projects/:id
DELETE /projects/:id
```

### 3. Tasks (Protected - Require `Authorization: Bearer <token>`)

```text
POST   /tasks
GET    /tasks?project_id=1&status=todo&priority=high
GET    /tasks/:id
PATCH  /tasks/:id
DELETE /tasks/:id
```

### 4. Users (Protected - Require `Authorization: Bearer <token>`)

```text
GET    /users/:id
PATCH  /users/:id
DELETE /users/:id
```

## Roadmap ถัดไป

- Response wrapper
- Validation error response
- Pagination และ sorting
- Logging และ request ID middleware
- Unit test และ Integration test