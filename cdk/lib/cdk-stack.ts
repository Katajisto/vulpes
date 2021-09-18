import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda'
import * as apigateway from '@aws-cdk/aws-apigateway'
import * as acm from '@aws-cdk/aws-certificatemanager'
import * as cf from '@aws-cdk/aws-cloudfront'
import * as origins from '@aws-cdk/aws-cloudfront-origins'
import { domain } from 'process';


export class CdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    const handler = new lambda.Function(this, 'VulpesMain', {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset('../lambda.zip'),
      handler: "main",
      memorySize: 500,
      environment: {
        FAUNADB_SERVER_SECRET:  process.env.FAUNADB_SERVER_SECRET || '',
      }
    });

    const vulpesApi = new apigateway.LambdaRestApi(this, 'VulpesApi', {
      handler: handler,
      proxy: true,
      endpointConfiguration: {
        types: [apigateway.EndpointType.REGIONAL]
      },
      domainName: {
        domainName: `vulpes.ktj.st`,
        certificate: acm.Certificate.fromCertificateArn(
          this,
          "VulpesCert",
          "arn:aws:acm:eu-west-1:693703930136:certificate/f3e0831e-488e-4f99-afe5-591a9474ec2b"
        ),
        endpointType: apigateway.EndpointType.REGIONAL,
      },
    })
  }
}