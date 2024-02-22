# Gin template
gin + gorm project template

## Feature Checklist
* [ ] Resources
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
datasource:
  driver: mysql
  url: tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
  username: 
  password: 
jwt:
  secret: test
  expire: 7
redis:
  addr: localhost:6379
  password:
  db: 0
```