package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace"
	//"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/cdncontainers"
	//"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/cdnobjects"
	osObjects "github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	//	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	//"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
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

	EachObject(serviceClient, raxContainer, fontFolder, handlerListHeaders)
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
	osObjects.Update(
		obj.Client,
		obj.Container,
		obj.Name,
		osObjects.UpdateOpts{Metadata: metadata},
	)

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
