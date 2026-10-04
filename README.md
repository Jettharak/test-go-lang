# Go addition API

API สำหรับรับจำนวนเต็มสองจำนวนผ่าน HTTP `POST` แล้วส่งผลบวกกลับเป็น JSON

## โครงสร้างไฟล์

- `main.go` เป็นจุดเริ่มต้นของโปรแกรม
- `server.go` ตั้งค่า route และเริ่ม HTTP server
- `add_handler.go` จัดการ request/response ของ endpoint `/add`
- `addition.go` มีฟังก์ชันบวกจำนวนเต็ม

## เริ่มเซิร์ฟเวอร์

กำหนด address สำหรับฟัง request ใน `config.yml`:

```yaml
server:
  address: ":8080"
```

เปลี่ยน `address` ได้ เช่น `":9090"` เพื่อใช้พอร์ตอื่น จากนั้นรัน:

```powershell
go run .
```

เซิร์ฟเวอร์ทำงานที่ address ที่กำหนด (ค่าเริ่มต้นคือ `http://localhost:8080`) โดย endpoint คือ `POST /add`

ตัวอย่าง request:

```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/add `
  -ContentType "application/json" -Body '{"num1":10,"num2":20}'
```

response:

```json
{"result":30}
```

ตรวจสอบโค้ดและ build:

```powershell
go test ./...
go build .
```
