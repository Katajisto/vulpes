import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda'
import * as apigateway from '@aws-cdk/aws-apigateway'

export class CdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    const handler = new lambda.Function(this, 'VulpesMain', {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset('lambda'),
    });

  }
}
