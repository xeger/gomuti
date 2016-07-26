v1
==

Not yet started.

v0
==

This is an exploratory/preview release with an unstable API. Mocks, stubs and
spies are implemented but their interfaces may still change. The following
issues are outstanding:

Rather than passing/returning raw interface slices to a Mock,
we might want to use a Tuple object that can cope with type checking, bounds
checking, etc ... need to consider how mongoose interacts with gomuti and where
the division of responsibility is. (Currently, mongoose generates code that
panics if an allowed call returns too few/many params, or params of the wrong
type.)
