#!/usr/bin/env bash
# Only works with Ubuntu 18.4


set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"
__base="$(basename ${__file} .sh)"
__root="$(cd "$(dirname "${__dir}")" && pwd)" # <-- change this as it depends on your app

arg1="${1:-}"


# Installing requirements for VirtualBox & Docker
sudo apt update
sudo apt -y upgrade
sudo apt -y dist-upgrade

sudo apt install -y build-essential linux-headers-$(uname -r) dkms apt-transport-https ca-certificates curl software-properties-common git make vim

echo "If you are running on virtual-box, install the virtual box tool now."
read -p "Press enter to continue"

# Installing Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo apt-key fingerprint 0EBFCD88
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt update
sudo apt -y upgrade
sudo apt -y install docker-ce docker-compose

# Allow Docker whitout sudo for current user
sudo usermod -aG docker $USER

#Download Software
git clone https://github.com/Sahzz/discord-raid-manager.git ~/

echo "Please add your api-key to the settings.json"
read -p "Press enter to continue"
vim ~/discord-raid-manager/discord_bot/settings.json

# Run Software
make start ~/discord-raid-manager/discord_bot
