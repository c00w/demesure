package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

func MakeLinodeRequest(action string, opts ...string) ([]byte, error) {
	api_key := os.Getenv("LINODE_API_KEY")
	v := url.Values{}
	v.Add("api_key", api_key)
	v.Add("api_action", action)
	for i, _ := range opts {
		if i%2 == 0 {
			continue
		}
		v.Add(opts[i-1], opts[i])
	}
	r, err := http.Get("https://api.linode.com/?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	resp, err := ioutil.ReadAll(r.Body)
	return resp, err
}

type DC struct {
	DataCenterId int
	Abbr         string
}

type dcRequest struct {
	ErrorArray []interface{}
	Action     string
	Data       []DC
}

func GetDCs() ([]DC, error) {
	output, err := MakeLinodeRequest("avail.datacenters")
	if err != nil {
		return nil, err
	}

	resp := &dcRequest{}

	err = json.Unmarshal(output, resp)
	if err != nil {
		return nil, err
	}

	if len(resp.ErrorArray) > 0 {
		return nil, errors.New(fmt.Sprint("Bad response", resp))
	}

	return resp.Data, nil
}

type linode struct {
	ErrorArray []interface{}
	Action     string
	Data       struct {
		LinodeID int64
	}
}

type linodediskcreate struct {
	ErrorArray []interface{}
	Action     string
	Data       struct {
		DiskID int64
	}
}

type linodeconfigcreate struct {
	ErrorArray []interface{}
	Action     string
	Data       struct {
		ConfigID int64
	}
}

type linodeboot struct {
	ErrorArray []interface{}
	Action     string
	Data       struct {
		JobID int64
	}
}

func CreateLinode(dc DC) (int64, error) {

	// Base Linode
	output, err := MakeLinodeRequest("linode.create", "DatacenterID", fmt.Sprint(dc.DataCenterId), "PlanID", "1")
	if err != nil {
		return 0, err
	}

	resp := &linode{}

	err = json.Unmarshal(output, resp)
	if err != nil {
		return 0, err
	}

	if len(resp.ErrorArray) > 0 {
		return 0, errors.New(fmt.Sprint("Bad response", resp))
	}

	id := resp.Data.LinodeID

	// Disk
	output, err = MakeLinodeRequest("linode.disk.createfromdistribution",
		"LinodeID", fmt.Sprint(id),
		"Label", "demesure-"+dc.Abbr+"-arch-disk",
		"DistributionID", "148",
		"Size", "20000",
		"rootPass", "ASDADASdasdasdersaddfsfvvffd12312312312^&*^&*",
		"rootSSHKey", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINZMnMpnL8YtHEaz7jlcscLa6ufspGv9vzNULkvcZIBS colin@rene",
	)

	diskresp := &linodediskcreate{}
	err = json.Unmarshal(output, diskresp)
	if err != nil {
		return 0, err
	}

	if len(diskresp.ErrorArray) > 0 {
		return 0, errors.New(fmt.Sprint("Bad response", diskresp))
	}

	diskid := diskresp.Data.DiskID

	// Swap
	output, err = MakeLinodeRequest("linode.disk.create",
		"LinodeID", fmt.Sprint(id),
		"Label", "demesure-"+dc.Abbr+"-arch-disk-swap",
		"Size", "480",
		"Type", "swap",
	)

	swapresp := &linodediskcreate{}
	err = json.Unmarshal(output, swapresp)
	if err != nil {
		return 0, err
	}

	if len(swapresp.ErrorArray) > 0 {
		return 0, errors.New(fmt.Sprint("Bad response", swapresp))
	}

	swapid := swapresp.Data.DiskID

	// Config
	output, err = MakeLinodeRequest("linode.config.create",
		"LinodeID", fmt.Sprint(id),
		"KernelID", "138",
		"Label", "demesure-"+dc.Abbr,
		"DiskList", fmt.Sprint(diskid)+","+fmt.Sprint(swapid),
		"helper_distro", "true",
	)

	configresp := &linodeconfigcreate{}
	err = json.Unmarshal(output, configresp)
	if err != nil {
		return 0, err
	}

	if len(configresp.ErrorArray) > 0 {
		return 0, errors.New(fmt.Sprint("Bad response", configresp))
	}

	configid := configresp.Data.ConfigID
	configid = configid

	// Boot
	output, err = MakeLinodeRequest("linode.boot",
		"LinodeID", fmt.Sprint(id),
		"ConfigID", fmt.Sprint(configid),
	)

	bootresp := &linodeboot{}
	err = json.Unmarshal(output, bootresp)
	if err != nil {
		return 0, err
	}

	if len(bootresp.ErrorArray) > 0 {
		return 0, errors.New(fmt.Sprint("Bad response", bootresp))
	}

	return resp.Data.LinodeID, nil
}

type linodelist struct {
	ErrorArray []interface{}
	Action     string
	Data       []struct {
		LinodeID int64
	}
}

func ListLinodes() ([]int64, error) {
	output, err := MakeLinodeRequest("linode.list")
	if err != nil {
		return nil, err
	}

	resp := &linodelist{}

	err = json.Unmarshal(output, resp)
	if err != nil {
		return nil, err
	}

	if len(resp.ErrorArray) > 0 {
		return nil, errors.New(fmt.Sprint("Bad response", resp))
	}

	out := []int64{}
	for _, l := range resp.Data {
		out = append(out, l.LinodeID)
	}
	return out, nil
}

type linodedelete struct {
	ErrorArray []interface{}
	Action     string
	Data       struct {
		LinodeID int64
	}
}

func DeleteLinode(id int64) error {
	output, err := MakeLinodeRequest("linode.delete", "LinodeID", fmt.Sprint(id), "skipChecks", "true")

	resp := &linodedelete{}

	err = json.Unmarshal(output, resp)
	if err != nil {
		return err
	}

	if len(resp.ErrorArray) > 0 {
		return errors.New(fmt.Sprint("Bad response", resp))
	}
	return nil
}

type linodeiplist struct {
	ErrorArray []interface{}
	Action     string
	Data       []struct {
		IPAddress string
		LinodeID  int64
	}
}

func GetAllIPs() ([]*net.IPAddr, error) {
	output, err := MakeLinodeRequest("linode.ip.list")

	resp := &linodeiplist{}

	err = json.Unmarshal(output, resp)
	if err != nil {
		return nil, err
	}

	if len(resp.ErrorArray) > 0 {
		return nil, errors.New(fmt.Sprint("Bad response", resp))
	}

	ips := map[int64]*net.IPAddr{}

	for _, v := range resp.Data {
		addr, err := net.ResolveIPAddr("ip", v.IPAddress)
		if err != nil {
			log.Printf("Error resolving", v.IPAddress)
			continue
		}
		ips[v.LinodeID] = addr
	}

	out := []*net.IPAddr{}

	for _, v := range ips {
		out = append(out, v)
	}
	return out, nil
}
