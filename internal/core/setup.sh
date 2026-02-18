#!/bin/bash
exec > >(tee /var/log/user-data.log | logger -t user-data -s 2>/dev/console) 2>&1


# ── Config ────────────────────────────────────────────────────────────────────
MOUNT_POINT="/mnt/{{ .VolumeName }}"
DEVICE_PATH="{{ .DevicePath }}"
MCLOUD_MEM_RATIO=70


# ── Memory ────────────────────────────────────────────────────────────────────
TOTAL_MEM=$(free -m | awk '/^Mem:/{print $2}')
MCLOUD_MEM=$(( TOTAL_MEM * MCLOUD_MEM_RATIO / 100 ))
MEMORY_SETTING="${MCLOUD_MEM}M"


# ── Volume ────────────────────────────────────────────────────────────────────
sleep 20
mkdir -p $MOUNT_POINT

if ! blkid $DEVICE_PATH &>/dev/null; then
    mkfs.ext4 -F $DEVICE_PATH
fi

mount -o discard,defaults,noatime $DEVICE_PATH $MOUNT_POINT
chmod 777 $MOUNT_POINT


# ── Docker & Server ───────────────────────────────────────────────────────────
curl -fsSL https://get.docker.com | sh

docker run -d \
    --name mcloud \
    --restart always \
    -p 25565:25565 \
    -e EULA=TRUE \
    -e TYPE={{ .ServerType }} \
    -e VERSION={{ .ServerVersion }} \
    -e MEMORY=$MEMORY_SETTING \
    -e WHITE_LIST=TRUE \
    -e VIEW_DISTANCE=6 \
    -e ENABLE_RCON=true \
    -e RCON_PORT=25575 \
    -e RCON_PASSWORD="{{ .RconPassword }}" \
    -v "${MOUNT_POINT}":/data \
    {{ .DockerImage }}

ufw allow 22/tcp
ufw allow 25565/tcp

echo "mcloud setup complete!"
