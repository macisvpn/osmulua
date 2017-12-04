#!/bin/bash

# initialisasi var
export DEBIAN_FRONTEND=noninteractive
OS=`uname -m`;
MYIP=$(wget -qO- ipv4.icanhazip.com);
MYIP2="s/xxxxxxxxx/$MYIP/g";
MYPORT="s/85/99/g";

# go to root
cd

#ca-certificates
apt-get install ca-certificates

# setting port ssh
sed -i 's/Port 22/Port 22/g' /etc/ssh/sshd_config
sed -i '/Port 22/a Port 143' /etc/ssh/sshd_config
service ssh restart

# install dropbear
apt-get -y install dropbear
sed -i 's/NO_START=1/NO_START=0/g' /etc/default/dropbear
sed -i 's/DROPBEAR_PORT=22/DROPBEAR_PORT=443/g' /etc/default/dropbear
sed -i 's/DROPBEAR_EXTRA_ARGS=/DROPBEAR_EXTRA_ARGS="-p 443 -p 80"/g' /etc/default/dropbear
echo "/bin/false" >> /etc/shells
echo "/usr/sbin/nologin" >> /etc/shells
service ssh restart
service dropbear restart

# install squid3
cd
apt-get -y install squid3
wget -O /etc/squid3/squid.conf "https://raw.githubusercontent.com/dathai/Wedssh/master/API/squid.conf"
sed -i $MYIP2 /etc/squid3/squid.conf;
service squid3 restart

# install webserver
cd
apt-get -y install nginx php5 php5-fpm php5-cli php5-mysql php5-mcrypt
rm /etc/nginx/sites-enabled/default && rm /etc/nginx/sites-available/default
mv /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup
mv /etc/nginx/conf.d/vps.conf /etc/nginx/conf.d/vps.conf.backup
wget -O /etc/nginx/nginx.conf "https://raw.githubusercontent.com/dathai/Wedssh/master/API/nginx.conf"
wget -O /etc/nginx/conf.d/vps.conf "https://raw.githubusercontent.com/dathai/Wedssh/master/API/vps.conf"
sed -i 's/cgi.fix_pathinfo=1/cgi.fix_pathinfo=0/g' /etc/php5/fpm/php.ini
sed -i 's/listen = \/var\/run\/php5-fpm.sock/listen = 127.0.0.1:9000/g' /etc/php5/fpm/pool.d/www.conf
sed -i $MYPORT /etc/nginx/conf.d/vps.conf;
useradd -m vps && mkdir -p /home/vps/public_html
rm /home/vps/public_html/index.html && echo "<?php phpinfo() ?>" > /home/vps/public_html/info.php
chown -R www-data:www-data /home/vps/public_html && chmod -R g+rw /home/vps/public_html
service php5-fpm restart && service nginx restart

# install openvpn
wget -O /etc/openvpn/openvpn.tar "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/openvpn-debian.tar"
cd /etc/openvpn/
tar xf openvpn.tar
wget -O /etc/openvpn/1194.conf "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/1194.conf"
service openvpn restart
sysctl -w net.ipv4.ip_forward=1
sed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/g' /etc/sysctl.conf
iptables -t nat -I POSTROUTING -s 192.168.100.0/24 -o eth0 -j MASQUERADE
iptables-save > /etc/iptables_yg_baru_dibikin.conf
wget -O /etc/network/if-up.d/iptables "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/iptables"
chmod +x /etc/network/if-up.d/iptables
service openvpn restart

# konfigurasi openvpn
cd /etc/openvpn/
wget -O /etc/openvpn/client.ovpn "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/client-1194.conf"
sed -i $MYIP2 /etc/openvpn/client.ovpn;
cp client.ovpn /home/vps/public_html/



# download script
cd /usr/bin
wget -O menu "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/menu.sh"
wget -O usernew "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/usernew.sh"
wget -O trial "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/trial.sh"
wget -O hapus "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/hapus.sh"
wget -O cek "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/user-login.sh"
wget -O member "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/user-list.sh"
wget -O resvis "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/resvis.sh"
wget -O speedtest "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/speedtest_cli.py"
wget -O info "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/info.sh"
wget -O about "https://raw.githubusercontent.com/dathai/SSH-OpenVPN/master/API/about.sh"

echo "0 0 * * * root /sbin/reboot" > /etc/cron.d/reboot

chmod +x menu
chmod +x usernew
chmod +x trial
chmod +x hapus
chmod +x cek
chmod +x member
chmod +x resvis
chmod +x speedtest
chmod +x info
chmod +x about

#FIGlet In Linux
sudo apt-get install figlet
yum install figlet

#RM file
rm -f myweb.tar
cd
rm -f install.sh

# About
clear
figlet "THAI-VPN"
echo "Script WebSSH Auto Install"
echo "-Web Server"
echo "-Squid Proxy Port 80,8080,3128"
echo "Squid     :  http://$MYIP:8080"
echo "Nginx      :  http://$MYIP:99"
echo "Web    :  http://$MYIP:99"
