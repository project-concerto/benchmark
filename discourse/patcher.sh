#!/bin/bash
cd /benchmark/discourse/discourse
git config --global user.email "project-concerto@protonmail.com"
git config --global user.name "project-concerto"
git checkout cdf45f4fe6c7e4803762d6945cb51c709a9a0eb3
git checkout -b concerto
git checkout -b aa-aht
git checkout -b aa-dbt
git checkout -b cbc-aht
git checkout -b cbc-dbt
git checkout -b rp-aht
git checkout -b rp-dbt-s
git checkout -b rp-dbt-w
git checkout -b rp-manual
git apply ../patches/rp/RP-MANUAL.patch
git add .
git commit -m "initial rp-manual branch"
git checkout rp-dbt-w
git apply ../patches/rp/RP-DBT-W.patch
git add .
git commit -m "initial rp-dbt-w branch"
git checkout rp-dbt-s
git apply ../patches/rp/RP-DBT-S.patch
git add .
git commit -m "initial rp-dbt-s branch"
git checkout rp-aht
git apply ../patches/rp/RP-AHT.patch
git add .
git commit -m "initial rp-aht branch"
git checkout cbc-dbt
git apply ../patches/cbc/discourse-CBC-DBT.patch
git add .
git commit -m "initial CBC-DBT branch"
git checkout cbc-aht
git apply ../patches/cbc/discourse-CBC-AHT.patch
git add .
git commit -m "initial CBC-AHT branch"

git checkout aa-aht
git apply ../patches/discourse/aa/api-key.path
git apply ../patches/discourse/aa/aht.path
git add .
git commit -m "initial aa-aht branch"

git checkout aa-dbt
git apply ../patches/discourse/aa/api-key.path
git apply ../patches/discourse/aa/dbt.path
git add .
git commit -m "initial aa-dbt branch"

cd ../discourse-solved
git checkout 584060102526cca953f270281751034c14b98d67
git checkout -b concerto
git checkout -b aht
git checkout -b dbt
git apply ../patches/cbc/discourse-solved-CBC-DBT.patch
git add .
git commit -m "initial cbc-dbt branch"
git checkout aht
git apply ../patches/cbc/discourse-solved-CBC-AHT.patch
git add .
git commit -m "initial cbc-aht branch"


cd ..
cp -r discourse-solved ./discourse/plugins
cd ./discourse
bundle install
redis-server &
disown $1
psql -dpostgres -c "create database discourse"
psql -ddiscourse < ../dump/discourse_cbc.dump 
RAILS_ENV=production bundle exec rake assets:precompile
