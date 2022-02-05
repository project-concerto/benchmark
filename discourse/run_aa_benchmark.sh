for branch in aht dbt
do
    cd /benchmark/discourse/discourse
    git checkout aa-${branch}
    for thread in 12 24 36 48 60 72 84 96 108 120 132 144 156
    do
        psql -d postgres -c "drop database discourse WITH (FORCE);"
        psql -d postgres -c "create database discourse;"
        psql -d discourse < /benchmark/discourse/dump/discourse_aa.dump
        cd /benchmark/discourse/discourse
        RAILS_ENV=production bundle exec rails s &
        sleep 30s
        cd /benchmark/discourse/clients/aa/
        THREADS={thread} go run ./discourse-like/bin/main.go -topics=7 -file-prefix=${branch} -append
        kill -9 $!

        psql -d postgres -c "drop database discourse WITH (FORCE);"
        psql -d postgres -c "create database discourse;"
        psql -d discourse < /benchmark/discourse/dump/discourse_aa.dump
        cd /benchmark/discourse/discourse
        RAILS_ENV=production bundle exec rails s &
        sleep 30s
        cd /benchmark/discourse/clients/aa/
        THREADS={thread} go run ./discourse-like/bin/main.go -file-prefix=NC${branch} -append -mode=no-contention
        kill -9 $!
    done
done
