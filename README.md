# GoSearch

演示地址：[http://s.install.ren](http://s.install.ren/) 
>感谢三丰云的免费云服务器，虽然要5天延期一下，但是使用还是很丝滑的

## 示意图
![](html/demo.png)

## 项目介绍

聚合各大搜索引擎结果，关键是秒开，UI清爽，支持**暗夜模式**，还支持**自定义权重排序**，哈哈！

**主要功能有：**

- 支持搜索引擎配置权重
- 支持域名配置权重
- 搜索引擎结果渲染优化

[里程碑计划](note/roadmap.md)

## 软件架构

采用go语言开发，前端和后端都在一个项目，转发请求到搜索引擎并对结果进行裁剪、聚合、展示。

**依赖组件：**

- goquery
- yaml.v2
- bootstrap5（因为页面简单，不需要独立部署前端服务）

[排序算法实现思路](note/algorithm.md)

## 使用说明

### 启动

#### Dock启动
```
docker run -d -p 80:80 zyhui98/gosearch:v1.0
```
访问地址：[http://127.0.0.1](http://127.0.0.1)


#### 本地启动
```
go run main.go 
```

### 配置文件

路径：``configs/config.yml``

```
server:
  debug: false
  port: 80

search:
  - name: Baidu
    domain: www.baidu.com
    weight: 1 #搜索引擎权重因子
    positionWeight: 1 #搜索引擎自然排序权重因子
    score: 0 #搜索引擎设置的附加得分
    enable: false #是否开启
  - name: Bing
    domain: cn.bing.com
    weight: 1
    positionWeight: 1
    score: 10
    enable: true
  - name: Google
    domain: www.google.com
    weight: 1
    positionWeight: 1
    score: 10
    enable: false
  - name: 微信公众号
    domain: weixin.sogou.com
    weight: 1
    positionWeight: 1
    score: 10
    enable: true

site:
  - domain: www.csdn.com
    weight: 1 #域名权重因子
    score: 0 #网站域名设置的附加得分
  - domain: zhuanlan.zhihu.com
    weight: 1
    score: 0
  - domain: www.yuanbiguo.com
    weight: 1
    score: 0
  - domain: juejin.cn
    weight: 1
    score: 0

```



## 参与贡献
- Fork 本项目
- 新建 Feat_xxx 分支
- 提交代码
- 新建 Pull Request
