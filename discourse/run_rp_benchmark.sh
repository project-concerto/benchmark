cd /benchmark/discourse/clients/rp/
mvn clean compile assembly:single
for branch in aht dbt-s dbt-w manual
do
    cd /benchmark/discourse/discourse
    git checkout rp-${branch}
    psql -d postgres -c "drop database discourse WITH (FORCE);"
    psql -d postgres -c "create database discourse;"
    psql -d discourse < /benchmark/discourse/dump/discourse_rp.dump
    rm -rf public/uploads
    cp -r /benchmark/discourse/dump/uploads /benchmark/discourse/discourse/public/
    RAILS_ENV=production bundle exec rails s &
    serverId=$!
    sleep 30s
    # avoding run clients for no contention workload, directly call downsize uploads scripts
    cd /benchmark/discourse/clients/rp/
    java -jar ./target/web-e2eb-1.0-SNAPSHOT-jar-with-dependencies.jar --emulators=64 --warmUp=20 --benchmark=10000 --coolDown=20 --interval=5 &
    clientId=$!
    cd /benchmark/discourse/discourse
    sleep 20s
    ts=$(date +%s%N)
    RAILS_ENV=production WORKER_COUNT=8 WORKER_ID=0 rails runner script/downsize_uploads.rb >> ${branch}-detail.res
    tt=$((($(date +%s%N) - $ts)/1000000))
    echo "Total Time: $tt" >> ${branch}-time.res
    kill -9 $serverId
    kill -9 $clientId
done
