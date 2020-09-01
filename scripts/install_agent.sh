#/bin/bash
if [ -z "$GOPATH" ]
then
      echo "\$GOPATH is empty"
      exit 1
else
      echo "\$GOPATH is NOT $GOPATH"
fi

GITLAB_TOKEN=zVcdyFsFaUVbNgoaG6Bw
REF=feature

# build
SRC_DIR=$GOPATH/src/code.htres.cn/casicloud
mkdir -p $SRC_DIR
cd $SRC_DIR
git clone https://oauth2:${GITLAB_TOKEN}@code.htres.cn/casicloud/alb.git
cd alb && git checkout $REF
cd $SRC_DIR/alb/cmd/lbagent
go build && go install

# install
BIN_DIR=$GOPATH/bin/alb
mkdir -p $BIN_DIR/logs
cp $GOPATH/bin/lbagent $BIN_DIR

#clean
cd $GOPATH && rm -rf $SRC_DIR

echo "lbagent now avaliable in $GOPATH/bin"