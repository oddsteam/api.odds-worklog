# Income monthly flow

Sequence ของ 1 เดือน: sync กยศ. → individual add/update income → admin export income

## Data sources on export

| Source | Mongo collection | Role |
|--------|------------------|------|
| Income records (primary) | `income` | Filtered by `role` + `submitDate` range; user fields are snapshotted at add/update |
| Student loan deduction | `studentloan` | Matched to income by `BankAccountName` / borrower name |
| Export audit log | `export` | Filename + date after successful export (not input data) |

Export does **not** re-read the live `user` collection.

## Sequence (one month)

```mermaid
sequenceDiagram
    autonumber
    actor Ops as Ops / Cron
    participant SLWeb as กยศ Web<br/>(slfrd.dsl.studentloan.or.th)
    participant Script as get_student_loan.go
    participant MongoSL as MongoDB<br/>studentloan
    actor Ind as Individual
    participant Web as Web App
    participant API as API<br/>(income handler)
    participant AddUC as AddIncome Usecase
    participant UpdUC as UpdateIncome Usecase
    participant ExpUC as ExportIncome Usecase
    participant UserRepo as MongoDB<br/>user
    participant IncRepo as MongoDB<br/>income
    participant ExpRepo as MongoDB<br/>export
    participant File as CSV / SAP File
    actor Admin as Admin

    Note over Ops,MongoSL: 1) Sync กยศ เข้ามา (เดือนปัจจุบัน พ.ศ.)
    Ops->>Script: รัน sync (SESSION + CSRF)
    Script->>SLWeb: POST EmployeeReport/getDataByPage<br/>(month, year พ.ศ.)
    SLWeb-->>Script: JSON รายชื่อ + ยอดหัก
    Script->>Script: CreateStudentLoanList + CreateIDForLoans
    Script->>MongoSL: Upsert โดย monthYear<br/>(SaveStudentLoans)

    Note over Ind,IncRepo: 2) Individual add income
    Ind->>Web: กรอกวันทำงาน / special income
    Web->>API: POST /v1/incomes
    API->>AddUC: AddIncome(req, userId)
    AddUC->>UserRepo: GetByID(userId)
    UserRepo-->>AddUC: user profile<br/>(ชื่อ, บัญชี, อัตราวัน, VAT…)
    AddUC->>IncRepo: GetIncomeUserByYearMonth<br/>(กันซ้ำเดือนเดียวกัน)
    IncRepo-->>AddUC: not found
    AddUC->>AddUC: CreatePayroll<br/>(snapshot ข้อมูล user ลง income)
    AddUC->>IncRepo: AddIncome(income)
    IncRepo-->>AddUC: saved
    AddUC-->>API: income
    API-->>Web: 200
    Web-->>Ind: บันทึกสำเร็จ

    Note over Ind,IncRepo: 3) Individual update income
    Ind->>Web: แก้ไข income
    Web->>API: PUT /v1/incomes/:id
    API->>UpdUC: UpdateIncome(id, req, userId)
    UpdUC->>UserRepo: GetByID(userId)
    UserRepo-->>UpdUC: user profile ปัจจุบัน
    UpdUC->>IncRepo: GetIncomeByID(id, userId)
    IncRepo-->>UpdUC: income record
    UpdUC->>UpdUC: UpdatePayroll<br/>(คำนวณใหม่ + อัปเดต snapshot)
    UpdUC->>IncRepo: UpdateIncome(income)
    IncRepo-->>UpdUC: saved
    UpdUC-->>API: income
    API-->>Web: 200
    Web-->>Ind: อัปเดตสำเร็จ

    Note over Admin,File: 4) Admin export income (เช่น individual เดือนนี้)
    Admin->>Web: Export CSV / SAP<br/>(role=individual, ช่วงเดือน)
    Web->>API: POST /v1/incomes/export<br/>หรือ /export/format/SAP
    API->>API: CanExportIncome?<br/>(full admin เท่านั้น)
    API->>ExpUC: ExportIncome…(role, start, end)

    ExpUC->>IncRepo: GetAllIncomeByRoleStartDateAndEndDate<br/>(role + submitDate ในช่วงเดือน)
    IncRepo-->>ExpUC: incomes[]

    ExpUC->>MongoSL: GetStudentLoans()<br/>(filter monthYear เดือนปัจจุบัน)
    MongoSL-->>ExpUC: student loan list

    ExpUC->>ExpUC: NewPayrollCycle(incomes, loans)<br/>จับคู่ loan ด้วย BankAccountName<br/>หักกยศ. ตอนสร้างแถวไฟล์
    ExpUC->>File: WriteFile (CSV หรือ SAP)
    File-->>ExpUC: filename
    ExpUC->>ExpRepo: AddExport(filename, date)
    ExpUC-->>API: filename
    API-->>Web: Attachment (ไฟล์)
    Web-->>Admin: ดาวน์โหลดไฟล์ payroll
```

## Notes

- กยศ. sync เข้า Mongo **แยกจาก** add income — ไม่ถูกผูกตอน individual กรอก
- Add/Update อ่าน `user` แล้ว **snapshot** ลง `income`
- Export อ่านแค่ `income` + `studentloan`
- Sync script: `scripts/get_student_loan.go` (แหล่งภายนอก `slfrd.dsl.studentloan.or.th`)
- Export usecase: `business/usecases/export_income.go`
