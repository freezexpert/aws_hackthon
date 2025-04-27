## Prerequisites
[go](https://go.dev/)
## How to use
```go=
cd backend
go run main.go
```
Then the backend server will run at localhost:8888
## How to test with Postman
1. Install [postman](https://www.postman.com/) and create an account
2. Create 2 requests(one GET one POST) in postman
3. In GET request (it will load chat messages from database)
![](https://hackmd.io/_uploads/rkyVPwUnn.png)
4. In POST request(it will modify the data on database)
![](https://hackmd.io/_uploads/S1BIPPLn2.png)
In POST request, you will need to send the whole history json file, you can take a look at chat_history.json in the backend diretory
Then checkout the [database](https://cloud.mongodb.com/v2/64cf2c094620f341ba711440#/metrics/replicaSet/64cf2c303d37c7777ae8e45e/explorer/Project)