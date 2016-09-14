package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rackspace/gophercloud"
	osObjects "github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace"
	"log"
	"os"
)

func main() {
	// Load the credentials
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	raxUsername := os.Getenv("RAX_USERNAME")
	if raxUsername == "" {
		log.Fatal("Missing the Rackspace username")
	}
	raxApiKey := os.Getenv("RAX_API_KEY")
	if raxApiKey == "" {
		log.Fatal("Missing your Rackspace API Key, find it under `Account Settings`")
	}

	// Parse the command line args
	var raxRegion string
	flag.StringVar(&raxRegion, "region", "LON", "Rax region")

	var raxContainer string
	flag.StringVar(&raxContainer, "container", "", "the rackspace container")

	var fontFolder string
	flag.StringVar(&fontFolder, "fontfolder", "fonts", "folder that contains the fonts")

	var updateCORS bool
	flag.BoolVar(&updateCORS, "update", false, "Add the CORS headers to objects (default is to list the objects)")

	flag.Parse()

	if raxContainer == "" {
		log.Fatal("Specify the rackspace container")
	}

	// And now do something useful!
	ao := gophercloud.AuthOptions{
		Username: raxUsername,
		APIKey:   raxApiKey,
	}
	provider, err := rackspace.AuthenticatedClient(ao)

	serviceClient, err := rackspace.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: raxRegion,
	})

	handler := handlerListHeaders
	if updateCORS {
		handler = handlerAddCORSHeaders
	}

	EachObject(serviceClient, raxContainer, fontFolder, handler)
}

type Object struct { // A RAX cloud object
	Client    *gophercloud.ServiceClient
	Container string // name of the object container
	Name      string // name of the object including folders (eg. fonts/something.woff2)
}

func handlerListHeaders(obj *Object) (bool, error) {
	fmt.Println(obj.Name)
	metadata := osObjects.Get(obj.Client, obj.Container, obj.Name, &osObjects.GetOpts{})

	fmt.Println(metadata)

	return true, nil
}

func handlerAddCORSHeaders(obj *Object) (bool, error) {
	fmt.Println(obj.Name)

	metadata := map[string]string{"Access-Control-Allow-Origin": "*"}

	url := obj.Client.ServiceURL(obj.Container, obj.Name)

	// we need to add a header called "Access-Control-Allow-Origin"
	// but osObjects.Update() will prepend it with "X-Meta-" ...
	// so we implement our own hacky version
	resp, err := obj.Client.Request("POST", url, gophercloud.RequestOpts{
		MoreHeaders: metadata,
		OkCodes:     []int{202},
	})
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	return true, nil
}

func EachObject(c *gophercloud.ServiceClient, raxContainer string, folder string,
	handler func(*Object) (bool, error)) error {
	opts := &osObjects.ListOpts{Full: true, Path: folder}

	pager := osObjects.List(c, raxContainer, opts)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		objectNames, err := osObjects.ExtractNames(page)
		if err != nil {
			log.Fatal(err)
		}
		for _, objectName := range objectNames {
			cont, err := handler(&Object{
				Name:      objectName,
				Container: raxContainer,
				Client:    c,
			})
			if err != nil {
				return false, err
			}
			if !cont {
				return false, nil
			}
		}

		return true, nil
	})

	return err
}
