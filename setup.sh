#!/bin/sh

sudo rm -rf /usr/local/go
sudo curl -O https://dl.google.com/go/go1.20.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bash_profile
echo "source ~/.bash_profile | re-login"