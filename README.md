# **mcloud** 
### Build a persistent Minecraft world on your personal cloud and play anytime.


`mcloud` allows you to run Docker-based Minecraft servers in `your private cloud` using continuous data storage. Keep world data secure on block storage volumes while destroying compute instances when not in use to `reduce costs`.


## 📦 Installation
### for Developers
```bash

git clone https://github.com/naturesh/mcloud.git

cd mcloud

go mod tidy
go build -o mcloud cmd/mcloud/main.go

```

### for Users
```
Check the release
```


## Providers

#### `DigitalOcean`

1. set Digital Ocean TOKEN
```
export MCLOUD_DIGITALOCEAN_TOKEN="dop_xxxxxxxxxxxx..."
```
2. Register the ssh public key of the computer you want to use with the DigitalOcean.
3. Add your private key to the ssh-agent

## Usage 

1. Create new server configuration
```
mcloud init server.yaml
```

2. Start the server
```
mcloud up server.yaml
```

3. Add a user to the whitelist. Whitelist is turned on by default.
```
mcloud console server.yaml "whitelist add <playername>"
```

4. If you want to shut down the server (keep the world)
```
mcloud down server.yaml
```

5. If you want to know the status of the server
```
mcloud status server.yaml
```

# Disclaimer
Minecraft EULA: By using this tool, you automatically agree to Mojang's EULA.

Cost Warning: This tool creates real paid resources on your cloud provider (DigitalOcean, etc). You are responsible for all costs.

  Compute (Instance): Billing stops when you run mcloud down.

  Storage (Volume): Even after running down, the Volume remains to keep your world data. This incurs a small monthly fee until you manually delete it in the Cloud Console.
