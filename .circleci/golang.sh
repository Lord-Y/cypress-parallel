#!/bin/sh

VERSION="$1"

install()
{
  mkdir -p ~/go/bin
  sudo rm -rf /usr/local/go ~/.cache/go-build/* ~/go/pkg/mod/* ~/go/src/* /tmp/go${VERSION}.linux-amd64.tar.gz
  wget https://golang.org/dl/go${VERSION}.linux-amd64.tar.gz -P /tmp
  sudo tar -C /usr/local -xzf /tmp/go${VERSION}.linux-amd64.tar.gz
  rm -rf /tmp/go${VERSION}.linux-amd64.tar.gz
  grep /usr/local/go/bin ~/.bashrc &> /dev/null
  if [ $? eq 1 ]
  then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
    echo 'export GO111MODULE=on' >> ~/.bashrc
    source ~/.bashrc
  fi
}

check()
{
  which go
}

check_installed_version()
{
  go version |grep $VERSION &>/dev/null
}

if [ -n "$VERSION" ]
then
  which go &> /dev/null
  if [ $? eq 1 ]
  then
    install
  fi
  check_installed_version
  if [ $? eq 1 ]
  then
    install
  fi
  go version
else
  echo "golang version is missing"
  exit 1
fi
