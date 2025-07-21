# goframe-ex
## 简介
* 基于Go Frame 对一些工具方法进行扩展封装

  


## 内置功能

1. **ecache**：主要为了后续db的cache使用，添加leveldb、file 类型的cache;***注意 ：file的实现用的hotgo的，部分接口实现为完成***。

2. **econf**：配置文件添加hocon格式文本；基于pongo,实现特定配置加密读取；根据gf.gcfg.file配置，额外读取配置文件。***注意 ：如果要使用，需要优先最优先初始化此配置，否则配置不生效***。

3. **ecron**：定时任务添加前置处理函数，方便设置有效期或者其他的，统一管理相应的处理服务；将goframe定时任务的next方法提出来，方便解析表达式后续几次时间。

4. **egoja**：对goja的一些处理进行封装，方便goja中引用其他函数或内置函数。

5. **edb**：基于gdb封装一些自己需要的函数方法。

   ```go
   // 1. ${} #{} 占位获取相应的处理数据
   var ccc, err = libdb.DB("h").Select(gctx.New(), `
       select  * from [TB_HR_HRJL] where  deptCode = #{deptCode}
       ?{ typeof lsh !== 'undefined' && lsh!=null  && lsh !=''  , and lsh= #{ lsh } }
   `, map[string]any{"deptCode": "MZ1"})
   
   	fmt.Println(ccc, err)
   // 2. 多结果集的一次性返回，主要是执行存储，存储中可能有多个结果集，能都收到
   // 3. 非model模式的sql语句缓存，实现比较简单，后期可能修改
   
   ```

   

6. **eswagger**：knife4j 风格前端集成。

7. **esoap**：根据ghttp，实现的webservice 服务端。应该是支持大部分其他语言调用的。

8. **einit**：程序有效期的简单设置。

   



## 后续计划
* edb添加其他数据库驱动
* edb 根据 gmeta 、gstructs 实现初始化自动建表
* ecache 部分数据库实现方法完善
* egoja 完善添加内置goframe其他函数
