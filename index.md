# 概述

Rigger是一个项目环境搭建工具，用于解决多人开发时的统一环境配置问题

本项目的示例都是用的项目中的[demo](https://github.com/ligang1109/rigger/tree/master/demo)目录中的内容

# 部署

请先安装[dep](https://golang.github.io/dep/)

```
git clone git@github.com:ligang1109/rigger.git
cd rigger
./deploy/deploy.sh host1 host2 ...
```

这会将rigger工具安装至目标主机的/usr/local/bin下

## 项目中使用

首先，需要填写指导rigger执行所需要的配置文件，详见：[Rigger配置文件]()

示例运行：

```
rigger -rconfDir=/home/ligang/devspace/rigger/demo/conf/rigger/ logLevel=1 prjHome=/home/ligang/devspace/rigger/demo
```

示例中的参数说明：

- -rconfDir：必选参数，放置rigger配置文件的目录，请用绝对路径

- logLevel: 可选，值见[logLevel](https://github.com/goinbox/golog/blob/master/base.go)

- prjHome：可选，和项目相关，这里是因为demo的配置文件[var.json]()中会用到
