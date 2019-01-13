# cloudping

cloudping identifies the regions geographically closest and returns them in order of lowest to highest "response time".

> Inspired by [CloudPing.info](https://www.cloudping.info/).

## Usage

```bash
cloudping --provider=aws --regions=us-east-1,us-east-2
```

## Why? 

The idea came from the need to download images from the geographically closest docker registry.
We operate our Kubernetes cluster in Sydney/Australia but use CircleCI operating in the US.
Because AWS ECR doesn't provide a common endpoint with geographically distributed backend we push our images to both locations.
Within our Makefile we can use `cloudping` to identify if images should be pulled from the US or Sydney. 

## Similar projects

* [CloudPing.info](https://www.cloudping.info/)
* [AWS Inter-Region Latency Monitoring](https://www.cloudping.co/)
* [GCP ping](http://www.gcping.com/)