# Contributing

Thanks for your interest in contributing to the Pinot Terraform Provider! Please take a moment to review this document before submitting a pull request.

## Pull requests

Please ask first before starting work on any significant new features.

Please ensure that the tests are passing when submitting a pull request. If you're adding new features, please include tests. 

See the terraform provider tutorial testing section for more information: [Provider Acceptance Testing](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-acceptance-testing#providers-plugin-framework-acceptance-testing)

## Where do I go from here?

- For any questions, support, or ideas, etc. please create a GitHub discussion.
- If you've noticed a bug, please submit an issue.

#### 1. Fork and create a branch

If this is something you think you can fix, then [fork terraform-provider-pinot] and create a branch with a descriptive name.

#### 2. Implement your fix or feature
At this point, you're ready to make your changes. Feel free to ask for help.

Be sure to also run `go generate` in the project root to generate documentation for your features.

#### 3. Create a Pull Request
At this point, if your changes look good and tests are passing, you are ready to create a pull request.

## Merging a PR (maintainers only)
A PR can only be merged into main by a maintainer if: 
1. CI is passing, and
2. The PR is up to date with the default branch. 

Any maintainer is allowed to merge a PR if these conditions ae met.

## 