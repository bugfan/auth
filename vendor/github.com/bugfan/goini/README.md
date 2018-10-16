# goini
golang 读取配置文件的工具包
根据配置文件的key读取value

# 简易使用方法
1. go get github.com/bugfan/goini
2. goini.LoadConfig("***/test.conf")
3. goini.Config.GetString("key")
4. goini.Config.GetInt64("key")


# 使用须知
1. 初始化设置配置（这一句：goini.LoadConfig("***/test.conf")）路径设置不对会报错，需要明确指定配置文件路径
2. 读不到数据会返回默认值
3. 支持直接获取字符串 和int类型
