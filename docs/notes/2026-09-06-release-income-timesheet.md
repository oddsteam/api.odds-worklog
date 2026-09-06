# Release note: income_from_timesheet

สรุปทุก commit ตั้งแต่ tag `v5.1.1` มาถึงตอนนี้ (2026-09-06) แบ่งเป็นของที่ปล่อยไปรอบก่อน กับของรอบนี้

## ครั้งก่อน (ทยอย merge ตั้งแต่ v5.1.1)

- **เปิด timesheet inbox** — `GET /v1/timesheet-event-logs` ดูเหตุการณ์จาก timesheet service ย้อนหลัง (ใครก็ log in ดูได้ ไม่ guard เฉพาะ admin เหมือน SAP failure log)
- **fix inbox แสดง 0 วันทำงาน/OT + ชื่อ site/employee ว่าง** — field tag ไม่ตรงกัน (JSON ที่ web อ่านเป็น camelCase แต่ struct ที่ใช้ร่วมกับ RabbitMQ consumer เป็น snake_case ของฝั่ง publisher) แก้ด้วยการ map เป็น response type แยกในชั้น api แทนที่จะไปแก้ tag ของ struct ที่ consumer ใช้
- **เพิ่มระบบ income_from_timesheet** — consume event `timesheet.monthly_summary.published` จาก RabbitMQ แล้ว upsert เข้า collection `income_from_timesheet` โดยใช้ payroll calculation เดียวกับ `income` ปกติ (`CreatePayroll`/`UpdatePayroll`) แยก collection กันเด็ดขาด ไม่แตะ `income`/`user` เดิม
- **Dual-write จาก manual add/edit income** — ทุกครั้งที่ user กรอก/แก้ income เองผ่านฟอร์ม ระบบจะ mirror ข้อมูลไปที่ `income_from_timesheet` ด้วย (`mirrorIncomeToTimesheet`) เพื่อให้ collection นี้มีข้อมูลครบทุกคน ไม่ใช่แค่คนที่มาจาก timesheet
- **Export endpoints สำหรับ income_from_timesheet** — เพิ่ม export หลายฟอร์แมต (CSV/SAP/PEAK) แล้ว refactor export usecase ให้ใช้ source adapter กลางร่วมกับของเดิม
- **เก็บ note ได้ทั้ง income และ income_from_timesheet**
- **ย้าย PEAK Code เข้าไปเก็บใน income เอง** — ตอน export ไม่ต้อง join กับ `user` แล้ว
- **เก็บ site name ใน income เอง** — เหตุผลเดียวกัน ไม่ต้อง join
- **WHT rate 3% → 5%**
- **Refactor**: รวม SAP export failure ports (driving/driven ที่หน้าตาเหมือนกัน) เป็น interface เดียว

## ครั้งนี้ (2026-09-06)

- **specialIncome เลิก hard code เป็น "0"** — ตอน sync จาก timesheet event คำนวณ special hourly rate จาก `user.DailyIncome / 8` แทน (`business/usecases/sync_income_from_timesheet.go`)
- **แปลงหน่วย OT จาก timesheet ให้ตรง** — `OvertimeDays` ที่ timesheet ส่งมาเป็น**วัน** คูณ 8 ก่อนเก็บลง `WorkingHours` ของ worklog ให้หน่วยเป็น**ชั่วโมง** ตรงกับที่ payroll ใช้คำนวณจริง และตรงกับฝั่ง manual add/edit ที่ user กรอกเป็นชั่วโมงอยู่แล้ว (เลยไม่กระทบ dual-write ฝั่งนั้น)
- อัปเดต unit test ที่เกี่ยวข้องให้ตรงกับ behavior ใหม่ — `go build ./...` และ `go test ./...` ผ่านหมด (422 tests)

## รอดำเนินการ

- ย้ายข้อมูลจาก `income_from_timesheet` ไป `income` เพื่อให้ขึ้นใน history — รอพี่จั๊วะทำ
- Sync ข้อมูลไป `income` จริง — รอพี่จั๊วะทำ
- เอา toggle ออก และเอาการ dual-write (save ลง `income` + `income_from_timesheet` พร้อมกัน) ออก — ทำได้หลัง sync ข้อมูลครบเท่านั้น

## ต้องแจ้ง user

- ถ้า OT rate ของใครไม่เท่ากับ rate ปกติ (`dailyRate / 8`) ให้ user ไปกรอก/แก้เวลาเพิ่มเองผ่านฟอร์ม add/edit income — ไม่ได้ทำ auto-detect/adjust ให้ในรอบนี้
