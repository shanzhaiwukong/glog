# Glog
> 一个简单的go log包 支持日志分级 日志文件大小设置 自动按日期编号文件 暂停/继续收集等...

## 常规设置
``` golang
    // 每一秒清空一次缓存区
    log := glog.New(1e9, "./logpath/log.log")
    //记录info类型的log信息 参数为（...interface{}）
    log.Info("glog",time.Now().Unix())
    //用格式化的方式记录info信息
    log.Infof("Hello this itme is %s",time.Now().Format("2006/01/02 15:04:05"))
    // ... 支持其他记录类型 （log.Warn log.Error log.Debug log.Panic） 此处省略
    //也可以用另外的形式记录
    log.Log(glog.LevelInfo,"glog",time.Now().Unix())
    log.Logf(glog.LevelInfo,"Hello this itme is %s",time.Now().Format("2006/01/02 15:04:05"))
    //暂停收集日志
    log.Pasue()
    //继续收集日志
    log.Continue()
    //清空缓存区
    log.Flush()
```
## 全局设置
``` golang
    //配置日志文件大小
    glog.SetFileSize(1024*1e4) //每个文件10M
    //配置日志当前环境 除console环境外 均将日志记录到文件中
    glog.SetEnvironment(glog.EnvConsole)//控制台
    //设置日期格式化 默认 "2006/01/02 15:04:05"
    glog.SetTimeFormat("06/01/02 15:04")
```
## 数据库设置
`实现DataBase接口Exec`
```golang
    type SomeDataBase struct{

    }
    func (SomeDataBase) Exec(sql string, params ...interface{}) error {
        // todo some logic
        // ...
        fmt.Println("----", sql, params)
        return nil
    }
    //设置接口
    glog.SetDataBase(SomeDataBase{})
    //使用DataBase 会在返回结果前记录一条日志信息到缓存区
    glog.Db("select * from tab_name where id=?","id_id")
```
## Todo
`日志分析`