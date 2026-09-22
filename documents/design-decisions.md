# Assumptions & Design Decisions

## Login error handling: ใช้ `invalid credentials` เดียวทุกกรณี

`Login` return `entity.ErrInvalidCredentials` เสมอ ไม่ว่าจะเป็นกรณี email ไม่มีในระบบ หรือ password ผิด ไม่แยก error ให้ client รู้ว่าสาเหตุคืออะไร

กันไว้ไม่ให้คนไล่เดา email ในระบบได้ ถ้าตอบคนละ error กันระหว่าง "ไม่มี email นี้" กับ "password ผิด" อาจจะมีผู้ไม่หวังดีจะไล่ยิง email ทีละตัวแล้วรู้ได้ทันทีว่าตัวไหนมีอยู่จริง ทั้งที่ยังไม่ต้องรู้ password เลย

## ค่า default ของ `GetList`

ถ้าผู้เรียกไม่ระบุค่ามา `GetList` จะใช้ค่า default ดังนี้:
- `limit = 20`
- ช่วงวันที่ = 1 เดือนล่าสุด (`to = ปัจจุบัน`, `from = to - 1 เดือน`)
ป้องกันไม่ให้ query ไม่มีขอบเขต (ดึง user ทั้งหมดในฐานข้อมูลออกมาหมด) เวลาที่ client เรียก `GET /users` โดยไม่ส่ง query param มาเลย

## TTL ของ token กับค่า config การต่อ MongoDB

`TOKEN_TTL` (default `24h`) ให้ปรับผ่าน env var ได้ ไม่ hardcode ไว้ในโค้ด เพราะแต่ละ environment อาจอยากได้อายุ token ต่างกัน (เช่น dev อยากให้ token อยู่นานๆไม่ต้อง login ซ้ำบ่อย ส่วน prod อาจอยากให้สั้นกว่านี้เพื่อความปลอดภัย) ปรับที่ env ได้เลยไม่ต้องแก้โค้ดแล้ว build ใหม่

ค่าต่อ MongoDB (`MONGO_MINPOOL_SIZE`, `MONGO_MAXPOOL_SIZE`, `MONGO_MAX_IDLE_TIME`, `MONGO_CONNECT_TIMEOUT`) ก็ให้ปรับผ่าน env เหมือนกัน โดยมี default ให้ (`3`, `30`, `30s`, `10s`) ที่เหมาะกับงานขนาดเล็ก แต่ก็เผื่อไว้เผื่อโหลดจริงมากขึ้นจะได้ทำการปรับ config ง่ายขึ้นโดยที่ไม่ต้อง build code ใหม่

## database name hoard code (constants value)

ชื่อ database เขียนเป็น constant ไว้ในโค้ดตรงๆ ไม่ได้ทำเป็น env var เหมือนค่า Mongo อื่นๆ เพราะปกติแล้วชื่อ database จะไม่มีการเปลี่ยนบ่อยๆ ต่างจาก host/port/credential ที่เปลี่ยนไปตาม environment ได้ตลอด

ถ้าจะเปลี่ยนชื่อ database จริงๆ ก็ต้อง build code ใหม่อยู่แล้ว เพื่อ test ให้แน่ใจว่า connect เข้า database ใหม่ได้จริงก่อน deploy ไม่ใช่แค่เปลี่ยน env var แล้วปล่อยผ่านไปเลยโดยไม่มีการตรวจสอบอะไร การ hardcode ไว้จึงบังคับให้ต้องผ่าน build/test flow ปกติทุกครั้งที่จะเปลี่ยนค่านี้

## การ hash password
- Password hash ด้วย **bcrypt** (`golang.org/x/crypto/bcrypt`) ซึ่งเป็นตัวเลือกมาตรฐานสำหรับเก็บ password ใน Go 

## การแบ่ง layer code
- `internal/core` (entity, domain DTO, service logic, และ interface ใน `port`) ไม่มีการ import MongoDB, HTTP, หรือ JWT library 
- กลุ่มที่มี library จะอยู่ที่ `internal/adapter` และ `internal/infra` แล้วนำมาประกอบกันที่ `cmd/main.go` เพื่อไม่ให้ Library อยู่ใน layter service logic ป้องกัน impact จากการอัพเดท Library และทำตามหลักการของ Hexagonal architechure