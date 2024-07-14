/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package config

import "github.com/zeromicro/go-zero/core/service"

type Config struct {
	service.ServiceConf

	ListenOn string
}
