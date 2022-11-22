/*
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-21 22:26:09
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-22 09:57:02
 * @FilePath: /film-server/Users/shixusheng/workspace/dbpool/pool/options.go
 * @Description: 客户端配置中心
 *
 * Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package pool

type OptionsGenerator interface {
	OptionsGenerator() (ClientOptions, error)
}
type ClientOptions interface {
	getOption()
}
