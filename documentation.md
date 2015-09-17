# Code

## Writing the code

New code, regardless of type (bugfix or new feature), must pass the same process:

- a new branch should be created from latest master
- changes should be reviewed before master commit (tests should be present, where needed)
- once a change is accepted it will be merged to master (or delayed and updated until such time the merge is possible)

## Testing the code

On the local machine you have to do the following:

```bash
cd $GOPATH/src/github.com/tapglue/backend/v02/server
CI=true go test -tags postgres -check.v
```

In the ```$GOPATH/src/github.com/tapglue/backend/v02/server``` there must be a file called ```config.json``` with the
correct configuration for the test runner.

WARNING!
The tests will wipe out the schemas and database information so be careful where the database is configured to point to.

Note:
The ```CI=true``` will switch the scrypt implementation to use weaker values so that the tests can run in an acceptable
ammount of time.

### Testing with Kinesis on localhost

First you must get the kinesalite emulator:
```bash
docker run -d -t -p 127.0.0.1:4567:4567 dlsniper/kinesalite:1.7.1
```

After this, in the ```config.json``` file you should have something like the following configuration:
```json
  "kinesis": {
    "auth_key": "demo",
    "secret_key": "demo",
    "region": "eu-central-1",
    "endpoint": "http://127.0.0.1:4567"
  },
```

Note the ```"endpoint"``` configuration will point the Kinesis to your local emulator. If you want to test with the real
kinesis from AWS, you can remove the ```"endpoint"``` key and provide the correct details.

```bash
CI=true go test -tags kinesis -check.v
```

## Breaking changes in the API

Breaking changes in the API should be done by bumping the API version in all cases but one.

The only exception to this rule will be security updates which can break the existing API versions after all the clients
have been notified accodingly.

## From code to artifacts

Every commit will be run thru the continuous integration pipeline.

The service that runs the CI tests is [CircleCI](https://circleci.com).

The steps are defined in [circle.yml](circle.yml) file in the root of the repository.

Once all the tests are run an artifact with the build ID is deployed to S3.

Each project will contain a releases.json file which will contain the current deployed version and the latest built
version for that project.

## From artifacts to production

In order for an artifact to be deployed the following steps must be taken:

- a list of currently active instance should be maintained
- the file releases.json for the corresponding application must be changed and the ```current_release``` must be updated
 to the desired build number
- the auto-scaling group desired number of instances should be increased by 1 so that the new instance is used in
 production (where the change allows it)
- once everything works ok with that instace, the instances with the old version should be killed one by one until
 everything was updated to the new release

# Application architecture

## Hardware

### Used AWS services

Tapglue runs in AWS VPC using the following services:

- EC2
- VPC
- ELB
- RDS (Postgres)
- Kinesis
- EC (Redis)

### Infrastructure layout

There are two VPCs, one production and one staging. Both VPCs are configured with the same layout.

There's also a third VPC present which is used for non-mission critical loads such as the corporate website.

The VPCs are configured like this:

- all availability zones in a region are used (e.g. if region is eu-central-1 then a & b are used, if region is eu-west-1 then a, b & c are used)
- there are 4 private subnets / availability zone devided by role:
 - one subnet holds the public facing elements like ELB, Internet Gateway, bastion host, etc.
 - one subnet holds the instances that are used in the frontend
 - one subnet holds the instances that are used in the backend loads
 - one subnet holds the instances that are used in the RDS workloads
- each subnet will be allowed to communicate only with the instances that it needs, via security groups (frontends will not be allowed to communicate with the database directly)
- each subnet type has it's own auto-scaling rules

Each auto-scaling configuration uses the ```tapglue-generic-installer-01``` AMI. This image comes preconfigured with
the ```tapglue``` which is used to run all the Tapglue software. For this user it can't be established a direct ssh connection.

For connecting to the instance via SSH, each team member has an individual user backed in the AMI with their 4096 bit RSA key.

## Software

```tapglue-generic-installer-01``` image comes with the following preinstalled / preconfigured sofware:
- rsyslog with Loggly integration for forwarding
- datadog agent
- mosh, git, iptraf, htop, lynx, zhs (+oh-my-zsh), fish, aws cli

On startup, the instance will run ```/home/tapglue/installer.sh``` which will read the instance tags, download the according software installer and run it.
