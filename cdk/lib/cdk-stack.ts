import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda'
import * as apigateway from '@aws-cdk/aws-apigatewayv2'
import * as integrations from '@aws-cdk/aws-apigatewayv2-integrations'
import { memoryUsage } from 'process';

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

    const api = new apigateway.HttpApi(this, "vulpes-api")

    const vulpesIntegration = new integrations.LambdaProxyIntegration({
      handler: handler
    })

    const httpApi = new apigateway.HttpApi(this, 'VulpesHttpApi', {
      defaultIntegration: vulpesIntegration
    });

    httpApi.addRoutes({
      path: '/',
      methods: [ apigateway.HttpMethod.ANY ],
      integration: vulpesIntegration,
    });

  }
}
