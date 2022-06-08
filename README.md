# DB

> package `github.com/fuyibing/db/v2`.

----

### 支持多实例

> 允许在配置文件 `config/db.yaml` 中配置多个数据库实例, 接受主从配置;
> 在 `dsn` 片段中, 第1个为 `Master`, 第 2 个开始为 `Slave`.
> 默认实例(即database下的第1个名称)务必使用 `db` 命名.

```yaml
# config/db.yaml
databases:
  db:
    driver: "mysql"
    dsn:
      - "user11:pass11@tcp(host11:port)/name1?charset=utf8"
      - "user11:pass11@tcp(host12:port)/name1?charset=utf8"
```

### 支持动态实例

> 设置实例

```text
db.Config.SetDatabase("key", &db.Database{
    Driver: "mysql",
    Dsn: []string{
        "user1:pass1@tcp(host1:port)/name?charset=utf8",
        "user1:pass1@tcp(host2:port)/name?charset=utf8",
    },
})
```

> 使用主库连接

```text
session := db.Manager.GetEngine("key").Master().NewSession()
session.QueryString("SELECT 1")
```

> 使用从库连接

```text
session := db.Manager.GetEngine("key").Slave().NewSession()
session.QueryString("SELECT 1")
```
