#!/bin/bash

# Make program install directory
mkdir /usr/local/flowvizer-edge

# Make program config directory
mkdir /etc/flowvizer-edge/

# Copy the service file to systemd
cp -f ./flowvizer-edge.service /etc/systemd/system/

# Copy content of tar to install directory
cp -rf ./* /usr/local/flowvizer-edge/

# Make changes to content in install directory
cd /usr/local/flowvizer-edge
mv ./livego ./flowvizer-edge
chmod +x ./flowvizer-edge

# Reload system services
systemctl daemon-reload

# Start the Service
systemctl start flowvizer-edge