# WPServer

**基于yaml插件实现软件WAF的反向代理网关**

用了nginx很久了，苦于没找到什么好用的软件waf，就想着自己模仿nginx写一个，并且能够实现基于插件灵活配置的软件waf，于是就有了这个项目。

原理其实就是gin的中间件（为什么不用net/http？），目前已经实现的功能有：

> 1. 反向代理
> 2. 负载均衡
> 3. Hearder替换

计划中的功能：

> 1. 证书管理，SSL转发
> 2. 基于插件的软件waf
> 3. web管理界面

项目已经处于一个初步可用的状态，希望各位师傅多提意见，包括新功能，新想法

ps. 啥时候给我发一个前端？
