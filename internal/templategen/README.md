# Template-Generator
Template Generator transforms Portal API metadata into valid Gateway Bundles that can either be bootstrapped or applied directly to running Gateways. This functionality is also present in the Portal Init Container.

1. go generate

# Issues
1. QuickTemplate is encoding apostrophe as &#39 , velocity/Portal encodes as &apos

