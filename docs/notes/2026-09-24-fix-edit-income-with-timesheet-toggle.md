# แก้บั๊ก: เปิด toggle "ใช้ข้อมูลจาก Timesheet" แล้วแก้ไข income ไม่ได้ (500)

## อาการ

เปิด toggle บนหน้า individual แล้วกดแก้ไข income → `PUT /v1/incomes/{id}` ตอบ
`500 {"message":"mongo: no documents in result"}`

## Root cause

`_id` ของ collection `income` กับ `income_from_timesheet` **ไม่เคยตรงกัน** — แม้แต่ record ที่
`mirrorIncomeToTimesheet` เขียน mirror ไว้ ก็ตั้งใจไม่ copy id มา (ดูคอมเมนต์ใน
`business/usecases/mirror_income_to_timesheet.go`)

พอเปิด toggle หน้าเว็บอ่าน record จาก `/v1/income-from-timesheet/current-month/:id` แล้วเก็บ
`res.id` ลง `IncomeFlag.id` ซึ่งเป็น id ของ **อีก collection** — แต่ตอนบันทึกยังยิง
`PUT /v1/incomes/{IncomeFlag.id}` ที่ไปหาใน collection `income` เสมอ จึงหาไม่เจอทุกครั้ง

## Flow เดิม (พัง)

```mermaid
flowchart TD
    A["ผู้ใช้เปิด toggle<br/>ใช้ข้อมูลจาก Timesheet"] --> B["add-income.checkStatusUser()"]
    B --> C["GET /v1/income-from-timesheet/current-month/:userId"]
    C --> D[("income_from_timesheet")]
    D --> E["res.id = _id ของ income_from_timesheet"]
    E --> F["IncomeFlag.id = res.id<br/>(ทั้งข้อมูลที่แสดง และเป้าหมายที่จะบันทึก)"]
    F --> G["กด Confirm → PUT /v1/incomes/{IncomeFlag.id}"]
    G --> H["updateIncomeUsecase.UpdateIncome()<br/>repo.GetIncomeByID(id, userId)"]
    H --> I[("income")]
    I --> J{"เจอ record ไหม?"}
    J -->|"ไม่เจอ — id เป็นของคนละ collection"| K["mongo: no documents in result"]
    K --> L["HTTP 500 ❌"]

    style L fill:#ffd7d7,stroke:#c0392b
    style F fill:#ffe9cc,stroke:#e08b00
```

## Flow ใหม่ (แก้แล้ว)

แยก **แหล่งข้อมูลที่แสดง** ออกจาก **record ที่จะบันทึก** — toggle คุมแค่ตัวเลขที่โชว์
ส่วนฟอร์มเขียนลง collection `income` เสมอ (ตรงกับแผน dual-write เดิม ที่ `income` เป็น
source of truth และ `income_from_timesheet` เป็น mirror)

```mermaid
flowchart TD
    A["ผู้ใช้เปิด toggle<br/>ใช้ข้อมูลจาก Timesheet"] --> B["add-income.checkStatusUser()"]
    B --> C["loadDisplayedIncome()<br/>GET /income-from-timesheet/current-month/:userId"]
    B --> D["loadEditTarget()<br/>GET /incomes/current-month/:userId"]

    C --> C1[("income_from_timesheet")]
    C1 --> C2["เติมตัวเลขลงฟอร์ม<br/>(ไม่แตะ IncomeFlag)"]

    D --> D1[("income")]
    D1 --> D2{"มี income ของเดือนนี้ไหม?"}
    D2 -->|"มี"| D3["IncomeFlag.id = income.id<br/>IncomeFlag.isUpdate = true"]
    D2 -->|"ไม่มี / error"| D4["IncomeFlag.id = ''<br/>IncomeFlag.isUpdate = false"]

    C2 --> E["กด Confirm"]
    D3 --> E
    D4 --> E

    E --> F{"IncomeFlag.isUpdate?"}
    F -->|"true"| G["PUT /v1/incomes/{income.id}"]
    F -->|"false"| H["POST /v1/incomes"]

    G --> I["UpdateIncome → เจอ record ใน income ✅"]
    H --> J["AddIncome → สร้าง record ใน income ✅"]

    I --> K["mirrorIncomeToTimesheet(period)<br/>อัปเดต record ใน income_from_timesheet<br/>(sites ไม่หาย)"]
    J --> K
    K --> L["HTTP 200 ✅"]

    style L fill:#d7f5d7,stroke:#27a05a
    style D fill:#d7e8ff,stroke:#2b6cb0
```

เคสที่เป็นไปได้ทั้งหมดเมื่อเปิด toggle:

| มี record ใน income_from_timesheet | มี record ใน income | ฟอร์มแสดง | กดบันทึกแล้ว |
|---|---|---|---|
| มี | มี | ตัวเลขจาก timesheet | `PUT /incomes/{income.id}` → 200 |
| มี | ไม่มี | ตัวเลขจาก timesheet | `POST /incomes` → 200 แล้ว mirror ทับ record timesheet ของ period นั้น |
| ไม่มี | มี | ฟอร์มเปล่า | `PUT /incomes/{income.id}` → 200 |
| ไม่มี | ไม่มี | ฟอร์มเปล่า | `POST /incomes` → 200 |

## สิ่งที่แก้

**web.odds-worklog**
- `src/app/shared/components/add-income/add-income.component.ts` — แยก `loadDisplayedIncome()` /
  `loadEditTarget()` และให้ `setEditTarget()` เป็นที่เดียวที่เขียน `IncomeFlag.id` / `isUpdate`
  (เดิม `isUpdate` ไม่เคยถูก set ตอนโหลดหน้า อาศัยค่า static default `true` เอา)
- `add-income.component.spec.ts` — เพิ่มเทสต์ครอบ 4 เคสข้างบน

**api.odds-worklog** (defense in depth — ให้ error อ่านออก ไม่ใช่ 500 ลอย ๆ)
- `business/usecases/update_income.go` — เพิ่ม `ErrIncomeNotFound`
- `repositories/income.go` — map `mongo.ErrNoDocuments` → `ErrIncomeNotFound`
- `api/income/handler.go` — ตอบ **404** แทน 500 เมื่อไม่เจอ record
- `docs/openapi.yaml` — ประกาศ response 404 ของ `PUT /v1/incomes/{id}`

## ยังเหลือ (ตาม release note 2026-09-06)

พอ migrate ข้อมูลเข้า `income` ครบและเอา toggle ออก `loadEditTarget()` จะไม่จำเป็นอีก —
`loadDisplayedIncome()` จะกลับไปอ่าน `income` ทางเดียวเหมือนเดิม
