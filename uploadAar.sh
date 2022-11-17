#!/bin/env sh
cd dist
echo "Upload .AAR"
mvn gpg:sign-and-deploy-file -Durl=https://s01.oss.sonatype.org/service/local/staging/deploy/maven2/ -DrepositoryId=ossrh -DpomFile=panktijapi-0.1.0-alpha.1.pom -Dfile=panktijapi-0.1.0-alpha.1.aar
echo "Upload sources.JAR"
mvn gpg:sign-and-deploy-file -Durl=https://s01.oss.sonatype.org/service/local/staging/deploy/maven2/ -DrepositoryId=ossrh -DpomFile=panktijapi-0.1.0-alpha.1.pom -Dfile=panktijapi-0.1.0-alpha.1-sources.jar -Dclassifier=sources
