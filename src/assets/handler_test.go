package assets

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testBlog = "../../testresources/testblog"

func testServer(t *testing.T) *httptest.Server {
	ch, err := NewCachedHandler(testBlog, false, "www.superblog.com")
	require.NoError(t, err)

	return httptest.NewServer(&ch)
}

func doGet(t *testing.T, srv *httptest.Server, path string) WebAsset {
	resp, err := srv.Client().Get(srv.URL+"/"+path)
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return WebAsset{
		MimeType: resp.Header.Get("Content-Type"),
		Body: body,
	}
}

func TestAtom(t *testing.T) {
	s := testServer(t)
	defer s.Close()

	for _, path := range []string {"/atom.xml", "atom.xml", "/atom.xml?some=stuff", "atom.xml?some=stuff"} {
		t.Run(path, func(t *testing.T) {
			wa := doGet(t, s, path)

			assert.Equal(t, "application/atom+xml", wa.MimeType)
			assert.Equal(t, "<feed xmlns=\"http://www.w3.org/2005/Atom\">"+
				"<title>Entries for www.superblog.com</title>"+
				"<id>www.superblog.com</id>"+
				"<link href=\"http://www.superblog.com\"></link>"+
				"<updated>2017-10-28T13:45:00+02:00</updated>"+
				"<entry>"+
				"<title>Gurbai!</title><id>1509191100</id>"+
				"<link href=\"http://www.superblog.com/entry/201710281345_gurbai.md\"></link>"+
				"<published>2017-10-28T13:45:00+02:00</published><updated></updated>"+
				"<summary type=\"text/html\">&lt;p&gt;Gurbai!&lt;/p&gt;</summary>"+
				"</entry>"+
				"<entry>"+
				"<title>Hello my frens!</title>"+
				"<id>1506599100</id>"+
				"<link href=\"http://www.superblog.com/entry/201709281345_hello-my-frens.md\"></link>"+
				"<published>2017-09-28T13:45:00+02:00</published><updated></updated>"+
				"<summary type=\"text/html\">&lt;p&gt;&lt;img src=&#34;/static/img.png&#34; alt=&#34;Image&#34;/&gt;&lt;/p&gt;</summary>"+
				"</entry>"+
				"<entry>"+
				"<title>Hello guy!</title>"+
				"<id>1477655100</id>"+
				"<link href=\"http://www.superblog.com/entry/201610281345_hello_guy.md\"></link>"+
				"<published>2016-10-28T13:45:00+02:00</published><updated></updated>"+
				"<summary type=\"text/html\">&lt;p&gt;Paragraph of hello guy&lt;/p&gt;</summary>"+
				"</entry>"+
				"</feed>",
				string(wa.Body))
		})
	}
}

func TestFile(t *testing.T) {
	s := testServer(t)
	defer s.Close()

	type testCase struct {
		path string
		expectedBody string
		expectedMime string
	}
	for _, tc := range []testCase{
		{
			path: "static/style.css",
			expectedMime: "text/css; charset=utf-8",
			expectedBody: "h1 {color: red;}",
	}, {
			path: "static/text/foot.txt",
			expectedMime: "text/plain; charset=utf-8",
			expectedBody: "bar!",
		}, {
			path: "static/text/foot.txt?foo=bar",
			expectedMime: "text/plain; charset=utf-8",
			expectedBody: "bar!",
		},
	} {
		t.Run(tc.path, func(t *testing.T) {
			wa := doGet(t, s, tc.path)
			assert.Equal(t, tc.expectedMime, wa.MimeType)
			assert.Equal(t, tc.expectedBody, string(wa.Body))
		})
	}
}

func TestIndex(t *testing.T) {
	s := testServer(t)
	defer s.Close()

	wa := doGet(t, s, "")
	assert.Equal(t, "text/html; charset=utf-8", wa.MimeType)
	assert.Equal(t, strings.Trim(`
<h3><a href="/entry/201710281345_gurbai.md">Gurbai!</a></h3>
<p>Posted on October 28, 2017 at 13:45</p>
<p>Gurbai!</p>

<h3><a href="/entry/201709281345_hello-my-frens.md">Hello my frens!</a></h3>
<p>Posted on September 28, 2017 at 13:45</p>
<p><img src="/static/img.png" alt="Image"/></p>

<h3><a href="/entry/201610281345_hello_guy.md">Hello guy!</a></h3>
<p>Posted on October 28, 2016 at 13:45</p>
<p>Paragraph of hello guy</p>
`, " \n\r"), strings.Trim(string(wa.Body), " \n\r"))
}

// TODO: remove <html><head></head><body> and </body></html> from HTML generation
func TestEntryPage(t *testing.T) {
	s := testServer(t)
	defer s.Close()
	type testCase struct {
		path string
		expectedBody string
	}
	for _, tc := range []testCase{
		{
			path: "/entry/201610281345_hello_guy.md",
			expectedBody: `<h2>Hello guy!</h2>

<div>Posted on October 28, 2016 at 13:45</div>


<html><head></head><body>
<p>Paragraph of hello guy</p>
<p>This text is going to be ignored in the index.</p>
</body></html>`,
		},
		{
			path: "/entry/201710281345_gurbai.md",
			expectedBody: `<h2>Gurbai!</h2>

<div>Posted on October 28, 2017 at 13:45</div>


<html><head></head><body>
<p>Gurbai!</p>
</body></html>`,
		},
		{
			path: "/entry/gurbai.md",
			expectedBody: `<h2>Gurbai page!</h2>


<html><head></head><body>
<p>Gurbai!</p>
</body></html>`,
		},
	} {
		t.Run(tc.path, func(t *testing.T) {
			wa := doGet(t, s, tc.path)
			assert.Equal(t, "text/html; charset=utf-8", wa.MimeType)
			assert.Equal(t,
				strings.Trim(tc.expectedBody, " \n\r"),
				strings.Trim(string(wa.Body), " \n\r"))
		})
	}
}

// TODO: test404, testInternalServerError