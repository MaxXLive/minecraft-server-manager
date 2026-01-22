# Minecraft Server Manager in Go

This is a Minecraft Server Manager written in Go that allows you to manage Minecraft servers by starting, stopping, and listing them. The server manager uses a combination of Go libraries and system commands (like `screen` and `java`) to control the Minecraft server processes. It also allows users to interact with servers, track their status, and handle logs in real-time.

## Features

- **Start Minecraft Servers**: Start a Minecraft server in the foreground or background with configurable memory.
- **Stop Minecraft Servers**: Stop the Minecraft server gracefully and monitor its status.
- **Restart Minecraft Servers**: Restart the server with a single command.
- **Health Checks**: Optionally verify server health via HTTP endpoint after startup.
- **Log Management**: Tail logs from the server in real-time and optionally write to a log file.
- **Server List**: Add, list, select, and remove servers from the manager.
- **Backup**: Stop server, commit world data to git, push to remote, and restart.
- **Self-Update**: Check for and install updates from GitHub releases.
- **Interactive Commands**: Spinner animation during server operations and real-time feedback.

## TODO
- Check for updates dependent on the server type (vanilla, paper, fabric)

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: [Go Installation Guide](https://golang.org/doc/install) (only needed for building from source)
- **Screen**: The Minecraft server manager relies on the `screen` command to run servers in the background.
- **Java**: You need Java installed to run Minecraft servers. Make sure Java is installed by running `java -version`.
- **Git**: Required for the backup command to push world data.
- **Curl**: Required for the self-update feature.

## Setup and Installation

### Download latest version (Example Linux AMD64)
```bash
wget https://github.com/MaxXLive/minecraft-server-manager/releases/download/v1.6.2/minecraft-server-manager_linux_arm64
chmod +x minecraft-server-manager_linux_arm64
```

### Move it to usr dir
```bash
mkdir -p /usr/minecraft-server-manager
mv minecraft-server-manager_linux_arm64 /usr/minecraft-server-manager/minecraft-server-manager
```

### Add to PATH to access it from everywhere
```bash
export PATH="$PATH:/usr/minecraft-server-manager"
function msm() { minecraft-server-manager "$@";}
```

### Clone and Build from Source
```bash
git clone https://github.com/maxxlive/minecraft-server-manager.git
cd minecraft-server-manager
chmod +x build.sh
./build.sh
```

## Commands

| Command | Description |
|---------|-------------|
| `start` | Start the minecraft server (foreground) |
| `start-bg` | Start the minecraft server in background |
| `stop` | Stop the minecraft server |
| `restart` / `r` | Restart the minecraft server |
| `kill` | Force kill the server (sends Ctrl+C first) |
| `console` / `c` | Attach to the minecraft server's console |
| `status` / `s` | Show the status of the minecraft server |
| `backup` / `b` | Stop server, create backup, push to git, restart |
| `list` | List saved servers in config |
| `select` | Select from servers in config |
| `add` | Add new server to config |
| `remove` | Remove a server from config |
| `logfile [enable\|disable]` | Show or set log file status |
| `help` | Show help message |
| `version` | Show version number |
| `check` | Check for updates |
| `update [--force]` | Update this app |

## Configuration

The configuration file `config.json` is stored in the same directory as the executable.

### Example config.json
```json
{
  "screen_name": "minecraft_server_",
  "servers": [
    {
      "id": "uuid-string",
      "name": "My Server",
      "max_ram": 4096,
      "jar_path": "/path/to/server.jar",
      "java_path": "java",
      "type": 0,
      "is_selected": true,
      "health_check_enabled": true
    }
  ],
  "log_file_enabled": true,
  "health_check_url": "http://localhost:8080/metrics/current"
}
```

### Configuration Options

| Option | Description |
|--------|-------------|
| `screen_name` | Prefix for screen session names |
| `log_file_enabled` | Enable/disable writing logs to `status.log` |
| `health_check_url` | URL to check for server health (global) |

### Server Options

| Option | Description |
|--------|-------------|
| `health_check_enabled` | Enable health checking for this server on start |

## Health Checks

When `health_check_enabled` is set to `true` for a server and `health_check_url` is configured:

1. After starting with `start-bg`, the manager will poll the health endpoint
2. If the server doesn't become healthy within 60 seconds, it kills the session and retries
3. Up to 5 retry attempts are made before giving up
4. The `status` command will also show health check status when enabled

This is useful with plugins like **AH Status API** that expose an HTTP endpoint when the server is fully loaded.

## Restart Script for In-Game /restart

For Spigot/Paper/Mohist servers, you can use the built-in `/restart` command by configuring `spigot.yml`:

```yaml
settings:
  restart-script: ./restart.sh
```

Create `restart.sh` in your server directory:
```bash
#!/bin/bash
/usr/minecraft-server-manager/minecraft-server-manager restart
```

## Add auto-start and auto-backup to crontab
```bash
@reboot /usr/minecraft-server-manager/minecraft-server-manager start-bg >> /tmp/crontab.log
0 6 * * * /usr/minecraft-server-manager/minecraft-server-manager backup >> /tmp/crontab.log
```