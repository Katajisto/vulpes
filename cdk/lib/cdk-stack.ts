import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda";
import * as apigateway from "@aws-cdk/aws-apigateway";
import * as acm from "@aws-cdk/aws-certificatemanager";
import * as cf from "@aws-cdk/aws-cloudfront";
import * as origins from "@aws-cdk/aws-cloudfront-origins";
import { domain } from "process";
import s3deploy = require("@aws-cdk/aws-s3-deployment");

import * as s3 from "@aws-cdk/aws-s3";

export class CdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    const handler = new lambda.Function(this, "VulpesMain", {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset("../lambda.zip"),
      handler: "main",
      memorySize: 500,
      environment: {
        FAUNADB_SERVER_SECRET: process.env.FAUNADB_SERVER_SECRET || "",
      },
    });

    const siteBucket = new s3.Bucket(this, "SiteBucket", {
      bucketName: "vulpes-static",
      websiteIndexDocument: "index.html",
      websiteErrorDocument: "error.html",
      publicReadAccess: true,

      // The default removal policy is RETAIN, which means that cdk destroy will not attempt to delete
      // the new bucket, and it will remain in your account until manually deleted. By setting the policy to
      // DESTROY, cdk destroy will attempt to delete the bucket, but will error if the bucket is not empty.
      removalPolicy: cdk.RemovalPolicy.DESTROY, // NOT recommended for production code
    });

    const vulpesApi = new apigateway.LambdaRestApi(this, "VulpesApi", {
      handler: handler,
      proxy: true,
      endpointConfiguration: {
        types: [apigateway.EndpointType.REGIONAL],
      },
    });

    const deployment = new s3deploy.BucketDeployment(
      this,
      "deploySiteStaticAssets",
      {
        sources: [s3deploy.Source.asset("../s3")],
        destinationBucket: siteBucket,
      }
    );

    const apiProxyPolicy = new cf.OriginRequestPolicy(
      this,
      "PasstroughEverythingPolicy",
      {
        cookieBehavior: cf.OriginRequestCookieBehavior.all(),
        queryStringBehavior: cf.OriginRequestQueryStringBehavior.all(),
        headerBehavior: cf.OriginRequestHeaderBehavior.all(),
      }
    );

    const dist = new cf.CloudFrontWebDistribution(this, "Distribution", {
      aliasConfiguration: {
        names: ["vulpes.ktj.st"],
        acmCertRef:
          "arn:aws:acm:us-east-1:693703930136:certificate/6545fe74-84b2-4a94-a1ae-8cd467426507",
      },
      defaultRootObject: "/",
      originConfigs: [
        {
          s3OriginSource: { s3BucketSource: siteBucket },
          behaviors: [
            {
              maxTtl: cdk.Duration.minutes(10),
              pathPattern: "/static/*",
              allowedMethods: cf.CloudFrontAllowedMethods.GET_HEAD,
            },
          ],
        },
        {
          customOriginSource: {
            originProtocolPolicy: cf.OriginProtocolPolicy.HTTPS_ONLY,
            domainName: `${vulpesApi.restApiId}.execute-api.${this.region}.${this.urlSuffix}`,
          },
          originPath: `/${vulpesApi.deploymentStage.stageName}`,

          behaviors: [
            {
              maxTtl: cdk.Duration.millis(0),
              pathPattern: "/*",
              isDefaultBehavior: true,
              allowedMethods: cf.CloudFrontAllowedMethods.ALL,
              forwardedValues: {
                queryString: true,
                cookies: { forward: "all" },
              },
            },
          ],
        },
      ],
    });

    new cdk.CfnOutput(this, "ApiUrl", { value: vulpesApi.url });
  }
}
