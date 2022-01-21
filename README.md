# GoBlog

A homecrafted Blog Written in Golang as a learning exercise. If you want to see how it looks, go to my home page: [macias.info](http://macias.info)

## Blog Structure

The `sample` folder contains an simplified example of the Root contents for a blog. You can override
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

## Environment variables

* `GOBLOG_DOMAIN` (default: value returned by `os.Hostname()`)
    * The domain of your blog (for Atom XML feeds)

* `GOBLOG_PORT` (default: 8080)
    * The port where the HTTP service listens

* `GOBLOG_ROOT` (default: `../sample`)
    * The root folder of the blog contents
    
## Version history

### v2.0 (under development)

* Updated to Go 1.18 to make use of generic stuff
* Addressed HTTPS server
* Redirects old HTTP to HTTPS
* GOBLOG_DOMAIN default value changed from HOSTNAME to `"localhost"`

### v1.0

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
