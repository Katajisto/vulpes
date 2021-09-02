set -e
echo STARTING DEPLOY
rm lambda.zip
echo DELETED OLD BUNDLE
echo STARTING COMPILE
GOOS=linux CGO_ENABLED=0 go build -tags prod -o main
echo COMPILE DONE
echo BUNDLING...
zip lambda.zip templates/ main -r
echo BUNDLE DONE
cd cdk
echo STARTING DEPLOY CDK
cdk deploy
echo CDK DONE! DEPLOY DONE!
cd ..