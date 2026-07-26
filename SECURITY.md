# Security Policy

## Supported versions

Security fixes apply to the current `main` branch and the latest generated
subscription files. Older commits and previously downloaded subscription
snapshots are not maintained.

## Reporting a vulnerability

Do not include exploit details, credentials, or live node URIs in a public
issue.

Use a private contact method listed on the
[maintainer's GitHub profile](https://github.com/Au1rxx). If none is
available, open a minimal issue asking the maintainer to arrange private
contact. Include no sensitive technical details in that issue.

In the private report, please include:

- the affected file, component, or commit;
- the impact and who could be affected;
- the smallest reproducible example or steps;
- any suggested mitigation; and
- whether the issue has been disclosed elsewhere.

The maintainer will acknowledge the report when practical, investigate it,
and coordinate a fix and disclosure based on its impact.

## Scope

Examples of security issues include committed credentials, unsafe generated
configuration, or a parser flaw that can be triggered by an upstream source.

Individual proxy nodes being offline, slow, blocked, or operated by an
untrusted third party are not security vulnerabilities in this repository.
They are inherent properties of a public aggregator; see the disclaimer in
the README.
