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
├── start-server.bat           // สคริปต์เริ่มต้นเซิร์ฟเวอร์
├── controllers/
│   ├── auth_controller.go      // Controller สำหรับ Login
│   └── user_controller.go      // Controller สำหรับจัดการ User CRUD
├── models/
│   └── user_model.go          // Model และฟังก์ชันฐานข้อมูล
├── middlewares/
│   └── auth_middleware.go     // Middleware ตรวจสอบ JWT
├── utils/
│   └── token.go               // ฟังก์ชันจัดการ JWT
└── tests/
    ├── test-api.ps1           // สคริปต์ทดสอบ API พื้นฐาน
    ├── test-crud.ps1          // สคริปต์ทดสอบ CRUD operations
    ├── test-errors.ps1        // สคริปต์ทดสอบ error handling
    ├── TEST_RESULTS.md        // รายงานผลการทดสอบ
    └── README.md              // คำแนะนำการทดสอบ
```

## การติดตั้งและใช้งาน

### 1. ติดตั้ง Dependencies

```bash
go mod tidy
```

### 2. ตั้งค่าฐานข้อมูล

แก้ไขไฟล์ `.env` ตามการตั้งค่าฐานข้อมูลของคุณ:

```env
DB_SERVER=your-server
DB_USER=sa
DB_PASSWORD=your-password
DB_PORT=1433
DB_NAME=your-database
JWT_SECRET=your-secret-key
PORT=8080
```

### 3. รันโปรแกรม

```bash
# วิธีที่ 1: รันด้วย Go command
go run main.go

# วิธีที่ 2: ใช้ batch file (Windows)
.\start-server.bat
```

Server จะทำงานที่ port 8080

## 🧪 การทดสอบ API

ใช้สคริปต์ทดสอบใน folder `tests/`:

```powershell
# ทดสอบฟังก์ชันพื้นฐาน
.\tests\test-api.ps1

# ทดสอบ CRUD operations
.\tests\test-crud.ps1

# ทดสอบ error handling
.\tests\test-errors.ps1
```

ดูรายละเอียดเพิ่มเติมใน `tests/README.md`

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
