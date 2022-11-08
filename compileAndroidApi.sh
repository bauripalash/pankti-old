#!/bin/env sh
mkdir -p dist
ANDROID_NDK_HOME=$HOME/Android/Sdk/ndk/25.1.8937393/ gomobile bind -v -o dist/androidapi.aar -target=android -javapkg=in.palashbauri.panktijapi -androidapi 19 -tags noide bauri.palash/pankti/androidapi
