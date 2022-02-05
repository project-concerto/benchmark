#!/bin/bash
trap 'ret=$?; test $ret -ne 0 && printf "failed\n\n" >&2; exit $ret' EXIT

set -e
log_info() {
  printf "\n\e[0;35m $1\e[0m\n\n"
}

if [ ! -f "$HOME/.bashrc" ]; then
  touch $HOME/.bashrc
fi

log_info "Installing redis ..."
apt-get install -y redis-server

log_info "Installing postgresql ..."
apt-get install -y gnupg
sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt focal-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
apt-get update
apt-get install -y postgresql-13 postgresql-contrib 
sed -i "s/md5/trust/" /etc/postgresql/13/main/pg_hba.conf
sed -i "s/max_connections = 100/max_connections = 1024/" /etc/postgresql/13/main/postgresql.conf
service postgresql start
runuser -l postgres -c "psql -c 'CREATE USER root WITH SUPERUSER'"

log_info "Updating Packages ..."
  apt-get update

log_info "Installing build essentials ..."
  apt-get -y install build-essential

log_info "Installing libraries for common gem dependencies ..."
  apt-get -y install libxslt1-dev libcurl4-openssl-dev libksba8 libksba-dev libreadline-dev libssl-dev zlib1g-dev libsnappy-dev

log_info "Installing sqlite3 ..."
 apt-get -y install libsqlite3-dev sqlite3

log_info "Installing ImageMagick ..."
  apt-get -y install libtool
  bash <(wget -qO- https://raw.githubusercontent.com/discourse/discourse_docker/master/image/base/install-imagemagick)

log_info "Installing image utilities ..."
  apt-get -y install advancecomp libjpeg-progs pngcrush

if [[ ! -d "$HOME/.rbenv" ]]; then
  log_info "Installing rbenv ..."
    git clone https://github.com/rbenv/rbenv.git ~/.rbenv

    if ! grep -qs "rbenv init" ~/.bashrc; then
      printf 'export PATH="$HOME/.rbenv/bin:$PATH"\n' >> ~/.bashrc
      printf 'eval "$(rbenv init - --no-rehash)"\n' >> ~/.bashrc
    fi

    export PATH="$HOME/.rbenv/bin:$PATH"
    eval "$(rbenv init -)"
fi

if [[ ! -d "$HOME/.rbenv/plugins/ruby-build" ]]; then
  log_info "Installing ruby-build, to install Rubies ..."
    git clone https://github.com/rbenv/ruby-build.git ~/.rbenv/plugins/ruby-build
fi

ruby_version="2.6.5"

log_info "Installing Ruby $ruby_version ..."
  rbenv install "$ruby_version"

log_info "Setting $ruby_version as global default Ruby ..."
  rbenv global $ruby_version
  rbenv rehash

log_info "Updating to latest Rubygems version ..."
  gem update --system

log_info "Installing Rails ..."
  gem install rails

log_info "Installing Bundler ..."
  gem install bundler

log_info "Installing MailHog ..."
  wget -qO /usr/bin/mailhog https://github.com/mailhog/MailHog/releases/download/v1.0.1/MailHog_linux_amd64
  chmod +x /usr/bin/mailhog

log_info "Installing Node.js 14 ..."
  curl -sL https://deb.nodesource.com/setup_14.x | bash -
  apt-get -y install nodejs
  npm install -g yarn

log_info "Installing go"
  wget https://go.dev/dl/go1.17.6.linux-amd64.tar.gz
  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz
  rm go1.17.6.linux-amd64.tar.gz
  export PATH=$PATH:/usr/local/go/bin
