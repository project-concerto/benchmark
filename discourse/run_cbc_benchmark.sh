cd clients/cbc/
mvn clean compile assembly:single
for branch in aht dbt
do
    cd /benchmark/discourse/discourse
    git checkout cbc-${branch}
    cd /benchmark/discourse/discourse/plugins/discourse-solved
    git checkout ${branch}
    for thread in 128 144 160
    do
        psql -d postgres -c "drop database discourse WITH (FORCE);"
        psql -d postgres -c "create database discourse;"
        psql -d discourse < /benchmark/discourse/dump/discourse_cbc.dump
        cd /benchmark/discourse/discourse
        RAILS_ENV=production bundle exec rails s &
        sleep 30s
        cd /benchmark/discourse/clients/cbc/
        java -jar ./target/web-e2eb-1.0-SNAPSHOT-jar-with-dependencies.jar --emulators=${thread} --warmUp=20 --benchmark=40 --coolDown=20 --interval=5 | tee -a ${branch}-${thread}.res
        kill -9 $!
    done
done
