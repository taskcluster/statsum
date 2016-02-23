Statistical Summarization Service
=================================

Simple service that accepts metrics such as counter and values. Counters are
aggregated the usual way, but for the value metrics the service will estimate
percentiles using t-digests.

For every 5 minute and 1 hour interval the service will forward estimated
percentiles (and aggregated counters) to a time series service like signalfx,
datadog or similar service.


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
