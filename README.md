# super-cp

`super-cp` 是一个基于 Golang 实现的静态资源发布工具。

他针对将静态资源部署到 s3 的场景，提供了缓存控制和发布顺序管理能力

```plaintext
$ super-cp --help

Copy files to object storage as a static website

Usage:
  super-cp [flags] [job-names...]

Flags:
  -a, --analyze           print analyze output (default true)
  -j, --concurrency int   concurrency level (default 8)
  -c, --config string     config file path (default ".super-cp.yml")
      --dry-run           print all actions without actually do
  -h, --help              help for super-cp
  -v, --verbose           print verbose output
```

## 核心概念

### 1. jobs

表示不同的传输任务，使用命令行的 args 指定实际使用哪一个, 默认是 `default`。

```yaml
jobs:
  test:
  prod:
```

### 2. job.source

文件源的过滤规则，使用 Unix Shell 风格的 glob 语法，支持 `**` 通配符

可以使用 strip 去除一部分前缀，下例的 `dist/index.html` 会被传输到目标存储的 `index.html` 中。

```yaml
source: # 源文件(可以有多个)
  pattern: "dist/**/*" # 匹配 dist 目录下的所有文件
  strip: "dist" # 去掉 dist 前缀
```

### 3. job.dist

目标存储 DSN

```yaml
dist:
  type: "s3"
  dsn: "s3://-/mybucket/mypath/"
```

### 4. job.rules

对待发布文件应用规则，参考 `.example.yml` 中的示例。

```yaml
rules:
    # glob 匹配
  - pattern: "{icons/**,favicon.ico}"
    
    # 添加 header
    headers:
      "Cache-Control": "public, max-age=300"

    # 标记文件为排除
    exclude: true

    # 设置上传顺序，默认是 0，先小后大
    index: 1
```

## Development

```bash
go build -o super-cp main.go
./super-cp -config .super-cp.yml test
```
