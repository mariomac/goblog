# goblog
A homecrafted Blog Written in Golang as a learning exercise

## Blog Structure

The `sample` folder contains a sample of the Root contents for a blog. You can override the
`GOBLOG_ROOT` environment variable to point the root to another folder.

The Root contents folders is structured as follows:

* `entries` subfolder contains MarkDown entries for each entry of the blog. The name
    * Entries whose file name starts with a timestamp `YYYYMMDDname.md` will be automatically added to the
      index.
    
    * Entries whose file name starts with other pattern will be treated as pages, and will need to link
      them manually in the template or another entry.

* `static` subfolder contains static assets (CSS, images, Javascript files...)

* `template` contains the HTML templates in the Golang templating format.

## How to

Put static assets in `static/` folder. They will be accessible through the `/static/` URL path.

Put blog entries in `entries/` folder as MarkDown documents. They will be accessible through the 
`/entry/` URL path (without extension).

Edit blog template files under the `template/` folder.

## Entries (todo)


## Environment variables

* `GOBLOG_PORT` (default: 8080)
    * The port where the HTTP service listens

* `GOBLOG_ROOT` (default: `../sample`)
    * The root folder of the blog contents

## TO DO

* Reload templates or entries when containing folders change

* Main page shows paginated entries

* Add 404 page