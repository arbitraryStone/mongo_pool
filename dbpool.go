/***
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-25 16:20:11
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-25 16:20:51
 * @FilePath: /dbpool/dbpool.go
 * @Description:数据库连接池对外接口
 * @
 * @Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package dbpool

import (
	"time"

	"github.com/arbitrarystone/dbpool/mongo"
	"github.com/arbitrarystone/dbpool/pool"
)

/***
 * @description:
 * @param {string} dbName 数据库名(mongo,redis,mysql,...)
 * @param {string} url 数据库地址
 * @param {int32} init 初始化连接数
 * @param {int32} cap 连接池容量
 * @param {int} mode pool.PoolGetModeLoose(连接池满后继续创建) pool.PoolGetModeStrict(连接池满后不继续创建)
 * @return {*}
 */
func NewPool(dbName string, url string, init int32, cap int32, mode int) (*pool.Pool, error) {
	var p *pool.Pool
	var err error
	switch dbName {
	case MongoDBName:
		p, err = newMongoDBPool(dbName, url, init, cap, mode)
	case RedisDBName:
	case MysqlDBName:
	}
	return p, err
}

func newMongoDBPool(dbName string, url string, init int32, cap int32, mode int) (*pool.Pool, error) {
	p, err := pool.NewPool(dbName, init, cap, mode)
	if err != nil {
		return nil, err
	}
	//创建mongo客户端生成器
	g := mongo.NewGenerator(url)
	//注册mongo客户端生成器
	p.RegisterClientGenerator(g)
	//初始化连接池
	if err := p.InitPool(); err != nil {
		return nil, err
	}
	return p, nil
}

func GetMongoClient(p *pool.Pool, timeout time.Duration) (*mongo.MongoClient, error) {
	v, err := p.Get(timeout)
	if err != nil {
		return nil, err
	}
	return v.(*mongo.MongoClient), nil
}

//后续更新
func GetRedisClient() {

}
func GetMysqlClient() {

}
