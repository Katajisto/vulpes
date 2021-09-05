import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda'
import * as apigateway from '@aws-cdk/aws-apigateway'
import * as cf from '@aws-cdk/aws-cloudfront'
import * as origins from '@aws-cdk/aws-cloudfront-origins'


export class CdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    const handler = new lambda.Function(this, 'VulpesMain', {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset('../lambda.zip'),
      handler: "main",
      memorySize: 500,
    });

    const vulpesApi = new apigateway.LambdaRestApi(this, 'VulpesApi', {
      handler: handler,
      proxy: true
    })

  }
}