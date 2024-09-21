# GoBlog

A blogging platform written by a coder, for coders.

If you want to see how it looks, go to my home page: [macias.info](http://macias.info)

## Building

```
make compile
```
or
```
go build -o goblog ./src
```

## Running sample

```
make sample
```

or

```
./goblog -cfg /path/to/goblog/config.yml
```

If you want to run a local copy of [my own blog at macias.info](https://macias.info), you can run:

```
GOBLOG_ROOT=./etc/macias.info make sample
```

## Configuring

Default configuration can be overridden by a YAML file config and/or Environment variables.

The YAML file config path must be passed by the `-cfg` command-line argument or the `GOBLOG_CONFIG`
environment variable.

Environment variables take precedence over YAML configuration.

* env: `GOBLOG_ROOT`, yaml: `rootPath`
  * The root folder of the blog contens (see [Blog Structure](#blog-structure))
  * Default: `./sample`
* env: `GOBLOG_HTTPS_PORT`, yaml: `httpsPort`
  * Port to serve the secure HTTPS content. If set to <0, HTTPS will be disabled.
  * Default: `8443`
* env: `GOBLOG_HTTP_PORT`, yaml: `httpPort`
  * Port to listen for any HTTP request. If set to <0, HTTP will be disabled.
  * Default: `8080`
* env: `GOBLOG_HTTPS_REDIRECT`, yaml: `httpsRedirect`
  * If set to `true`, any HTTP request will be redirected to HTTPS. If set to `false`, content is served
    directed in an insecure port.
  * Default: `true`
* env: `GOBLOG_DOMAIN`, yaml: `domain`
  * Domain/hostname/IP where the blog is going to be visible from
  * Default: `localhost`
* env: `GOBLOG_TLS_CERT`, yaml: `tlsCertPath`
* env: `GOBLOG_TLS_KEY`, yaml: `tlsKeyPath`
  * Paths of the TLS certificate and key for HTTPS serving
  * Default: empty
* env: `GOBLOG_CACHE_SIZE_BYTES`, yaml:`cacheSizeBytes`
  * Size, in bytes, of the HTTP cache to minimize disk loads and template renderings
  * Default: 32MB

TODO: explain `redirect` map in yaml.

## Blog Structure

The `sample` folder contains a simplified example of the Root contents for a blog. You can override
the `GOBLOG_ROOT` environment variable to point the root to another folder.

The Root contents folders is structured as follows:

* `entries` subfolder contains MarkDown entries for each entry of the blog.
    * Entries whose file name starts with a timestamp `YYYYMMDDHHMMname.md` will be automatically added to the
      index.
    
    * Entries whose file name starts with other pattern will be treated as pages, and will need to link
      them manually in the template or another entry.

* `static` subfolder contains static assets (CSS, images, Javascript files...)

* `template` contains the HTML templates in the Golang templating format.
    * The template MUST contain at least two files: `index.html` for the main index page, and
      `entry.html` for the blog entry page.

## How to create your blog

Put static assets in `static/` folder. They will be accessible through the `/static/` URL path.

Put blog entries in `entries/` folder as MarkDown documents. They will be accessible through the 
`/entry/` URL path (without extension).

Edit blog template files under the `template/` folder.

## How to add an entry to your blog

Just add a file in the `entries/` folder in a timestamped format. E.g. `201711281330_hello.md`
will create an entry created at November 28th, 2017 at 13:30.

The markdown file MUST contain a First-level header (e.g. `# Post title`), that will be used
as title of the entry in the entry heading and links.

At this early stage of the blog, you *MUST* restart the blog process before changes are visible. 
