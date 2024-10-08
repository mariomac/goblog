## Version history

### v2.x

* Added configurable per-client rate limiter
* Updated goldmark rendering.
* Support github-flavored tables.
* HTTPS is not mandatory but optional. HTTP can be configured
  either for redirecting to HTTPS or for serving the blog.
* Configure Log level

#### Breaking changes

* In index template, you can't iterate entries anymore with
  `{{ range . }}` but you need to use `{{ range .Entries }}`.

### v1.0.0

* Default root path is local folder
* Added `redirect` YAML-only configuration option
* Replaced standard logger by logrus
* Web cache configurable with "GOBLOG_CACHE_SIZE_BYTES" yaml:"cacheSizeBytes"`
* Breaking changes:
  * In entry templates `{{ if .Time }}` must be replaced by `{{ if not .Time.IsZero }}`
  * Index templates will fail if they contain the `{{ if .Time }}` directive
* Invalidating web cache

### v0.2.0

* HTTPS server
* Redirects old HTTP to HTTPS
* GOBLOG_DOMAIN default value changed from HOSTNAME to `"localhost"`
* New YAML/Env dual configuration system

### v0.1.0

* Changed markdown processor
* Added code syntax highlighting
* Modified makefile
* Addressed few code lintings

### v0.0.10

* Added previews to blog entries
    - Updated atom.xml to show previews
    - Updated index.html to show previews

### v0.0.9

* XML Atom Feeds

### v0.0.8

* Basic behaviour. Functional, simple blog

## TO DO

* Reload templates or entries when containing folders change

* Extract first paragraph of entries so you can do a preview

* Main page shows paginated entries

* Add 404 page

* Migrate `log` to `github.com/golang/glog`

* Provide `DockerFile`

* Github triggers to upload blog
