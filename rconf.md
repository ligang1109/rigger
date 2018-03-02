# Rigger配置文件

## 概述

Rigger的执行过程，通过下面提到的一些配置文件来指定，请将这些配置文件统一放置在一个目录下，如：

```
ligang@vm-centos7 ~ $ ls -l /home/ligang/tmp/misc/rigger/demo/conf/rigger
total 12
-rw-r--r--. 1 ligang wheel  275 Mar  2 14:24 action.json
drwxr-xr-x. 2 ligang wheel   61 Mar  2 14:24 tpl
-rw-r--r--. 1 ligang wheel  337 Mar  2 14:24 tpl.json
-rw-r--r--. 1 ligang wheel 1097 Mar  2 14:24 var.json
```

## var.json

每个人的开发环境都会有不同的地方，例如：路径、port等，可以通过在这个配置文件中定义一些变量，然后在模版中使用的方式来解决这些不同，从而做到配置的统一

### 示例说明

```
{
    "IS_DEV": "true",          //IS_DEV即是定义的变量名，可以在配置文件中及模版文件中用到
    "USER": "${USER}",         //${USER}是使用变量的语法：${varName}
    "PRJ_NAME": "rdemo",
    "PRJ_HOME": "__ARG__(prjHome)", //这里使用rigger命令行执行时用户指定的参数值prjHome

    "PRJ_CONF_ROOT": "${PRJ_HOME}/conf",
    "TPL_CONF_ROOT": "${PRJ_CONF_ROOT}/rigger/tpl",
    "NGX_DATA_ROOT": {                  //每个人的开发环境都可能不同，可以通过这种形式区分
        "ligang": "/data/nginx",        //这里是指用户ligang执行时的值
        "default": "/usr/local/nginx"   //没有特殊定义，则用户执行时的值
    },

    "FRONT_DOMAIN": "${USER}.rdemo.com",
    "FRONT_ACCESS_LOG": "${FRONT_DOMAIN}.log",
    "FRONT_ERROR_LOG": "${FRONT_DOMAIN}.error.log",
    "FRONT_HTTP_CONF_TPL": "${TPL_CONF_ROOT}/tpl_front_http.conf",
    "FRONT_HTTP_CONF_DST": "${PRJ_CONF_ROOT}/http/${USER}_front_http.conf.ngx",
    "FRONT_HTTP_CONF_LN": "${NGX_DATA_ROOT}/conf/include/${FRONT_DOMAIN}.conf",

    "ACCESS_LOG_BUFFER": "1k",
    "NGX_PORT": "80",
    "GO_PORT": "__MATH__(6000+${UID})",    //这里的值通过数学运算获得

    "SERVER_CONF_TPL": "${TPL_CONF_ROOT}/tpl_server_conf.json",
    "SERVER_CONF_DST": "${PRJ_CONF_ROOT}/server/${USER}_server_conf.json",
    "SERVER_CONF_LN": "${PRJ_CONF_ROOT}/server_conf.json",

    "NGX_EXEC_PREFIX": {      //依旧是不同用户环境的区分
        "ligang": "/usr/local/bin/dexec nginx",
        "default": "sudo /usr/local/nginx/sbin/nginx"
    }
}
```

### 重点说明

- 定义的值，请都使用字符串的方式，如"true"，"1"
- 定义的值如果是路径，请使用绝对路径
- ${varName}取值时，如果没有在本配置中定义，会尝试从用户的环境变量中获取该值
- `__ARG__(argName)`，这个值是rigger执行时用户通过argName=value这种方式指定的
- `__MATH__(lvalue[+-*/]rvalue)`，这个值的结果是lvalue和rvalue的运算结果，目前支持+-*/

## tpl.json

定义模版文件的解析，模版文件使用的demo中的tpl：

[tpl_server_conf.json](https://github.com/ligang1109/rigger/blob/master/demo/conf/rigger/tpl/tpl_server_conf.json)

[tpl_front_http.conf](https://github.com/ligang1109/rigger/blob/master/demo/conf/rigger/tpl/tpl_front_http.conf)

### 示例说明

```
{
    "server_conf": {                        //server_conf只是一个标识，没有特别的意义
        "tpl": "${SERVER_CONF_TPL}",        //模版文件路径
        "dst": "${SERVER_CONF_DST}",        //通过模版生成的文件路径
        "ln": "${SERVER_CONF_LN}",          //如需对生成文件做软链，这里设置软链路径，不需要则不设置
        "sudo": false                       //做软链时是否需要sudo
    },
    "front_http_conf": {
        "tpl": "${FRONT_HTTP_CONF_TPL}",
        "dst": "${FRONT_HTTP_CONF_DST}",
        "ln": "${FRONT_HTTP_CONF_LN}",
        "sudo": true 
    }
}
```

## action.json

通过模版生成项目所需要的配置文件后，通常会还需要有些初始化类型的操作，如临时目录的创建、webserver的重启等，这个配置文件就用于定义这些方面的行为

### 示例说明

```
{
    "mkdir": [    //这里定义目录的创建
    {
        "dir": "${PRJ_HOME}/tmp",      //要创建的目录名称
        "mask": "777",                 //目录权限
        "sudo": false                  //是否sudo创建
    },
    {
        "dir": "${PRJ_HOME}/logs/front",
        "mask": "777",
        "sudo": false
    },
    {
        "dir": "${PRJ_HOME}/logs/task",
        "mask": "755",
        "sudo": false
    }
    ],
    "exec": [       //这里定义哪些命令会被执行
        "${NGX_EXEC_PREFIX} -s reload"      //要执行的命令
    ]
}
```
