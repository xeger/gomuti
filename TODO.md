v1
==

Create explicit error types.

Provide a better error message when we panic due to:
  - ambiguous mock match
  - no mock match (speculate about best partial match...)
  - wrong parameters/results for `Do()` functions
