FROM heroku/cedar:14
MAINTAINER Fabio Rehm "fgrehm@gmail.com"

ENV HOME="/home/developer" \
    PATH="/home/developer/bin:/home/developer/android-sdk-linux/tools:/home/developer/android-sdk-linux/platform-tools:$PATH" \
    ANDROID_HOME="/home/developer/android-sdk-linux" \
    JAVA_HOME="/usr/lib/jvm/java-7-openjdk-amd64"

RUN set -x \
    && echo "developer:x:1000:1000:developer,,,:/home/developer:/bin/bash" >> /etc/passwd \
    && echo "developer:x:1000:" >> /etc/group \
    && mkdir -p $HOME \
    && curl -L http://dl.google.com/android/android-sdk_r24.4-linux.tgz | tar -zx -C $HOME \
    && echo y | android update sdk --all --no-ui --force --filter android-22 \
    && echo y | android update sdk --all --no-ui --force --filter tools \
    && echo y | android update sdk --all --no-ui --force --filter platform-tools \
    && echo y | android update sdk --all --no-ui --force --filter build-tools-22.0.1 \
    && echo y | android update sdk --all --no-ui --force --filter extra-android-m2repository \
    && echo y | android update sdk --all --no-ui --force --filter extra-google-m2repository \
    && echo y | android update sdk --all --no-ui --force --filter extra-google-google_play_services \
    && echo y | android update sdk --all --no-ui --force --filter extra-google-play_licensing \
    && echo y | android update sdk --all --no-ui --force --filter extra-google-gcm \
    && chown 1000:1000 -R $HOME

RUN dpkg --add-architecture i386 \
    && curl -sL https://deb.nodesource.com/setup_0.12 | bash - \
    && apt-get install -y --force-yes \
                       expect \
                       ant \
                       wget \
                       libc6-i386 \
                       lib32stdc++6 \
                       lib32gcc1 \
                       lib32ncurses5 \
                       lib32z1 \
                       nodejs \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/*

ENV GOROOT="/usr/lib/go" \
    GOPATH="/go" \
    GOBIN="/go/bin" \
    PATH="./node_modules/.bin::/go/bin:/usr/lib/go/bin:$PATH"

RUN set -x \
    && curl -L https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz \
       | tar xz -C $(dirname $GOROOT) \
    && go get github.com/parkghost/watchf/... \
    && go get github.com/constabulary/gb/... \
    && go get github.com/jteeuwen/go-bindata/... \
    && go get github.com/elazarl/go-bindata-assetfs/... \
    && rm -rf /tmp/* \
    && mkdir -p $GOPATH \
    && chown 1000:1000 -R $GOPATH \
    && mkdir -p $GOROOT \
    && chown 1000:1000 -R $GOROOT

USER developer

WORKDIR /code
