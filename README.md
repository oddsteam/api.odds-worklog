# Docker 

Build `docker build -t api-odds-worklog .`

Run `docker run --name api-odds-worklog -p 8080:8080 api-odds-worklog`


## API 

link  `http://worklog.odds.team`

GET `/api/userinfo`  mock ขึ้นมาเองครั

GET `/api/user` ดึงจากถังใน mongoครับ

POST `/api/insertUser`
      body:  `x-www-form-urlencoded`
     
keyvalue : 

      fullname
      email
      bankAccountName
      bankAccountNumber
      totalIncome
      submitDate
      cardNumber
DELETE `/api/user/:id`
       example: 
          DELETE  `http://worklog.odds.team/api/delete/5bb9798e106b940001443df0`