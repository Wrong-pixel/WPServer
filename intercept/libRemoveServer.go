/*
* 没啥好说的，把Server信息隐藏一下，基本操作了
*
* There is nothing to say, just hide the server information, basic operation
 */

package intercept

import (
	"github.com/gin-gonic/gin"
)

func RemoveServer(c *gin.Context) {
	c.Writer.Header().Del("Server")
	c.Writer.Header().Add("Server", "WPServer v1.3.1")
	c.Next()
}
