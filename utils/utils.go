package utils

import(
	"time"
	"strconv"
)

var (
	uuidCounter int64
	uuidNode int
)

func SetUUIDNode(node int){
	uuidNode = node
}

func Atoi(str string)int{
	v, _ := strconv.Atoi(str)
	return v
}

/**
 * 1位符号
 * 31位时间戳(最大可表示到2038年)
 * 10位毫秒
 * 10位服务器ID(最大可表示1024)
 * 12位自增id(最大值是4096)
 * 共64位，每秒可生成400w条不同ID
 */
 func UUID() int64{
	uuidCounter++
	return ((time.Now().UnixNano() / 1000000) << 22) | int64((uuidNode & 0x7FFF) << 7) | (uuidCounter & 0x7F)
}