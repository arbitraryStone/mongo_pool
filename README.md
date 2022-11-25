## 数据库连接池 
### Usage
```go
```
#### mongo连接池Usage
```go
func main() {
	p, err := dbpool.NewPool("mongo", "mongodb://localhost:27017", 1, 10, pool.PoolGetModeStrict)
	if err != nil {
		panic(err)
	}
    defer p.Close()
	client, err := dbpool.GetMongoClient(p, 1*time.Second)
    defer client.Close()
	if err != nil {
        fmt.Printf("get mongo client failed:%v",err)
        return
	}
	_, err = client.GetCount("FilmDB", "area")
	if err != nil {
		fmt.Printf("get count error:%v", err)
		return
	}
}
```