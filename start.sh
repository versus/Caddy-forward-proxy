#!/usr/bin/env bash

sudo pwd
sudo ls -la

sudo mv /proxy/proxy /usr/local/bin/proxy
sudo chown root:root /usr/local/bin/proxy
sudo chmod 755 /usr/local/bin/proxy
sudo ulimit -n 8192 /usr/local/bin/proxy
sudo cp /proxy/proxy.crt /etc/ssl/proxy.crt
sudo cp /proxy/key.key /etc/ssl/key.key
sudo chown www-data /etc/ssl/proxy.crt
sudo chown www-data /etc/ssl/key.key
sudo chmod 400 /etc/ssl/proxy.crt
sudo chmod 400 /etc/ssl/key.key
sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/proxy

sudo mv /proxy/proxy.service /etc/systemd/system/proxy.service

sudo /usr/bin/systemctl enable proxy.service
sudo /usr/bin/systemctl status proxy.service
sudo /usr/bin/systemctl start proxy.service