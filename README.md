# aip

Composable AI pipeline CLI for logs and text.

## Install

Build from source:

```sh
go build ./cmd/aip
```

Or build cross-platform artifacts:

```sh
./scripts/build.sh
```

## Quick start

Configure API access:

```sh
aip config wizard
```

Summarize input (streaming by default):

```sh
cat file.log | aip summary "summarize errors and next steps"
```

Normalize logs into signatures:

```sh
cat file.log | aip norm --profile postgres 
```

Cluster normalized signatures:

```sh
cat norm.jsonl | aip cluster
```

## Configuration

Config file: `~/.aip/config.toml`

```toml
base_url = "https://api.openai.com"
api_key = "sk-..."
model = "gpt-4o-mini"
```

Environment overrides:

```sh
AIP_BASE_URL=https://my-gateway.example.com \
AIP_API_KEY=... \
AIP_MODEL=... \
aip summary "summarize"
```

## Commands

Implemented:

- `summary <prompt> [file]` — single-pass LLM summary (streaming text by default)
- `norm [file]` — normalize logs into signatures (`--profile`, `--rules`, `--emit`)
- `cluster [file]` — simhash clustering for signatures (`--format`)
- `config` — manage config (`show/path/get/set/wizard`)
- `version`

Not yet implemented:

- `map`, `watch`, `sample`, `diagnose`

## Examples

Summarize top errors from PostgreSQL logs:

```sh
rg "ERROR|FATAL|PANIC" postgresql.log \
  | aip norm --profile postgres  \
  | aip cluster  \
  | aip summary "summarize root causes and suggested fixes"
```

Quickly scan one sample per cluster:

```sh
cat norm.jsonl | aip cluster --format sample
```


