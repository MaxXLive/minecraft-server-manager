# Minecraft Server Manager in Go

This is a Minecraft Server Manager written in Go that allows you to manage Minecraft servers by starting, stopping, and listing them. The server manager uses a combination of Go libraries and system commands (like `screen` and `java`) to control the Minecraft server processes. It also allows users to interact with servers, track their status, and handle logs in real-time.

## Features

- **Start Minecraft Servers**: Start a Minecraft server in the background with configurable memory.
- **Stop Minecraft Servers**: Stop the Minecraft server and monitor its status.
- **Log Management**: Tail logs from the server in real-time and display them to the user.
- **Server List**: Add, list, and remove servers from the manager.
- **Interactive Commands**: Spinner animation during server startup and real-time feedback.

## TODO
- Check for updates dependent on the server type (vanilla, paper, fabric)
- Automatic backup option (notify, stop, git push, start)
- Start the server in the background

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: [Go Installation Guide](https://golang.org/doc/install)
- **Screen**: The Minecraft server manager relies on the `screen` command to run servers in the background.
- **Java**: You need Java installed to run Minecraft servers. Make sure Java is installed by running `java -version`.

## Setup and Installation

### Clone the Repository

```bash
git clone https://github.com/maxxlive/minecraft-server-manager.git
cd minecraft-server-manager
chmod +x build.sh
./build.sh
