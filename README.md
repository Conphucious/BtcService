# Btc Service
Pulls from a BTC restful API using GOlang (marshalls the JSON) and uploads it to a NoSQL database (DynamoDB). The GOlang program is contained via a Docker container which is then pushed into an AWS EC2 instance with limited IAM policies managed by AWS ECR. A web end-point was accessible to view our EC2 instance's data retrieval and AWS CodePipeline was used to manage CI/CD deployment.
