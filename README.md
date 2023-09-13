## 测试

[文档](https://gin-gonic.com/docs/testing/)

```go
import (
  "mangosteen/internal/router"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
  r := router.New()
  w := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/ping", nil)
  r.ServeHTTP(w, req)

  assert.Equal(t, 200, w.Code)
  assert.Equal(t, "pong1", w.Body.String())
}
```

## 数据库连接

### 安装数据库

```bash
# pg
docker run -d --name pg-for-go-mangosteen -e POSTGRES_USER=mangosteen -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=mangosteen_dev -e PGDATA=/var/lib/postgresql/data/pgdata -v pg-go-mangosteen-data:/var/lib/postgresql/data --network=network1 postgres:14

# mysql
docker run -d --network=network1 --name mysql-for-go-mangosteen -e MYSQL_DATABASE=mangosteen_dev -e MYSQL_USER=mangosteen -e MYSQL_PASSWORD=123456 -e MYSQL_ROOT_PASSWORD=123456 -v mysql-go-mangosteen-data:/var/lib/mysql -d mysql:8 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

### 连接 PG

[网站](https://www.connectionstrings.com/)

需要先安装 `PG` 驱动：`_ "github.com/lib/pq"`

```go
func Connect() {
  connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, passrord, dbname)
  defer db.Close()
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    log.Fatalln(err)
  }
  DB = db
  err = db.Ping()
  if err != nil {
    log.Fatalln(err)
  }
  log.Println("Connected to database")
}
```

### 连接 MySQL

需要先安装 `MySQL` 驱动：`_ "github.com/go-sql-driver/mysql"`

```go
func MySQLDatabase() {
  connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", passrord, "go-uccs-1", 3306, dbname)
  db, err := sql.Open("mysql", connStr)
  if err != nil {
    log.Fatalln(err)
  }
  DB = db
  err = db.Ping()
  if err != nil {
    log.Fatalln(err)
  }
  defer db.Close()
}
```

## 工具

### cobra 

使用 `go build .; ./mangosteen server`

```go
import (
	"github.com/spf13/cobra"
)

func Run() {
	rootCmd := &cobra.Command{
		Use: "mangosteen",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	srvCmd := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			RunServer()
		},
	}

	rootCmd.AddCommand(srvCmd)
	rootCmd.Execute()
}
```

### 安装 migrate

[文档](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md#further-reading)

```go
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

创建迁移文件：

```bash
migrate create -ext sql -dir config/migrations -seq create_users_table
```

运行迁移文件：

```bash
migrate -database "postgres://mangosteen:123456@go-mangosteen:5432/mangosteen_dev?sslmode=disable" -source "file://$(pwd)/config/migrations" up 
```

### 发送邮件

[使用说明](https://wx.mail.qq.com/list/readtemplate?name=app_intro.html#/agreement/authorizationCode)

```go
import (
  "gopkg.in/gomail.v2"
)
func Send() {
  m := gomail.NewMessage()
  m.SetHeader("From", "1500846601@qq.com")
  m.SetHeader("To", "1500846601@qq.com")
  m.SetHeader("Subject", "Hello!")
  m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
  d := gomail.NewDialer("smtp.qq.com", 465, "1500846601@qq.com", "xxxx")
  if err := d.DialAndSend(m); err != nil {
    panic(err)
  }
}
```

使用环境变量保存邮箱信息：

```go
var (
	EMAIL_STMP_HOST = os.Getenv("EMAIL_STMP_HOST")
	EMAIL_STMP_PORT = os.Getenv("EMAIL_STMP_PORT")
	EMAIL_USERNAME  = os.Getenv("EMAIL_USERNAME")
	EMAIL_PASSWORD  = os.Getenv("EMAIL_PASSWORD")
)

func Send() {
	m := gomail.NewMessage()
	m.SetHeader("From", "1500846601@qq.com")
	m.SetHeader("To", "1500846601@qq.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello 张凯君")
	port, err := strconv.Atoi(EMAIL_STMP_PORT)
	if err != nil {
		log.Fatalln(err)
	}
	d := gomail.NewDialer(EMAIL_STMP_HOST, port, EMAIL_USERNAME, EMAIL_PASSWORD)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
```

在本地的 `.zshrc` 文件中添加环境变量：

```bash
export EMAIL_STMP_HOST="smtp.qq.com"
export EMAIL_STMP_PORT=465
export EMAIL_USERNAME="1500846601@qq.com"
export EMAIL_PASSWORD="mihdblyiqyosgffa"
```

### vscode 调试

下面相当于运行 `go run . email` 命令：

```json
"configurations": [
  {
    "name": "Launch Package",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${workspaceFolder}",
    "args": ["email"]
  }
]
```

## 代码

1. 随机字符串：
    ```go
    var letterRuns = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    func init() {
      rand.New(rand.NewSource(time.Now().UnixNano()))
    }
    func RandStringRunes(n int) string {
      b := make([]rune, n)
      for i := range b {
        b[i] = letterRuns[rand.Intn(len(letterRuns))]
      }
      return string(b)
    }
    ```