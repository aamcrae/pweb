# pweb
pweb is a generator for web based photo albums. It uses a [config](example/photos/web) file
that provides the parameters.

pweb uses a web-assembly program to dynamically generate the web pages
using XML data files as the source. There are 2 types of XML files used, album and gallery.
An album essentially points to a list of other albums and galleries, and is used
as a navigation menu. A gallery is a collection of images presented as a set of thumbnail pages;
separate images can be selected, viewed, downloaded etc.

EXIF tags and metadata is used to add titles, or selectively filter the list by rating etc.

[My own photo web site](http://amcrae.au/photos/index.html) is an example of an extensive web site
built by pweb.

The main goals for pweb is:
- Minimise the effort required to make a web gallery from a set of images.
- Take advantage of the hierarchial layout of albums and galleries.
- Use the available workflow tools to add ratings and captions.
- Easily maintain and update existing galleries.
- Optimise the speed of building galleries.

Using pweb, it is possible to take a set of images and insert these as a new gallery into an existing web site
by just a couple of simple commands. Maximum advantage is taken of concurrent image processing so that
creating or updating galleries can be very fast.

## Navigation

Web pages built by pweb are simple to navigate. For albums, a list of sub-albums and galleries is presented,
and clicking on the item will navigate to the selected album or gallery.
Clicking the title of the album will typically navigate back to the upper referencing album. A swipe down or
pressing the Up or Left arrow will also navigate back to the referencing album.

Galleries are where the images are displayed. A gallery is presented as one or more pages of thumbnails.
Each thumbnail may have a headline or caption displayed as part of the thumbnail. Thumbnail pages can be navigated
using keyboard shortcuts such as arrow keys. Home and End will jump to the first or last of the images.
If there are multiple pages of thumbnails, a page list will be displayed at the top - clicking the page number
will jump to that page of thumbnails.
Page Up and Page Down will also jump to the previous or next page of thumbnails.
Hitting Enter when an image is highlighted will display a full page view
of that image (clicking any of the thumbnails will also show the image).

Once the full sized image page is shown, right/left arrow keys will allow next/previous images to be displayed
(there are also link in the corners). A link at the top allows returning back to the thumbnail page
(pressing Up arrow will also return back to the thumbnail page).

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
                  |_ t
                  |_ p
                  |_ d
           australia
             |_ index.html
                gallery.xml
                  |_ t
                  |_ p
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

Gallery directories have subdirectories created for thumbnails (```t```), medium sized preview images (```p```)
and if configured, download (```d```). These directories are created or removed as necessary.

The [index.html](assets/index.html) file is used for both albums and galleries, and is
basically just used for loading the web assembly program.
When a gallery is created or updated, pweb will append the new gallery to the album
that is meant to reference the gallery, and generate the necessary gallery.xml file and
install the ```index.html``` file to the target web page directory.
Top level albums that reference other albums will need to be initially created (usually by copying
and modifying an existing ```album.xml``` file).

## Workflow

A typical workflow for using pweb to generate new galleries may be:

1. Take photos.
2. Process the photos using your favourite raw converter, adding titles and a rating.
3. Create a config file that contains all the appropriate information that pweb requires to generate the web pages (see below). Often it's easiest to cut'n'paste an existing config file.
4. Run ```pweb``` with the config file as an argument (if none is provided, the default file ```.web``` is used).

The web pages and resized image files are generated and placed in the location provided in the config file, and
the album referencing the gallery is updated.
```pweb``` may be run at any time using the same config file to update or regenerate the gallery, typically
if photos need to be added or removed from the gallery. The modification time on the images is used to check
whether the web images need to be regenerated.

The EXIF metadata on the images may be used to provide image headlines/captions, and the XMP Rating can be used
to filter the selected photos.

The XMP rating is a useful method of filtering images, and is supported by most raw conversion processors.
In the event that images are not being processed through a workflow that allows setting of
a XMP rating and captions, there is a separate program called [ptag](https://github.com/aamcrae/ptag) that can be used
to add a rating and a caption to an existing image.

## Image Selection

A key part of pweb is the selection of images to be used in the gallery being generated.
Images are selected in the config file via the ```include``` and ```exclude``` directives.
The order of images is maintained according to the config file, unless
date or name sorting is selected.

The process for selecting the final list of images is:
- Use one or more ```include``` directives in the config file to provide a list of image filenames. These filenames may be
wildcarded (along with brace expansion). If no ```include``` directives are configured, a default list is set (```*.jpg```).
- A similar set of ```exclude``` directives allows selected files to be excluded
- The ```after``` and ```before``` keywords allow wildcarded files to be added after or before a specific image.
- A ```rating``` and ```select``` directive allows filtering by XMP Rating values. ```rating``` operates as a value
(i.e all images rated that value or higher are included). ```select``` will only include images that have the matching rating
(multiple ratings values may be selected). ```select``` is useful when ratings are used to group images in separate categories.
Galleries can then be created with combinations of the categories.

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
All other config directives either have defaults or are optional.

The directives are:

| Keyword | Arguments | Example | Description |
|---------|-----------|---------|-------------|
| dir | directory-name | hiking/usa/yosemite | The ```dir``` keyword defines the directory where the generated web pages will be written. The directory is relative to the base web directory set in the ```pweb``` flags.|
| title | Gallery title | Yosemite Hiking | The title that is placed on the gallery. If no title is specified, "Photo Album" is used.|
| up | link to referring album | ../index.html | Indicates the album that is referencing this gallery. If set, the path is used to find the ```album.xml``` file that refers to this gallery, and a link is added to the album to this gallery (if none already exists). If this directive is not present, no change is made to any referring album, and no link back from this gallery is generated (this is useful to create a private or orphaned gallery, inaccessible from the main album navigation).|
| include | filenames | day-{2,3}/img_2*.jpg | A list of filenames (which may be wildcards) indicating the images to be included in this gallery. Multiple ```include``` lines may be used. If no ```include``` directives are present, the default include of ```*.jpg``` is used.|
| exclude | filenames | */img_234[5-7].jpg | A list of filenames that are to be excluded from the gallery. Multiple exclude lines are allowed.|
| after | file filenames | img_1234.jpg other/*.jpg | Insert the list of selected files after the file specified. This allows files to be placed in a particular order.|
| before | file filenames | img_4321.jpg other/*.jpg | Similar to ```after``` except the files are placed immediately before the file selected.|
| rating | 0 - 5 | 3 | Selects images that have a XMP rating this value or higher. Images that have XMP Rating metadata or with rating values less than the selected value are excluded.|
| select | 0 - 5 | 2 4 5| Selects images where the XMP rating matches one of the of rating values in the list. Only one of ```rating``` or ```select``` may be used, they are mutally exclusive.|
| download | static,symlink | | Allow the original images to be downloaded via a link in the generated web pages. Also, unless ```nozip``` is set, create a ```photos.zip``` file containing all of the photos in the gallery, and provide a link to download this zip file. No argument or ```symlink``` will use symlinks to the original. ```static``` will place a copy of the original image into the download directory.|
| nozip | | | If set, do not generate a ```photos.zip``` file for download.|
| sort | date,name | date | ```name``` will sort the images by their filename. ```date``` will sort the images by date. The date used is extracted from the EXIF of the image, or the modification time if no EXIF date is available. By default the images are placed in the order they are included.|
| reverse | | | If set, add the link to this gallery to the end of the list in the referring album; otherwise, the link to the gallery will be placed at the start of the album list. By default, album entries are considered to be newest first. By using ```reverse```, newer entries are placed at the end. Typically this is done when processing a set of galleries that are associated together, and the processing is done in chronological order (with the album entries also put in chronological order).
| caption | file title | img1234.jpg Nice flowers | Use this title string for the caption on the image; any EXIF captions are ignored.|
| large | | | If set, generate a larger image to be displayed for the image. Default image size is 1500 x 1200, large image size is 1800 x 1500.|
| nocaption | | | If set, do not generate captions for the images.|
| thumb | size | 200 | Set the width and height of the thumbnails generated to this value. The default is 160.|

## Flags

There are some runtime flags that may be used to control ```pweb```. These can be displayed
via ```pweb --help```. The main ones are:
- ```--force```: Remove the gallery completely and rebuild it. This is useful when changing the thumbnail size etc.
- ```--base```: Used to set the web pages base directory (default /var/www/html/photos).
- ```--assets```: Directory containing template web files such as the album and gallery ```index.html``` files etc. These can be locally customised (default /usr/share/pweb).

Other flags exist for various diagnostic functions.

## Initial installation

To install ```pweb```:
- Build and install the program via ```go build; sudo cp pweb /usr/local/bin```
- Copy the asset files to an appropriate location e.g
```
cp assets/*.xml assets/*.html /usr/share/pweb
cp assets/css/* /var/www/html/pweb
```
- Build and install the WASM support and binaries (optionally, tinygo can be used - see below):
```
(cd wasm; GOOS=js GOARCH=wasm go build -o /var/www/html/pweb/pweb.wasm)
cp $(go env GOROOT)/misc/wasm/wasm_exec.js /var/www/html/pweb
```
- The album-template.xml and gallery-template.xml files may be customised to add a copyright owner.

If image downloading with symlinks is configured, download links are generated as symlinks to the
original files, and your web server must be configured to allow the following of the symlinks
to access the file (e.g for apache2, ```Options FollowSymLinks``` must be set for the photos directory).
If download is configured as ```static```, a copy of the image is placed into the download directory, which
is useful if the web site is going to be copied/synced to a separate server.

## tinygo

The WASM binary built with the standard Go compiler is quite large, over 8.5Mb.
[Tinygo](https://tinygo.org/) can be used instead, reducing to size to under 1Mb.
To build the wasm binary with tinygo:
```
(cd wasm; tinygo build -target wasm -o /var/www/html/pweb/pweb.wasm)
cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js /var/www/html/pweb
```

It is important that the appropriate ```wasm_exec.js``` file is installed for whichever compiler is used.
