##### Pull Request Guidelines

These are recommendations for pull requests.
They are strictly guidelines to help manage expectations.

##### PR preparation
Run `make pr-prep` from the root of the repository to run formatting, linting and tests.

##### Should this be an issue instead
- [ ] is it a convenience method? (no new functionality, streamlines some use case)
- [ ] exposes a previously private type, const, method, etc.
- [ ] is it application specific (caching, retry logic, rate limiting, etc)
- [ ] is it performance related.

##### API changes

Since API changes have to be maintained they undergo a more detailed review and are more likely to require changes.

- no tests, if you're adding to the API include at least a single test of the happy case.
- If you can accomplish your goal without changing the API, then do so.
- dependency changes. updates are okay. adding/removing need justification.

###### Examples of API changes that do not meet guidelines:
- in library cache for users. caches are use case specific.
- Convenience methods for Sending Messages, update, post, ephemeral, etc. consider opening an issue instead.
