#!/bin/env sh
cd dist
export PANTKI_VERSION="0.2.0-alpha.1"
echo "v${PANTKI_VERSION}"
echo "Upload .AAR"
mvn gpg:sign-and-deploy-file -Durl=https://s01.oss.sonatype.org/service/local/staging/deploy/maven2/ -DrepositoryId=ossrh -DpomFile=panktijapi-$PANTKI_VERSION.pom -Dfile=panktijapi-$PANTKI_VERSION.aar
echo "Upload sources.JAR"
mvn gpg:sign-and-deploy-file -Durl=https://s01.oss.sonatype.org/service/local/staging/deploy/maven2/ -DrepositoryId=ossrh -DpomFile=panktijapi-$PANTKI_VERSION.pom -Dfile=panktijapi-$PANTKI_VERSION-sources.jar -Dclassifier=sources
