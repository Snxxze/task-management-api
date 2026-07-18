# Task Management API Design

ระบบนี้เป็น REST API สำหรับจัดการงาน (Task Management)

ประกอบด้วย 3 Module

- Authentication
- Projects
- Tasks

API ทั้งหมดออกแบบตามหลัก RESTful API และใช้ JSON เป็นรูปแบบในการ request & response

---

## Base URL

API Version

/api/v1

ตัวอย่าง Endpoint

GET /api/v1/tasks

POST /api/v1/tasks

---

## รูปแบบข้อมูล

Request และ Response ทุกตัวใช้รูปแบบ JSON

Content-Type

application/json

---

## Authentication

API ส่วนใหญ่ต้องยืนยันตัวตนด้วย JWT
Bearer <access_token>

---

## หลักการออกแบบ API

- URL ใช้คำนาม (Resource)
- ใช้ HTTP Method เพื่อระบุการกระทำ
- Response ใช้ HTTP Status Code มาตรฐาน
- Request และ Response ใช้ JSON
- Endpoint ทุกตัวต้องอยู่ภายใต้ /api/v1

---

## การออกแบบ API

แต่ละ Endpoint จะมีรายละเอียดเพิ่มเติมในหัวข้อถัดไป

### Authentication

| Method | Endpoint | Description | Auth |
|---------|----------|-------------|------|
| POST | /api/v1/auth/register | สมัครสมาชิก | ไม่ |
| POST | /api/v1/auth/login | เข้าสู่ระบบ | ไม่ |
| GET | /api/v1/auth/me | ข้อมูลผู้ใช้ปัจจุบัน | ใช่ |

### Projects

| Method | Endpoint | Description | Auth |
|---------|----------|-------------|------|
| GET | /api/v1/projects | รายการ Project ทั้งหมด | ใช่ |
| GET | /api/v1/projects/{id} | รายละเอียด Project | ใช่ |
| POST | /api/v1/projects | สร้าง Project | ใช่ |
| PATCH | /api/v1/projects/{id} | แก้ไข Project | ใช่ |
| DELETE | /api/v1/projects/{id} | ลบ Project | ใช่ |

### Tasks

| Method | Endpoint | Description | Auth |
|---------|----------|-------------|------|
| GET | /api/v1/tasks | รายการ Task ทั้งหมด | ใช่ |
| GET | /api/v1/tasks/{id} | รายละเอียด Task | ใช่ |
| POST | /api/v1/tasks | สร้าง Task | ใช่ |
| PATCH | /api/v1/tasks/{id} | แก้ไข Task | ใช่ |
| DELETE | /api/v1/tasks/{id} | ลบ Task | ใช่ |