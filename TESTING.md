# Testing Guide: Task Management API

คู่มือการทดสอบระบบและสถาปัตยกรรม Automated Testing สำหรับโปรเจกต์ `task-management-api`

---

## Testing Architecture Overview

ระบบทดสอบของโปรเจกต์นี้ได้รับการออกแบบตามแนวคิด **Testing Pyramid** และ **Layered Architecture** เพื่อให้ทดสอบได้ครอบคลุม รวดเร็ว และแม่นยำสูง

```text
┌──────────────────────────────────────────────────────────┐
│ 4. Full E2E Router Tests (internal/routes/api_e2e_test)  │
├──────────────────────────────────────────────────────────┤
│ 3. Repository Tests (internal/repositories + testcontainers)│
├──────────────────────────────────────────────────────────┤
│ 2. Middleware Tests (internal/middleware/*_test)        │
├──────────────────────────────────────────────────────────┤
│ 1. Unit Tests (internal/services, models, util)         │
└──────────────────────────────────────────────────────────┘
```

---

## โครงสร้างไฟล์และสโคปการทดสอบ

| Package | ไฟล์ทดสอบ | ชนิดการทดสอบ | รายละเอียดการทดสอบ |
| :--- | :--- | :--- | :--- |
| `internal/models` | `enums_test.go` | Unit Test | ตรวจสอบการ Validation ค่า Enum ของ `TaskStatus` และ `TaskPriority` |
| `internal/util` | `pagination_test.go` | Unit Test | ตรวจสอบคำนวณ `Offset` และ `PaginationMeta` |
| `internal/middleware` | `auth_middleware_test.go` <br> `request_id_middleware_test.go` | Middleware Test | ตรวจสอบการ Parse JWT Bearer Header, 401 Unauthorized, และการออก `X-Request-ID` |
| `internal/services` | `mocks_test.go` <br> `auth_service_test.go` <br> `project_service_test.go` <br> `user_service_test.go` <br> `task_service_test.go` | Unit Test | ตรวจสอบ Business Logic โดยใช้ `testify/mock` จำลอง Repository 100% |
| `internal/repositories` | `testhelper_test.go` <br> `user_repository_test.go` <br> `task_repository_test.go` | Integration Test | ตรวจสอบ GORM Queries บน **PostgreSQL Container จริง** ผ่าน `testcontainers-go` |
| `internal/routes` | `api_e2e_test.go` | E2E Router Test | ทดสอบ Full User Flow (Register -> Login -> Bearer Token -> Project CRUD -> Task CRUD -> **IDOR Guard Protection**) |

---

## รูปแบบและมาตรฐานการเขียน Test (Go Testing Standards)

### 1. AAA Pattern (Arrange / Act / Assert)
ทุก Test Case จะต้องจัดสัดส่วนโค้ดออกเป็น 3 ขั้นตอนชัดเจน:
```go
// 1. Arrange: ตั้งค่า Input DTO, Mock Expectation
projRepo.On("FindByID", ctx, projectID, userID).Return(&models.Project{UserID: userID}, nil)

// 2. Act: เรียกใช้งานฟังก์ชันที่ต้องการทดสอบ
res, err := service.Create(ctx, userID, req)

// 3. Assert: ยืนยันผลลัพธ์
assert.NoError(t, err)
require.NotNil(t, res)
```

### 2. Table-Driven Tests
ใช้ Struct Array สำหรับจัดกลุ่มสภาวะทดสอบ (Happy Path, Edge Cases, Error Paths):
```go
tests := []struct {
    name          string
    req           taskdto.CreateTaskRequest
    mockSetup     func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository)
    expectedError error
}{ ... }
```

### 3. Mocking Only Dependencies
Mock เฉพาะ External Interfaces (`TaskRepository`, `UserRepository`, `ProjectRepository`) และใช้ Concrete Types สำหรับ Data Structures / Context

---

## คำสั่งการรันการทดสอบ (How to Run Tests)

### 1. รันการทดสอบทั้งหมดในโปรเจกต์
```powershell
go test -v ./...
```

### 2. รันการทดสอบเฉพาะ Service Unit Tests
```powershell
go test -v ./internal/services/...
```

### 3. รันการทดสอบเฉพาะ Middleware Tests
```powershell
go test -v ./internal/middleware/...
```

### 4. รันการทดสอบเฉพาะ E2E API Integration Test
```powershell
go test -v ./internal/routes/...
```

---

## Repository Integration Tests กับ `testcontainers-go`

ไฟล์ `internal/repositories/testhelper_test.go` มีระบบสปิน-อัป PostgreSQL Container ชั่วคราวบน Docker อัตโนมัติ:

- **หากเครื่องมี Docker Daemon ทำงานอยู่**: ระบบจะรัน Integration Test บน Postgres 16-Alpine จริง 100%
- **หากไม่มี Docker**: ระบบจะทำการ `t.Skipf()` ข้ามการทดสอบโดยอัตโนมัติ เพื่อไม่ให้การรัน Unit Test ใน CI/CD หรือเครื่อง Dev ขัดข้อง
