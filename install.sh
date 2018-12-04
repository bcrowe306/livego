mkdir /usr/local/flowvizer
mkdir /etc/flowvizer-edge/
cp -rf ./flowvizer-edge.service /etc/systemd/system/
cp -rf ./* /usr/local/flowvizer-edge/
cd /usr/local/flowvizer-edge
mv ./livego ./flowvizer-edge
chmod +x ./flowvizer-edge
systemctl daemon-reload
systemctl start flowvizer-edge