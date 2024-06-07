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

The main goals for pweb is:
- Minimise the effort required to make a web gallery from a set of images.
- Take advantage of the hierachial layout of albums and galleries.
- Use the available workflow tools to add ratings and captions.
- Easily maintain and update existing galleries.
- Optimise the speed of building galleries.

Using pweb, it is possible to take a set of images and insert these as a new gallery into an existing web site
by just a couple of simple commands. Maximum advantage is taken of concurrent image processing so that
creating or updating galleries can be very fast.
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
3. Create a config file that contains all the appropriate information that pweb requires to generate the web pages (see below). Often it's easiest to cut'n'paste an existing config file.
4. Run ```pweb``` with the config file as an argument (if none is provided, the default file ```.web``` is used).

The web pages and resized image files are generated and placed in the location provided in the config file, and
the album referencing the gallery is updated.
```pweb``` may be run at any time using the same config file to update or regenerate the gallery, typically
if photos need to be added or removed from the gallery. The modification time on the images is used to identify
if the web images need to be regenerated.

The EXIF metadata on the images may be used to provide image headlines/captions, and the XMP Rating can be used
to filter the selected photos.

## Image Selection

A key part of pweb is the selection of images to be used in the gallery being generated.
There are a number of ways that images can be selected via the config file.
The order of images is maintained according to the order they are selected in the config file, unless
date/time sorting is required.

The process for selecting the final list of images is:
- One or more ```include``` directives in the config file provides a list of image filenames. These filenames may be
wildcarded (along with brace expansion). If no ```include``` config exists, a default list is used as ```*.jpg```.
- A similar ```exclude``` set allows selected files to be excluded from the list.
- The ```after``` and ```before``` keywords allow wildcarded files to be added after or before a specific image resp. 
- A ```rating``` and ```select``` directive allows filtering by XMP Rating values.

## Config file

The config file is a series of lines, with each line containing a keyword followed by a ':', and then optional
arguments. Empty lines and lines starting with '#' are ignored.

A typical config file appears:
```
title: Day 2 & 3: Climbing Mount Kinabalu
dir: hiking/asia/kinabalu-2015/climb
include: day-2/*.jpg
include day-3/*.jpg
up: ../index.html
rating: 2
reverse:
```

There is only **one** required configuration keyword, ```dir```, that specifies the directory of
the web site (relative to the base directory) where the generated gallery is to be placed.
All other config directives have defaults or are optional.

The directives are:
| Keyword | Arguments | Example | Description |
|---------------------------------------------|
| dir | directory-name | hiking/usa/yosemite | The ```dir``` keyword defines the directory where the generated web pages will be written. The directory is relative to the base web directory set in the ```pweb``` flags.|
| title | Gallery title | Yosemite Hiking | The title that is placed on the gallery. If no title is specified, "Photo Album" is used.|
| up | link to referring album | ../index.html | Indicates the album that is referencing this gallery. If set, the path is used to find the ```album.xml``` file that refers to this gallery, and a link is added to the album to this gallery (if none already exists).|
| include | filenames | day-{2,3}/img_2*.jpg | A list of filenames (which may be wildcards) indicating the images to be included in this gallery. Multiple ```include``` lines may be used.|
