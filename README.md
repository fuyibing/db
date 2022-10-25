# DB

> 基于 `xorm` 构建的 `DB` 工具.



### Transaction

```go

func main(){
    ctx := trace.New()
    sess := db.Selector.GetMasterWithContext(ctx, name)

    db.TransactionWithSession(ctx, sess, func(){})
}




```


