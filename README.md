# Gin template
gin + gorm project template

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![GitHub release](https://img.shields.io/github/tag/SShnoodles/gin-admin-template.svg?label=release)](https://github.com/SShnoodles/gin-admin-template/releases)
[![LICENSE](https://img.shields.io/github/license/SShnoodles/gin-admin-template.svg)](https://github.com/SShnoodles/gin-admin-template/blob/main/LICENSE)

## Frontend
[Example Frontend(zh-CN)](https://github.com/SShnoodles/gin-admin-frontend-template)

## Features
* **Multi** organizational permission design
* **Restful** API design
* Database repository **Gorm**
* Log repository **zap**
* Base authentication **JWT**
* Api Docs repository **Swagger**
* Return result message based on **i18n**
* Configuration file repository **viper**

## Interface List
* [x] Resources
* [x] Menus
* [x] Organizations
* [x] Roles
* [x] Users

## Configuration File
config.yml

```yaml
server:
  port: 8080
logging:
  level: info
  file:
    name: app.log
    path: logs
# en/zh-CN
language: en
datasource:
  driver: mysql
  url: tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
  username: 
  password: 
jwt:
  secret: test
  expire: 7
verification:
  resourceEnabled: false
redis:
  addr: localhost:6379
  password:
  db: 0
```

## Swagger
### install
```shell
go install github.com/swaggo/swag/cmd/swag@latest
```
### init or update
```shell
swag init
```
### UI URL
```shell
http://localhost:8080/swagger/index.html
```
