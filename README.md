# What is this

This is the cli Keycloak client.
In this first release:
- Retrieve the federal user ID from Keycloak

More funcs will be add in the near furture: 
- interactive cli.
- restAPI(Fiber).

This CLI tool does not support interactive mode at the moment, so it strictly requires input and hardcodes the output for the first value only.

# How to 
...
```
├── Sample
│   ├── get.json
│   ├── get1.json
│   ├── subtask.json
│   ├── task.json
│   └── task2.json
├── cmd
│   └── root.go
├── credentials.json
├── handler
│   └── handler.go
├── main.go
├── pkgs
│   ├── google
│   │   └── google.go
│   ├── jira
│   │   └── jira.go
│   └── keycloak
│       └── keycloak.go
```

# API sử dụng:
## Google Chat API
Auth: Sử dụng Service account, yêu cầu đặt tên: `credentials.json`

Mỗi tin gửi đi sẽ được  đánh `MessageId` (chỉ duy nhất) và `ThreadId` dựa theo Jira Issue Key, với Issue Type là `Sub-DevOps`, `MessageId` sẽ đánh thêm `ChangeLogId`

Vì tin nhắn edit sẽ không mention lại user, nên các changelog liên quan tới task sẽ được gửi mới trong tin nhắn cha.

# Trouble shoot


