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

	/*	cdnClient, err := rackspace.NewObjectCDNV1(provider, gophercloud.EndpointOpts{
		Region: raxRegion,
	}) */

	opts := &osObjects.ListOpts{Full: true, Path: fontFolder}

	pager := osObjects.List(serviceClient, raxContainer, opts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		// Get a slice of strings, i.e. object names
		objectNames, err := osObjects.ExtractNames(page)
		if err != nil {
			log.Fatal(err)
		}
		for _, n := range objectNames {
			fmt.Println(n)
		}

		return true, nil
	})

	/*
		metadata := map[string]string{"some-key": "some-data"}
		_, err := objects.Update(
			serviceClient,
			"{containerName}",
			"{objectName}",
			objects.UpdateOpts{Metadata: metadata},
		).ExtractHeader()
	*/

}
