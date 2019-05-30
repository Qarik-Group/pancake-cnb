# Pancake Buildpack

This tiny buildpack does a simple thing for your application - it flattens the crazy `$VCAP_SERVICES` service instance credentials into many specific environment variables.

In the following example, in addition to `$VCAP_SERVICES` the application will also have variables starting with `MYSQL_` for each credential (such as `MYSQL_HOSTNAME`, `MYSQL_USERNAME`, `MYSQL_URI`, for credentials `hostname`, `username`, and `uri`):

## Usage

```plain
rm -rf pancake-cnb_*
./scripts/package.sh
pack build phpap --path integration/fixtures/phpapp \
  --buildpack pancake-cnb_* \
  --buildpack https://github.com/cloudfoundry/php-cnb/releases/download/v0.0.3/php-cnb-0.0.3.tgz \
  --buildpack https://github.com/cloudfoundry/httpd-cnb/releases/download/v0.0.2/httpd-cnb-0.0.2.tgz \
  --buildpack https://github.com/cloudfoundry/php-web-cnb/releases/download/v0.0.4/php-web-cnb-0.0.4.tgz \
  --env VCAP_APPLICATION={} \
  --env "VCAP_SERVICES=$(cat integration/fixtures/vcap_services/two-services.json)"
```