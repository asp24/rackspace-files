package main

import (
	"bufio"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/rackspace/gophercloud"
	osObjects "github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
	"io"
	"log"
	"os"
)

func getServiceClient(userName string, apiKey string, region string) (*gophercloud.ServiceClient, error) {
	ao := gophercloud.AuthOptions{
		Username: userName,
		APIKey:   apiKey,
	}

	provider, err := rackspace.AuthenticatedClient(ao)
	if err != nil {
		return nil, err
	}

	serviceClient, err := rackspace.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: region,
	})

	if err != nil {
		return nil, err
	}

	return serviceClient, nil
}

func getObjectsList(serviceClient *gophercloud.ServiceClient, contaier string) {
	options := &osObjects.ListOpts{Full: true}

	objects.List(serviceClient, contaier, options).EachPage(func(page pagination.Page) (bool, error) {
		nameList, err := objects.ExtractNames(page)
		if err != nil {
			return false, err
		}

		for _, name := range nameList {
			// ...
			fmt.Printf("%+v\n", name)
		}

		return true, nil
	})
}

func getWriterForPath(path string) (*bufio.Writer, error) {
	if path == "-" {
		writer := bufio.NewWriter(os.Stdout)

		return writer, nil
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(f)

	return writer, nil
}

func downloadObject(serviceClient *gophercloud.ServiceClient, container string, objectName string, writer *bufio.Writer) error {
	dr := objects.Download(serviceClient, container, objectName, nil)
	if dr.Err != nil {
		return dr.Err
	}

	// content, err := result.ExtractContent()
	// result.Body is also an io.ReadCloser of the file content that may be consumed as a stream.

	bytes, _ := io.Copy(writer, dr.Body)
	fmt.Printf("Writed %d\n", bytes)

	return dr.Body.Close()
}

func getReaderForPath(path string) (*bufio.Reader, error) {
	if path == "-" {
		reader := bufio.NewReader(os.Stdin)

		return reader, nil
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	return reader, nil
}

// func uploadObject() error {
//   //
// }

type BaseFileCommand struct {
	File   flags.Filename `short:"f" required:"true" long:"file" default:"-"`
	Object string         `short:"o" required:"true" long:"object"`
}

type DownloadCommand struct {
	*BaseFileCommand
}

func (x *DownloadCommand) Execute(args []string) error {
	//x.BaseFileCommand
	fmt.Printf("Adding (all=%+v): %#v\n", x.BaseFileCommand.Object, args)
	return nil
}

type UploadCommand struct {
	*BaseFileCommand
}

type ListCommand struct {
}

func (x *ListCommand) Execute(args []string) error {
	fmt.Printf("Adding (all=%v): %#v\n", x, args)
	return nil
}

var options struct {
	UserName  string `short:"u" required:"true" long:"user-name" description:"Username"`
	ApiKey    string `short:"k" required:"true" long:"api-key" description:"Api key"`
	Region    string `short:"r" required:"false" long:"region" description:"Container region" default:"LON"`
	Container string `short:"c" required:"true" long:"container" description:"Container"`

	DownloadCommand struct {
		*BaseFileCommand
	} `command:"download"`

	UploadCommand struct {
		*BaseFileCommand
	} `command:"upload"`

	ListCommand struct {
	} `command:"list"`
}

func main() {
	parser := flags.NewParser(&options, flags.Default)

	// var listCommand ListCommand
	// parser.AddCommand("list", "List objects", "List of objects in container", &listCommand)

	// var uploadCommand UploadCommand
	// parser.AddCommand("upload", "Upload file", "Upload file to RackSpace", &uploadCommand)

	// var downloadCommand DownloadCommand
	// parser.AddCommand("download", "Download file", "Download file from RackSpace", &downloadCommand)

	// Parse flags from `args'. Note that here we use flags.ParseArgs for
	// the sake of making a working example. Normally, you would simply use
	// flags.Parse(&opts) which uses os.Args

	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)

		panic(err)
	}

	fmt.Printf("=> %+v %+v\n", options, parser.Active)

	serviceClient, err := getServiceClient(options.UserName, options.ApiKey, options.Region)
	if err != nil {
		log.Fatal(err)
	}

	switch parser.Active {
	case parser.Command.Find("list"):
		getObjectsList(serviceClient, options.Container)
	case parser.Command.Find("upload"):

	case parser.Command.Find("download"):
		writer, err := getWriterForPath(string(options.DownloadCommand.File))
		if err != nil {
			log.Fatal(err)
		}
		downloadObject(serviceClient, options.Container, options.DownloadCommand.Object, writer)
	}
}
