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

func getObjectsList(serviceClient *gophercloud.ServiceClient, contaier string) error {
	options := &osObjects.ListOpts{Full: true}

	objects.List(serviceClient, contaier, options).EachPage(func(page pagination.Page) (bool, error) {
		nameList, err := objects.ExtractNames(page)
		if err != nil {
			return false, err
		}

		for _, name := range nameList {
			fmt.Printf("%+v\n", name)
		}

		return true, nil
	})

	return nil
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

	defer writer.Flush()

	// content, err := result.ExtractContent()
	// result.Body is also an io.ReadCloser of the file content that may be consumed as a stream.

	_, err := io.Copy(writer, dr.Body)

	if err != nil {
		return err
	}

	// fmt.Printf("Writed %d\n", bytes)

	return dr.Body.Close()
}

func getReaderForPath(path string) (io.ReadSeeker, error) {
	if path == "-" {
		return os.Stdin, nil
	}

	return os.Open(path)
}

func uploadObject(serviceClient *gophercloud.ServiceClient, container string, objectName string, reader io.ReadSeeker) error {
	res := objects.Create(serviceClient, container, objectName, reader, nil)

	return res.Err
}

type BaseFileCommand struct {
	File   flags.Filename `short:"f" required:"true" long:"file" default:"-"`
	Object string         `short:"o" required:"true" long:"object"`
}

type Options struct {
	UserName  string `short:"u" required:"true" long:"user-name" description:"Username"`
	ApiKey    string `short:"k" required:"true" long:"api-key" description:"Api key"`
	Region    string `short:"r" required:"false" long:"region" description:"Container region" default:"LON"`
	Container string `short:"c" required:"true" long:"container" description:"Container"`

	DownloadCommand struct {
		*BaseFileCommand
	} `command:"download" description:"Downaload object from storage"`

	UploadCommand struct {
		*BaseFileCommand
	} `command:"upload" description:"Upload object to storage"`

	ListCommand struct {
	} `command:"list" description:"Get list ob objects in storage"`
}

func doAction() error {
	var options Options
	parser := flags.NewParser(&options, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		return err
	}

	serviceClient, err := getServiceClient(options.UserName, options.ApiKey, options.Region)
	if err != nil {
		return err
	}

	switch parser.Active {
	case parser.Command.Find("list"):
		return getObjectsList(serviceClient, options.Container)

	case parser.Command.Find("upload"):
		reader, err := getReaderForPath(string(options.UploadCommand.File))
		if err != nil {
			return err
		}
		return uploadObject(serviceClient, options.Container, options.UploadCommand.Object, reader)

	case parser.Command.Find("download"):
		writer, err := getWriterForPath(string(options.DownloadCommand.File))
		if err != nil {
			return err
		}
		return downloadObject(serviceClient, options.Container, options.DownloadCommand.Object, writer)
	}

	return nil
}

func main() {
	err := doAction()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
