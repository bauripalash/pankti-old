#!/bin/env sh
mkdir -p dist
export PANTKI_VERSION="0.1.1-alpha.1"
ANDROID_NDK_HOME=$HOME/Android/Sdk/ndk/25.1.8937393/ gomobile bind -v -o dist/panktijapi-$PANTKI_VERSION.aar -target=android -javapkg=in.palashbauri.panktijapi -androidapi 19 -tags noide go.cs.palashbauri.in/pankti/androidapi
cp ./androidapi/androidapi.pom dist/panktijapi-$PANTKI_VERSION.pom
