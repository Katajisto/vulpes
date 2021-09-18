set -e
FAUNADB_SERVER_SECRET="fnAESiHogGACTNAxuOV1_1fbP2BP7m0L5FfMJO83"
export FAUNADB_SERVER_SECRET
echo STARTING DEPLOY
echo STARTING COMPILE
GOOS=linux CGO_ENABLED=0 go build -tags prod -o main
echo COMPILE DONE
echo BUNDLING...
zip lambda.zip views/ main -r
echo BUNDLE DONE
cd cdk
echo STARTING DEPLOY CDK
cdk deploy
echo CDK DONE! DEPLOY DONE!
rm lambda.zip
cd ..