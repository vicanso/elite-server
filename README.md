# elite

基于`elton`的脚手架，实现了数据校验、行为统计等功能。


## commit

feat：新功能（feature）

fix：修补bug

docs：文档（documentation）

style： 格式（不影响代码运行的变动）

refactor：重构（即不是新增功能，也不是修改bug的代码变动）

test：增加测试

chore：构建过程或辅助工具的变动

## 启动数据库

### postgres

```
docker pull postgres:alpine

docker run -d --restart=always \
  -v $PWD/elite:/var/lib/postgresql/data \
  -e POSTGRES_PASSWORD=A123456 \
  -p 5432:5432 \
  --name=elite \
  postgres:alpine

docker exec -it elite sh

psql -c "CREATE DATABASE elite;" -U postgres
psql -c "CREATE USER vicanso WITH PASSWORD 'A123456';" -U postgres
psql -c "GRANT ALL PRIVILEGES ON DATABASE elite to vicanso;" -U postgres
```

## redis

```
docker pull redis:alpine

docker run -d --restart=always \
  -p 6379:6379 \
  --name=redis \
  redis:alpine
```