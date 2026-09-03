# Examples

Each example is a standalone Go program. Set your LunarCrush API key before
running one:

```powershell
$env:LUNARCRUSH_API_KEY = "your-api-key"
go run ./examples/basic-topic
```

Available examples:

- [`basic-topic`](basic-topic) — retrieve the social summary for Bitcoin.
- [`coin-list`](coin-list) — request and print a filtered, sorted coin list.
- [`error-handling`](error-handling) — configure retries and distinguish common
  API errors.
