# Gin template
gin + gorm 项目模版

## 功能清单
* [ ] 机构
* [ ] 菜单
* [ ] 角色
* [ ] 用户

## 配置文件
config.yml

```yaml
server:
  port: 8080
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