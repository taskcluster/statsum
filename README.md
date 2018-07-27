Statistical Summarization Service
=================================

Simple service that accepts metrics such as counter and values. Counters are
aggregated the usual way, but for the value metrics the service will estimate
percentiles using t-digests.

For every 5 minute and 1 hour interval the service will forward estimated
percentiles (and aggregated counters) to a time series service like signalfx,
datadog or similar service.

### Submitting Data

Data is submitted as follows using a project specific `jwt-token` and a request specific `uuid`.
The `uuid` should be reused when a request is retried due to connection error, this avoids duplication of values.

```
POST /v1/project/<projectName>
content-type:         application/json
accept:               application/json
content-length:       <length>
authorization:        bearer <jwt-token>
x-statsum-request-id: <uuid>
{
  "counters": [ // values for a counter are summed up
    {"k": "name.of.counter", v: 5},
  ],
  "measures: [ // values for a meaure is fed to a t-digest for percentile estimation
    {"k": "name.of.measure", v: [3,4,5,65,5,5,7,8,9]},
  ]
}
```

As a rule of thumb, it is reasonable to accumulate metrics for 30 to 90 seconds before flushing. As statsum will only reports every 5min.

### Motivation
Services like signalfx, datadog, stathat, etc. cannot compute or estimate
percentiles. Any function in their dashboard or analysis tools pretending to do
so is false.

These services aggregates data-points as averages (typically averages per second).
Regardless of the resolution, "percentiles over averages" is not a thing.
It is not **statically sound** (period!).

If this service is used to aggregate data-points before they are forwarded to
signalfx, datadog or similar service, you should end up with valid estimates of
the 25'th, 50'th, 75'th, 95'th and 99'th percentile.

But beware when displaying these metrics in signalfx, datadog or similar
services. These services may still aggregate the numbers when rendering graphs.
Any such aggregation is naturally incorrect.

### Building Statsum Server
The server is located in `/cmd/statsum` to build it you must fetch govendor
locked dependencies with `govendor sync`, then `go build ./cmd/statsum` will
build the binary.

Alternatively, the `make build` will do the same and build the docker image too.

### Configuring the Server

When deploygin the server as built above, it requires the following
environment variables:

 * `JWT_SECRET_KEY`, symmetric secret that JWTs are signed with,
 * `TLS_CERTIFICATE`, (optional) TLS certificate chain,
 * `TLS_KEY`, (optional) TLS private key,
 * `PORT`, port to listen on,
 * `SIGNALFX_TOKEN`, signalfx token for sending metrics,
 * `SENTRY_DSN`, sentry DSN for reporting errors,
 * `PUBLIC_URL`, public URL to this server,

## Service Owner

Service Owner: dustin@mozilla.com
