package vm

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/cppforlife/bosh-virtualbox-cpi/driver"
	bnet "github.com/cppforlife/bosh-virtualbox-cpi/vm/network"
)

const (
	// Attaching NICs to running VM is not allowed, so 4 NICs will always be connected.
	maxNICs = 4
)

type NICs struct {
	driver driver.Driver
	vmID   string
}

func (n NICs) Configure(nets Networks, host Host) (Networks, error) {
	if len(nets) > maxNICs {
		return nil, bosherr.Errorf("Exceeded maximum # of NICs (%d)", maxNICs)
	}

	nicIdx := 1
	netsWithMACs := Networks{}

	for netName, net := range nets { // todo there is no network order?
		mac, err := n.addNIC(strconv.Itoa(nicIdx), net, host)
		if err != nil {
			return nil, err
		}

		net.MAC = mac
		netsWithMACs[netName] = net

		nicIdx++
	}

	return netsWithMACs, nil
}

func (n NICs) addNIC(nic string, net Network, host Host) (string, error) {
	// http://www.virtualbox.org/manual/ch06.html#network_nat_service
	// https://www.virtualbox.org/ticket/6176
	// `VBoxManage setextradata VM_NAME "VBoxInternal/Devices/pcnet/0/LUN#0/Config/Network" "172.23.24/24"`
	// `VBoxManage setextradata VM_NAME "VBoxInternal/Devices/pcnet/0/LUN#0/Config/DNSProxy" 1`
	args := []string{"modifyvm", n.vmID, "--nic" + nic}

	switch net.CloudPropertyType {
	case bnet.NATType:
		args = append(args, []string{"nat"}...)

	case bnet.NATNetworkType:
		actualNet, err := host.FindNetwork(net)
		if err != nil {
			return "", err
		}
		args = append(args, []string{"natnetwork", "--nat-network" + nic, actualNet.Name()}...)

	case bnet.HostOnlyType:
		actualNet, err := host.FindNetwork(net)
		if err != nil {
			return "", err
		}
		args = append(args, []string{"hostonly", "--hostonlyadapter" + nic, actualNet.Name()}...)

	default:
		return "", bosherr.Errorf("Unknown network type: %s", net.CloudPropertyType)
	}

	mac, err := n.randomMAC()
	if err != nil {
		return "", err
	}

	args = append(args, []string{"--macaddress" + nic, strings.ToUpper(fmt.Sprintf("%02x", mac))}...)

	_, err = n.driver.Execute(args...)

	return n.userFriendly(mac), err
}

func (NICs) randomMAC() ([]byte, error) {
	// http://stackoverflow.com/questions/21018729/generate-mac-address-in-go
	buf := make([]byte, 6)

	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	// VirtualBox uses '[0-9A-Fa-f][02468ACEace][0-9A-Fa-f]{10}' to validate MACs
	// Also set local bit, ensure unicast address
	buf[0] = 2

	return buf, nil
}

func (NICs) userFriendly(buf []byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}
