# super-cp

### 能力 & 原理

`super-cp` 是一个基于 Golang 实现的静态资源发布工具。super-copy 配合 yml 规则定义文件（environment、source、dist、rule）一起实现了对前端项目的自动化部署到 S3。原理是 super-cp 读取当前需要部署的项目根目录下的 yml 配置，先通过 source 源规则过滤需要上传的文件列表，将这个文件列表处理成的 S3 key 形式，再对每个需要上传的文件（target file）进行遍历所有 rules，给 target file 预备好 metadata，最后上传到 dist 定义的 S3 bucket 中，这一步需要把 target file + metadata 一起处理并上传。

### 三个重要概念

#### environment（环境）

表示不同环境，在使用工具的时候使用 -e 进行指定

```yaml
environments:
  test:
  
  prod:
```

#### source（本地源）

表示多个源文件的过滤规则，下面是两个源规则（`dist/**/*`表示递归 dist 下面的所有文件，`**/*.html`表示当前前端项目的所有 .html 文件），目前 opetions 只支持 dot 和 noCase。

```yaml
source: # 源文件(可以有多个)
  - pattern: "dist/**/*" # 匹配 dist 目录下的所有文件
    strip: "dist" # 去掉 dist 前缀
  - pattern:
      glob: "**/*.html" # 匹配当前.所有 html 文件
      options:
        dot: true # 是否上传隐藏文件
        noCase: true # 是否区分大小写
```

#### dist（远程目标存储）

需要上传的目标存储配置

```yaml
dist:
  type: "@s3"
  dsn: "s3://appid:appsecret@10.69.64.54:19000/super-cp/test/?use-path-style=true&protocol=http"
```

#### rule（文件处理规则）

对通过了 source 过滤需要上传的文件进行规则处理，包括加 Header 等操作。

```yaml
rules:
  - pattern: "*.html" # 匹配所有 html 文件
    headers:
      "Cache-Control": "public, max-age=60" # 设置缓存控制，60 秒后过期
      "Content-Disposition": "inline" # 设置 Content-Disposition 为 inline
    auto-mime-type: true # 自动根据文件名设置 Content-Type
```

### 测试

```bash
$ go build -o super-cp main.go
# 测试 test 环境下的规则
$ ./super-cp -config .super-cp.yml -e test
```