#!/usr/bin/env bash
set -x

sudo yum update -y
sudo yum install git -y
curl -LO https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xvzf go1.11.5.linux-amd64.tar.gz
mkdir -p ~/projects/{bin,pkg,src}
mkdir -p ~/projects/src/github.com/Sahzz/discord-raid-manager/discord_bot/
export PATH=$PATH:/usr/local/go/bin
export GOBIN="$HOME/projects/bin"
export GOPATH="$HOME/projects/src"
curl https://raw.githubusercontent.com/Sahzz/discord-raid-manager/master/discord_bot/main.go -o ~/projects/src/github.com/Sahzz/discord-raid-manager/discord_bot/discord_bot.go
cd ~/projects/src/github.com/Sahzz/discord-raid-manager/discord_bot/
go get -u
go get ./...
go install ./discord_bot.go
sudo ln -s $HOME/projects/bin/discord_bot /usr/local/bin/discord_bot
sudo cp ~ec2-user/discordbot.service /usr/lib/systemd/system/discordbot.service
sudo systemctl daemon-reload
sudo systemctl enable discordbot
