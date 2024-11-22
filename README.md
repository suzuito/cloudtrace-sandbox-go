# cloudtrace-sandbox-go

CloudTraceと戯れるレポジトリ。

CloudTraceとは[Google Cloudが提供するTraceサービス](https://cloud.google.com/trace/docs/overview)。

Go言語で書かれたプロダクトにおいては、[opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go)を介してCloudTraceへデータを送信できる。

Google Cloud以外の環境からCloudTraceに対してデータを送信する場合、認証を通す必要がある(いつものあれ。`GOOGLE_APPLICATION_CREDENTIALS`)。

```bash
go run cmd/server001/main.go
```
