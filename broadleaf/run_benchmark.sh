#!/bin/bash
cd /benchmark/broadleaf/client
mvn clean compile assembly:single
for branch in aht dbt
do
    cd /benchmark/broadleaf/broadleafcommerce
    git checkout ${branch}
    mvn clean install -DskipTests
    for thread in 24 32 40 48
    do
        /usr/local/mysql/bin/mysql -uroot -p123456 -e "drop database broadleaf; create database broadleaf;"
        /usr/local/mysql/bin/mysql -uroot -p123456 broadleaf < /benchmark/broadleaf/broadleaf.dump
        cd /benchmark/broadleaf/DemoSite/site
        mvn spring-boot:run > /dev/null &
        sleep 120s
        cd /benchmark/broadleaf/client
        java -jar ./target/web-e2eb-1.0-SNAPSHOT-jar-with-dependencies.jar --emulators=${thread} --warmUp=20 --benchmark=60 --coolDown=20 --interval=5 | tee -a ${branch}_benchmark.res
        kill -9 $!
        # use following commands to ensure application is always stopped correctly
        kill -9 $(ps aux | grep 'DemoSite/site'| grep -v grep | awk '{print $2}')
        kill -9 $(ps aux | grep 'jdk' | grep -v grep | awk '{print $2}')
        sleep 10s
    done
done
