# Database Design

## ภาพรวม

ระบบ Task Management ประกอบด้วย Entity หลักทั้งหมด 3 Entity ได้แก่

- User
- Project
- Task

ความสัมพันธ์ของข้อมูลเป็นแบบ One-to-Many ดังนี้

- ผู้ใช้งาน (User) 1 คน สามารถสร้าง Project ได้หลาย Project
- Project 1 ตัว สามารถมี Task ได้หลาย Task

---

## Entity Relationship Diagram (ERD)

```text
User
│
├── Projects []Project
│
▼
Project
│
├── UserID
├── User
├── Tasks []Task
│
▼
Task
│
├── ProjectID
└── Project
```

---

## ความสัมพันธ์ของ Entity

### User → Project

**Relationship**

```
1 User : N Projects
```

**Foreign Key**

```
projects.user_id
```

**คำอธิบาย**

- ผู้ใช้งาน 1 คน สามารถเป็นเจ้าของ Project ได้หลาย Project
- Project ทุกตัวต้องมีเจ้าของเพียง 1 คน

---

### Project → Task

**Relationship**

```
1 Project : N Tasks
```

**Foreign Key**

```
tasks.project_id
```

**คำอธิบาย**

- Project 1 ตัว สามารถมี Task ได้หลาย Task
- Task ทุกตัวต้องอยู่ภายใต้ Project เพียง 1 Project

---

## Model Relationship

### User

```go
Projects []Project `gorm:"foreignKey:UserID"`
```

**อธิบาย**

ใช้สำหรับโหลด Project ทั้งหมดของ User

---

### Project

```go
UserID uint
User   User `gorm:"foreignKey:UserID"`

Tasks []Task `gorm:"foreignKey:ProjectID"`
```

**อธิบาย**

- UserID คือ Foreign Key ที่อ้างอิงไปยัง User
- User คือ Association สำหรับโหลดข้อมูลเจ้าของ Project
- Tasks คือรายการ Task ทั้งหมดภายใน Project

---

### Task

```go
ProjectID uint
Project   Project `gorm:"foreignKey:ProjectID"`
```

**อธิบาย**

- ProjectID คือ Foreign Key
- Project คือ Association สำหรับโหลดข้อมูล Project ที่ Task สังกัดอยู่

---

## Relationship Summary

| Parent | Child | Type | Foreign Key |
|---------|--------|------|-------------|
| User | Project | One-to-Many | user_id |
| Project | Task | One-to-Many | project_id |

---

## Database Tables

```
users
│
├── id
├── name
├── email
└── password_hash

projects
│
├── id
├── name
├── description
└── user_id

tasks
│
├── id
├── title
├── description
├── status
└── project_id
```