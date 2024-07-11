/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package ctxdata

import "github.com/golang-jwt/jwt/v4"

const Identify = "imooc.com"

func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["identify"] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
