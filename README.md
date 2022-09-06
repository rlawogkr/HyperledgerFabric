# HyperledgerFabric


How to start Caliper
1. npm install --only=prod @hyperledger/caliper-cli@0.5.0
2. npx caliper bind --caliper-bind-sut fabric:2.2
3. npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/workload/config.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
