# dyn-ip-monitor

If your ISP has alloted you dynamic IP and you want to know how often it changes, and it what range, you are at the right place.

dyn-ip-monitor is a very simple service which fetches you public IP and stores it in mongo database.

## Overview

It stores the entries in `dynip` database and the collection is `ip_log`.

### Environment Variable

- `MONGODB_URI` - connection string to the database
- `INTERVAL` - sleep time between execution
- `TZ` - without this, it will fallback to UTC in your db entries

## How to use

**Note**: This app is in alpha, and right now I'm married to my setup. I'm running docker in swarm mode, and I have a central mongodb server. Contributions are welcomed, but right now I'll focus on MVP.

Right now, I pass following spec in portainer to create a stack:

```yaml
version: '3.8'

services:
  server:
    image: sntshk/dyn-ip-monitor
    environment:
      - INTERVAL=900  # 15 mins
      - TZ=Asia/Kolkata
```

## FAQ

1. **Why mongo if this could have been done with file based database?**

Because I started it on my existing homelab infrastructure. Pull requests are welcomed to add more databases.
