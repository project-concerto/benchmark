#!/bin/bash
cd /
wget https://downloads.mysql.com/archives/get/p/23/file/mysql-5.7.35-linux-glibc2.12-x86_64.tar.gz
groupadd mysql
useradd -r -g mysql -s /bin/false mysql
cd /usr/local
tar zxvf /mysql-5.7.35-linux-glibc2.12-x86_64.tar.gz
ln -s mysql-5.7.35-linux-glibc2.12-x86_64 mysql
cd mysql
mkdir mysql-files
chown mysql:mysql mysql-files
chmod 750 mysql-files
bin/mysqld --initialize-insecure --user=mysql
bin/mysql_ssl_rsa_setup
bin/mysqld_safe --user=mysql &
sleep 60s
bin/mysql -uroot mysql -e 'update user set authentication_string=PASSWORD("123456") where user="root"; flush privileges;'
bin/mysql -uroot -p123456 -e "create database broadleaf"
bin/mysql -uroot -p123456 broadleaf < /benchmark/broadleaf/broadleaf.dump
