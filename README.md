# goblog
A homecrafted Blog Written in Golang as a learning exercise


## How to

Put static assets in `static/` folder. They will be accessible through the `/static/` URL path.

Put blog entries in `entries/` folder as MarkDown documents. They will be accessible through the 
`/entry/` URL path (without extension).

Edit blog template files under the `template/` folder.

## Entries (todo)

* Entries whose file name starts with a timestamp `YYYYMMDDname.md` will be automatically added to the
  index.

* Entries whose file name starts with other pattern will be treated as pages, and will need to link
  them manually in the template or another entry.