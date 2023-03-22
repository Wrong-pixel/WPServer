# WPServer(废弃项目)

**基于yaml插件实现软件WAF的反向代理网关**

用了nginx很久了，苦于没找到什么好用的软件waf，就想着自己模仿nginx写一个，并且能够实现基于插件灵活配置的软件waf，于是就有了这个项目。
本来是想用协程和连接池来提高并发的，但是golang本身并发支持就挺高的，况且我也不会写（摆），暂时搁置吧，目前在1000请求/秒的并发下丢包率在0.3%
<img width="980" alt="image" src="https://user-images.githubusercontent.com/43137902/208852671-f5f2614b-afd5-42b4-9f44-73749d74e13b.png">
完成了SQL注入、log4j、fastjson、xss等漏洞的简单防御
<img width="1417" alt="image" src="https://user-images.githubusercontent.com/43137902/210288213-9b12c711-c83c-4f44-8511-4e5fe67c9173.png">
拦截后跳转页面（抄了长亭waf的页面）
<img width="1162" alt="image" src="https://user-images.githubusercontent.com/43137902/210288277-e0eea361-d6c3-43a0-bdd9-ea40b399f3b5.png">


原理其实就是gin的中间件（为什么不用net/http？），目前已经实现的功能有：

> 1. 反向代理
> 2. 负载均衡
> 3. Hearder替换
> 4. 漏洞插件

计划中的功能：

> 1. 443端口监听，证书管理，SSL(这里卡主了，还没想好怎么把证书卸载再装载)
> 2. web管理界面（也需要微服务？或者其他方式，可能要大改，暂时搁置）

项目已经处于一个初步可用的状态，希望各位师傅多提意见，包括新功能，新想法

ps. 啥时候给我发一个前端？
