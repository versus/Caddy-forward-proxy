#!/usr/bin/env bash

sudo yum update -y

sudo yum -y install wget
wget https://inspector-agent.amazonaws.com/linux/latest/install -P /tmp/
sudo bash /tmp/install

sudo yum install -y https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/linux_amd64/amazon-ssm-agent.rpm || true
sudo /usr/bin/systemctl enable amazon-ssm-agent
sudo /usr/bin/systemctl start amazon-ssm-agent

sudo groupadd  www-data || true
sudo adduser -g www-data  www-data || true