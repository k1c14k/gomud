# Go Mud
GoMud is a [Multi-user dungeon](https://en.wikipedia.org/wiki/Multi-user_dungeon) server (driver) written in [Go Programming Language](https://go.dev/). Key concepts come from [Arkadia LP MUD](https://arkadia.rpg.pl/) (MUD, world interaction) and [LPC scripting language](https://www.genesismud.org/lpc/lpc.pdf).

## Configuration

The server can be configured via a YAML file (`mud.yaml` by default) or via CLI flags.
When both are provided, CLI flags take precedence over the configuration file.

### CLI Flags

- `--config`: Path to the config file (default: `mud.yaml`)
- `--host`: Server host address (e.g., `0.0.0.0`)
- `--port`: Server port (e.g., `2323`)
- `--mudlib`: Path to the mudlib directory

**Example:**
```bash
go run ./cmd/server/main.go --port 8080 --host 127.0.0.1
```

### Configuration File (`mud.yaml`)

```yaml
server_config:
  host: "0.0.0.0"
  port: 2323
mudlib_config:
  mudlib_path: "mudlib/"
```

## Directory Structure

- `cmd/`: Command line applications (`compiler`, `server`)
- `docs/`: Project documentation
- `internal/`: Internal packages and game logic
- `mudlib/`: Default mudlib (scripts, locations, players)
