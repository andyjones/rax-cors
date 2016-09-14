rax-cors
========

Adds an `Access-Origin-Control` header to all objects
in a folder in Rackspace cloud files.

Synopsis
--------

List objects in the `fonts` folder in the `ws-styles-dev` container:

    ./rax-cors -container ws-styles-dev

Add an `Access-Origin-Control` header to all objects:

    ./rax-cors -container ws-styles-dev -update

Description
-----------

Browsers will only load the fonts from a different domain
if that domain includes an `Access-Origin-Control` header.
Rackspace Cloud Files supports custom headers but they get
wiped on every deploy because DeployBot uploads new objects.

After a deploy, run this script and it will add the correct
headers to fonts.

You can tell if the script needs to be run by checking
the network tab when you load https://www.weswap.com.
If it is full of red errors, you need to run this script!

The one off setup
-----------------

    # Build the binary
    go get
    go build

And create a `.env` file containing your credentials:

```
RAX_USERNAME="<your rackspace username>"
RAX_API_KEY="<your rackspace API key>"
```

You can also use environment variables if you prefer.
