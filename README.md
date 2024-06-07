# pweb
pweb is a generator for web based photo albums. It uses a [config](example/photos/web) file
that provides the parameters.

pweb uses web-assembly programs to dynamically generate the web pages
using XML data files as the source. There are 2 types of XML files used, album and gallery.
An album essentially points to a list of other albums and galleries, and is used
as a navigation menu. A gallery is a collection of images presented as a set of thumbnail pages;
separate images can be selected, viewed, downloaded etc.

EXIF tags and metadata is used to add titles, or selectively filter the list by rating etc.

[My own photo web site](http://amcrae.au/photos/index.html) is an example of an extensive web site
built by pweb.
## Navigation
Web pages built by pweb are simple to navigate. For albums, a list of sub-albums and galleries is presented,
and clicking on the item will navigate to the selected album or gallery.
Clicking the title of the album will typically navigate back to the upper referencing album.

Galleries are where the images are displayed. A gallery is presented as one or more pages of thumbnails.
Each thumbnail may have a headline or caption displayed as part of the thumbnail. Thumbnail pages can be navigated
using keyboard shortcuts such as arrow keys. Home and End will jump to the first or last of the images.
If there are multiple pages of thumbnails, a page list will be displayed at the top - clicking the page number
will jump to that page of thumbnails. Hitting Enter when an image is highlighted will display a full page view
of that image (clicking any of the thumbnails will also show the image).

Once the full sized image page is shown, right/left arrow keys will allow next/previous images to be displayed
(there are also link in the corners). A link at the top allows returning back to the thumbnail page.

For mobile devices, left/right swipes on the thumbnail page will go to the next/previous thumbnail pages.
When showing the full sized image, swipes will show the next or previous images.
Swiping down will return to the thumbnail page.

The web pages are responsive to window resizing, displaying more or less thumbnails as necessary.
## Directory Organisation

The intent is that photos are stored separately from the actual
web pages, and that pweb will perform all the necessary processing to
create a web gallery for the photos, resizing and copying them to the web site.

A typical web site directory layout may be:
```
  [base directory]
    |_ index.html
       album.xml
       travel
        |_ index.html
           album.xml
           usa
             |_ index.html
                gallery.xml
           australia
             |_ index.html
                gallery.xml
       family
        |_ index.html
           album.xml
           birthdays
             |_ index.html
                gallery.xml
...
```
Each directory either contains an album or a gallery. The directory layout
basically follows the navigation layout of the web site (clicking on albums or galleries will
step to that directory).

There are separate ```index.html``` files for
(albums)[assets/album-index.html] and (galleries)[assets/gallery-index.html],
with the only difference being the WASM file loaded (TODO: merge into a single binary).
When a gallery is created or updated, pweb will append the new gallery to the album
that is meant to reference the gallery, and generate the necessary gallery.xml file and install the appropriate
```index.html``` to the target web page directory.
Top level albums that reference other albums will need to be initially created (usually by copying an
existing album.xml file.

## Workflow

A typical workflow for using pweb to generate new galleries may be:

1. Take photos.
2. Process the photos using your favourite raw converter, adding titles and a rating.
3. Create a ```pweb``` config file that contains all the appropriate information that pweb requires to generate the web pages (see below). Often it's easy to copy a config file from another gallery.
4. Run ```pweb``` with the config file as an argument (if none is provided, the default file ```.web``` is used).

The web pages and resized image files are generated and placed in the location provided in the config file, and
the album referencing the gallery is updated.
```pweb``` may be run at any time using the same config file to update or regenerate the gallery, typically
if photos need to be added or removed from the gallery. The modification time on the images is used to identify
if the web images need to be regenerated.

The EXIF metadata on the images may be used to provide image headlines/captions, and the XMP Rating can be used
to filter the selected photos.
