# Analysis model

Analysis model ต่างกับ design model คือมันจะไม่ได้สะท้อนสิ่งที่อยู่ในซอฟต์แวร์ของเรา แต่จะสะท้อน
mental model ในหัว user เป็นหลัก

## Domain Payroll ทั่วไป

ใน domain ของ Payroll ทั่ว ๆ ไป model ที่ chatGPT generate มาให้ หน้าตาประมาณนี้

```mermaid
classDiagram

%% Core Domain
class Employee {
  employee_id
  first_name
  last_name
  hire_date
  salary
  tax_id
}

class BankAccount {
  bank_name
  account_number
  account_type
}

class TimeSheet {
  pay_period_start
  pay_period_end
  hours_worked
  overtime_hours
  approved_by
}

class Payroll {
  pay_period_start
  pay_period_end
  gross_pay
  net_pay
  status
}

class Payment {
  payment_date
  amount
  payment_method
  status
}

class BatchTransfer {
  batch_transfer_id
  transfer_date
  total_amount
  status
  bank_reference_id
}

class BankTransaction {
  transaction_id
  transaction_date
  status
  bank_message
}

class Deduction {
  deduction_type
  description
  amount
}

class Benefit {
  benefit_type
  description
  employer_contribution
  employee_contribution
}

class Tax {
  tax_type
  rate
  amount
}

class PayrollCycle {
  month
}

%% Relationships
Employee "1" -- "1" BankAccount
Employee "1" -- "many" TimeSheet
Employee "1" -- "many" Payroll
Employee "1" -- "many" Benefit
Payroll "1" -- "1" Payment
Payroll "1" -- "many" Deduction
Payroll "1" -- "many" Tax
Payment "1" -- "1" BankTransaction
Payment "many" -- "1" BatchTransfer
PayrollCycle "1" -- "1" BatchTransfer
BankTransaction "many" -- "1" BatchTransfer
PayrollCycle "1" -- "many" Payroll
```

พนักงาน 1 คน จะมี timesheet ได้หลายใบ (เช่นเดือนละใบ) แต่ละเดือน พนักงานก็จะได้รับค่าจ้าง
(Payroll) ของเดือนนั้น ๆ และอาจจะมีสวัสดิการ

ค่าจ้าง (Payroll) 1 อันจะมี การหักลบหนี้ (Deduction) และภาษีที่ (Tax) เกิดขึ้นจากรายได้ครั้งนั้น
เช่น

- ภาษีมูลค่าเพิ่ม (Value Added Tax, VAT) หรือ
- ภาษีหัก ณ ที่จ่าย (Witholding Tax) เป็นต้น

Payroll ของพนักงานทุกคนในเดือนนั้น ๆ จะถูกเอาไปรวมกันเป็น PayrollCycle ของเดือนนั้น ๆ
แล้วตอนทำจ่าย ก็ไปทำ BatchTransfer ที่ธนาคารเพื่อส่งรายการว่าต้องโอนให้บัญชี BankAccount ไหน
ยอดเท่าไหร่บ้าง เกิดเป็น BankTransaction สำหรับแต่ละรายการ ซึ่ง BankTransaction
ที่เกิดขึ้นจะสะท้อน Payment ของเดือนนั้น ๆ ของพนักงานแต่ละคน

## Odds Payroll

เพื่อให้มั่นใจว่าไม่มีใครในออดส์อยู่เพราะเสียดายสวัสดิการ พวกเราเลยไม่มีสวัสดิการ (Benefit)

Deduction หลัก ๆ ที่เรามี คือบางคนที่มีหนี้กองทุนกู้ยืมเพื่อการศึกษา (กยศ) บริษัทจะต้องหักหนี้ กยศ
ออกจากค่าจ้าง ก่อนจ่ายให้กับพนักงาน

ภาษีหัก ณ ที่จ่าย (Withholding Tax) เป็นหน้าที่ของบริษัทต้องหักและนำส่งสรรพากรแทนพนักงาน
ส่วนที่เหลือให้พนักงานไปจ่ายเพิ่ม หรือขอคืนเองตามภาษีเงินได้ประจำปีของพนักงานคนนั้น ๆ

สำหรับบุคคล/นิติบุคคลที่จดทะเบียนภาษีมูลค่าเพิ่ม (VAT) มีภาระที่จะต้องเก็บภาษีมูลค่าเพิ่มจาก Payroll
นี้เพื่อไปนำส่งภาษีมูลค่าเพิ่มต่อไป

แต่ละ PayrollCycle ของเราจะถูกแบ่งจ่ายออกเป็น 3 BatchTransaction คือ

1. สำหรับ Individual user
1. สำหรับ Corporate user (ที่จดนิติบุคคล) และมี BankAccount ของ TTB เอง และ
1. Coporate user ที่ใช้บัญชีต่างธนาคาร

ที่ต้องแยก Corporate ต่างธนาคารออกไป เพราะการโอนภายใน TTB จะ effective ทันที
แต่ถ้าโอนต่างธนาคารจะต้องรอ 2 business days ถึงจะมีผล เลยต้องแยก file BatchTransaction
ออกจากกันเป็นคนละ file

หน้าตา analysis model สำหรับออดส์เลยเป็นแบบนี้ (ไม่มี benefit และ 1 PayrollCycle จะมี 3
BatchTransactions)

```mermaid
classDiagram

%% Core Domain
class Employee {
  employee_id
  first_name
  last_name
  hire_date
  salary
  tax_id
}

class BankAccount {
  bank_name
  account_number
  account_type
}

class TimeSheet {
  pay_period_start
  pay_period_end
  hours_worked
  overtime_hours
  approved_by
}

class Payroll {
  pay_period_start
  pay_period_end
  gross_pay
  net_pay
  status
}

class Payment {
  payment_date
  amount
  payment_method
  status
}

class BatchTransfer {
  batch_transfer_id
  transfer_date
  total_amount
  status
  bank_reference_id
}

class BankTransaction {
  transaction_id
  transaction_date
  status
  bank_message
}

class Deduction {
  deduction_type
  description
  amount
}

class Tax {
  tax_type
  rate
  amount
}

class PayrollCycle {
  month
}

%% Relationships
Employee "1" -- "1" BankAccount
Employee "1" -- "many" TimeSheet
Employee "1" -- "many" Payroll
Payroll "1" -- "1" Payment
Payroll "1" -- "many" Deduction
Payroll "1" -- "many" Tax
Payment "1" -- "1" BankTransaction
Payment "many" -- "1" BatchTransfer
PayrollCycle "1" -- "3" BatchTransfer
BankTransaction "many" -- "1" BatchTransfer
PayrollCycle "1" -- "many" Payroll
```

## Process ในการได้ file PayrollCycle มา

```mermaid
gantt
    title สมมติว่าเป็นเดือน Oct 2025 ของออดส์
    dateFormat  YYYY-MM-DD
    section Employee
    กรอก income           :a1, 2025-10-01, 24d
    worklog กรอกครบแล้ว :milestone, 0d
    section Admin
    บันทึกหนี้ตามที่ กยศ แจ้งมา: 2025-10-5, 1d
    Export payroll cycle      :e1, 2025-10-25, 0d
    ทำ file batch transation ส่งธนาคาร    :t1, after e1, 0d
    section Bank 
    โอนภายในธนาคาร: 1d
    โอนต่างธนาคาร:2025-11-02, 1d
```
