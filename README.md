# Simple RESTful API with Go

REST API ที่สร้างด้วยภาษา Go ตามโครงสร้างแบบ MVC (Model-View-Controller) ใช้ Gin Gonic เป็นเว็บเฟรมเวิร์กหลักและเชื่อมต่อกับ Microsoft SQL Server

## ฟีเจอร์หลัก

- **Authentication**: JWT-based authentication with bcrypt password hashing
- **User Management**: CRUD operations for user management
- **Database**: Microsoft SQL Server integration
- **Security**: Protected routes with JWT middleware
- **Architecture**: MVC pattern implementation

## โครงสร้างโปรเจกต์

```
├── main.go                     // ไฟล์เริ่มต้นโปรแกรม
├── controllers/
│   ├── auth_controller.go      // Controller สำหรับ Login
│   └── user_controller.go      // Controller สำหรับจัดการ User CRUD
├── models/
│   └── user_model.go          // Model และฟังก์ชันฐานข้อมูล
├── middlewares/
│   └── auth_middleware.go     // Middleware ตรวจสอบ JWT
└── utils/
    └── token.go               // ฟังก์ชันจัดการ JWT
```

## การติดตั้งและใช้งาน

### 1. ติดตั้ง Dependencies

```bash
go mod tidy
```

### 2. ตั้งค่าฐานข้อมูล

แก้ไข connection string ในไฟล์ `models/user_model.go`:

```go
connString := "server=localhost;user id=sa;password=YourPassword123;port=1433;database=TestDB"
```

### 3. รันโปรแกรม

```bash
go run main.go
```

Server จะทำงานที่ port 8080

## API Endpoints

### Authentication
- `POST /login` - เข้าสู่ระบบและรับ JWT token

### User Management (ต้องมี Bearer Token ยกเว้น POST /users)
- `POST /users` - สร้างผู้ใช้ใหม่ (ไม่ต้องมี token)
- `GET /users` - ดูข้อมูลผู้ใช้ทั้งหมด
- `GET /users/:id` - ดูข้อมูลผู้ใช้รายคน
- `PUT /users/:id` - อัปเดตข้อมูลผู้ใช้
- `DELETE /users/:id` - ลบผู้ใช้

## ตัวอย่างการใช้งาน

### 1. สร้างผู้ใช้ใหม่
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### 2. เข้าสู่ระบบ
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. ดูข้อมูลผู้ใช้ทั้งหมด (ต้องมี token)
```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. อัปเดตข้อมูลผู้ใช้
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "full_name": "Updated Name"
  }'
```

### 5. ลบผู้ใช้
```bash
curl -X DELETE http://localhost:8080/users/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## ความปลอดภัย

- รหัสผ่านถูกเข้ารหัสด้วย bcrypt ก่อนเก็บในฐานข้อมูล
- JWT token มีระยะเวลาหมดอายุ 24 ชั่วโมง
- Protected routes ต้องการ Bearer Token ใน Authorization header
- Middleware ตรวจสอบความถูกต้องของ token ทุกครั้ง

## การปรับแต่ง

- แก้ไข JWT secret key ในไฟล์ `utils/token.go`
- ปรับ connection string ฐานข้อมูลในไฟล์ `models/user_model.go`
- ปรับระยะเวลาหมดอายุของ token ในฟังก์ชัน `GenerateToken`

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/denisenkom/go-mssqldb` - SQL Server driver
- `github.com/dgrijalva/jwt-go` - JWT implementation
- `golang.org/x/crypto` - Password hashing with bcrypt
