#!/bin/bash
cd /benchmark/broadleaf/broadleafcommerce
git config --global user.email "project-concerto@protonmail.com"
git config --global user.name "project-concerto"
git checkout d4b48995dfeee46a4b227ce39783cff940254834
git checkout -b aht
git checkout -b dbt
git apply ../patches/broadleafcommerce-DBT.patch
git add .
git commit -m "initial database transaction branch"
git checkout aht
git apply ../patches/broadleafcommerce-AHT.patch
git add .
git commit -m "initial ad hoc transaction branch"
mvn clean install -DskipTests

cd ../DemoSite
git checkout 8e0b44d77d4b8de445ca7b6e7b6cb477fa392b27
git checkout -b benchmark
git apply ../patches/DemoSite.patch
git add .
git commit -m "initial benchmark branch"
mvn clean install

cd ../broadleaf-boot-starter-database
git checkout 07f31fb7ee8e2deb4ae9ad01d2ffeb91bb5680e3
git checkout -b benchmark
git apply ../patches/broadleaf-boot-starter-database.patch
git add .
git commit -m "initial benchmark branch"
mvn clean install
