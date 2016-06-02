前后台协议：
MSG <msghead><msgcontent>\n

type msghead struct{
	index  uint16 //防止重复攻击？记下上次处理的index编号，如果下次还是发这个编号过来，就认为是非法的
	length uint32 //报文长度
    cmd   uint16 //请求类型
}

