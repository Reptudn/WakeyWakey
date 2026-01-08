# WakeyWakey

A Discord bot to send Wake-on-LAN (WoL) packets to your PCs from Discord slash commands. Register your machines, manage them with aliases, and wake them up instantly — all from Discord.

> Credits to my homie [@Tubinex](https://github.com/Tubinex) for the idea and the name to this repo and bot <3

## Features

- **Register PCs** — Store PC aliases with their MAC addresses
- **Wake on Command** — Send WoL packets with autocomplete suggestions
- **Manage Devices** — Register, unregister, and list your PCs
- **User-Isolated** — Each user manages their own PC registry
- **Ephemeral Responses** — Only you see command feedback
- **SQLite Backend** — Lightweight, persistent storage
- **Docker Ready** — Ships with Dockerfile and docker-compose

## Prerequisites

- Go 1.25+
- Discord Bot Token (get from [Discord Developer Portal](https://discord.com/developers/applications))
- Discord Server ID (Guild ID)

## Setup

This Project is currently setup that is works with only Docker really.
But if you still want to execute it without Docker you would have to set the env vars manually first before running the bot.

### 1. Environment Variables

Create a `.env` file or copy and edit the `.env.example` and rename it to `.env`:

```bash
export DISCORD_BOT_TOKEN="your_bot_token_here"
export DISCORD_GUILD_ID="your_guild_id_here"
export SELF_USER_ID="your_discord_user_id_here"  # Optional: restrict bot to your commands only
```

### 2. Run with Docker

```bash
docker-compose up --build
```

Data persists in the `./data` volume.

## Usage

### `/register <alias> <mac-address>`

Register a new PC.

```
/register my-gaming-pc AA:BB:CC:DD:EE:FF
```

### `/wake <alias>`

Send a WoL packet to wake a registered PC. Supports autocomplete to list your aliases.

```
/wake my-gaming-pc
```

### `/unregister <alias>`

Remove a PC from your registry.

```
/unregister my-gaming-pc
```

### `/list`

Lists all you added Devices.

```
/list
```

## Dependencies

- `github.com/bwmarrin/discordgo` — Discord API bindings
- `modernc.org/sqlite` — Pure Go SQLite driver (no CGO)

## Security

- Bot only responds to the `SELF_USER_ID` if set (optional restriction)
- All commands are ephemeral (invisible to other users)
- MAC addresses stored locally, never exposed

## Troubleshooting

**Bot not appearing online:**

- Normal. Discord bots don't show "online" status. Your bot is working if commands execute.

**"Database is locked" errors:**

- Handled by busy timeout pragma. Ensure only one bot instance is running.

**Commands not showing:**

- Verify `DISCORD_BOT_TOKEN` and `DISCORD_GUILD_ID` are correct.
- Reinvite the bot with proper scopes: `applications.commands`, `send_messages`.

**WoL packets not received:**

- Verify MAC address format (with colons: `AA:BB:CC:DD:EE:FF`).
- Ensure target PC has WoL enabled in BIOS.
- Some networks block broadcast packets—check firewall rules.
