# Pancake Buildpack

This tiny buildpack does a simple thing for your application - it flattens the crazy `$VCAP_SERVICES` service instance credentials into many specific environment variables.

In the following example, in addition to `$VCAP_SERVICES` the application will also have variables starting with `MYSQL_` for each credential (such as `MYSQL_HOSTNAME`, `MYSQL_USERNAME`, `MYSQL_URI`, for credentials `hostname`, `username`, and `uri`):
