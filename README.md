# Task Management API

REST API สำหรับจัดการงานและโปรเจกต์ พัฒนาด้วย Go, Gin, GORM และ PostgreSQL เพื่อเรียนรู้แนวคิดการออกแบบตามแนวคิด Layered Architecture และ Swagger

## Tech Stack

- Go
- Gin
- GORM
- PostgreSQL
- Docker Compose
- Swagger / Swaggo

## โครงสร้างโปรเจกต์

```text
task-management-api/
├── cmd/server          # จุดเริ่มต้นของ application
├── docs                # เอกสารออกแบบ API, database, tech debt และ Swagger spec
├── internal/apperrors  # error กลางของระบบ
├── internal/bootstrap  # dependency injection
├── internal/config     # โหลด environment config
├── internal/database   # database connection, migration และ seed
├── internal/docs       # Swagger runtime docs.go สำหรับ gin-swagger
├── internal/dto        # request/response DTO
├── internal/handlers   # HTTP handler layer
├── internal/mappers    # แปลง model เป็น response DTO
├── internal/models     # GORM models
├── internal/repositories # database access layer
├── internal/routes     # register routes
└── internal/services   # business logic layer
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

### Projects

```text
GET    /projects
POST   /projects
PATCH  /projects/:id
DELETE /projects/:id
```

ตัวอย่าง request สำหรับสร้าง project:

```json
{
  "name": "Learning Go API",
  "description": "Practice layered architecture",
  "color": "#3B82F6"
}
```

### Users

```text
GET    /users/:id
PATCH  /users/:id
DELETE /users/:id
```

ตัวอย่าง request สำหรับแก้ไข user:

```json
{
  "name": "Updated User",
  "email": "updated@example.com"
}
```

## Seed Data

ตอนเริ่ม server ระบบจะ seed user สำหรับทดสอบให้อัตโนมัติ:

```text
name: Test User
email: test@example.com
password_hash: hashed-password
```

หมายเหตุ: ตอนนี้ Project API ยังใช้ `userID := uint(1)` ชั่วคราว จนกว่าจะเพิ่มระบบ Authentication/JWT

## Roadmap ถัดไป

- Authentication ด้วย JWT
- Task CRUD API
- Response wrapper
- Validation error response
- Pagination, search และ sorting
- Logging และ request ID
- Unit test และ handler test
- Dockerfile
