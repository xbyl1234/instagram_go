package anti

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"net"
	"strconv"
	"unsafe"
)

//我的电脑
var DeviceID float64 = 6967533473 + 2.1654324

func GetDeviceID() float64 {
	var flag float64
	guid, err := getMachineGuid()
	if err == nil {
		hash, _ := Hash([]byte(guid))
		fhash, _ := strconv.ParseFloat(hash, 64)
		flag += flag + fhash
	}

	cpuinfo, err := getCPUInfo()
	if err == nil && len(cpuinfo) > 0 {
		hash, _ := Hash([]byte(cpuinfo[0].VendorID + cpuinfo[0].PhysicalID))
		fhash, _ := strconv.ParseFloat(hash, 64)
		flag += flag + fhash
	}

	mac, err := getMACAddress()
	if err == nil {
		hash, _ := Hash([]byte(mac))
		fhash, _ := strconv.ParseFloat(hash, 64)
		flag += flag + fhash
	}
	return flag
}

func getMACAddress() (string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", nil
	}
	mac, macerr := "", errors.New("无法获取到正确的MAC地址")
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags&net.FlagUp) != 0 && (netInterfaces[i].Flags&net.FlagLoopback) == 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				ipnet, ok := address.(*net.IPNet)
				if ok && ipnet.IP.IsGlobalUnicast() {
					mac = netInterfaces[i].HardwareAddr.String()
					return mac, nil
				}
			}
		}
	}
	return mac, macerr
}

type cpuInfo struct {
	CPU        int32  `json:"cpu"`
	VendorID   string `json:"vendorId"`
	PhysicalID string `json:"physicalId"`
}

type win32_Processor struct {
	Manufacturer string
	ProcessorID  *string
}

func getCPUInfo() ([]cpuInfo, error) {
	var ret []cpuInfo
	var dst []win32_Processor
	q := wmi.CreateQuery(&dst, "")
	if err := wmiQuery(q, &dst); err != nil {
		return ret, err
	}

	var procID string
	for i, l := range dst {
		procID = ""
		if l.ProcessorID != nil {
			procID = *l.ProcessorID
		}

		cpu := cpuInfo{
			CPU:        int32(i),
			VendorID:   l.Manufacturer,
			PhysicalID: procID,
		}
		ret = append(ret, cpu)
	}

	return ret, nil
}

// WMIQueryWithContext - wraps wmi.Query with a timed-out context to avoid hanging
func wmiQuery(query string, dst interface{}, connectServerArgs ...interface{}) error {
	ctx := context.Background()
	if _, ok := ctx.Deadline(); !ok {
		ctxTimeout, cancel := context.WithTimeout(ctx, 3000000000) //超时时间3s
		defer cancel()
		ctx = ctxTimeout
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- wmi.Query(query, dst, connectServerArgs...)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func getMachineGuid() (string, error) {
	// there has been reports of issues on 32bit using golang.org/x/sys/windows/registry, see https://github.com/shirou/gopsutil/pull/312#issuecomment-277422612
	// for rationale of using windows.RegOpenKeyEx/RegQueryValueEx instead of registry.OpenKey/GetStringValue
	var h windows.Handle
	err := windows.RegOpenKeyEx(windows.HKEY_LOCAL_MACHINE, windows.StringToUTF16Ptr(`SOFTWARE\Microsoft\Cryptography`), 0, windows.KEY_READ|windows.KEY_WOW64_64KEY, &h)
	if err != nil {
		return "", err
	}
	defer windows.RegCloseKey(h)

	const windowsRegBufLen = 74 // len(`{`) + len(`abcdefgh-1234-456789012-123345456671` * 2) + len(`}`) // 2 == bytes/UTF16
	const uuidLen = 36

	var regBuf [windowsRegBufLen]uint16
	bufLen := uint32(windowsRegBufLen)
	var valType uint32
	err = windows.RegQueryValueEx(h, windows.StringToUTF16Ptr(`MachineGuid`), nil, &valType, (*byte)(unsafe.Pointer(&regBuf[0])), &bufLen)
	if err != nil {
		return "", err
	}

	hostID := windows.UTF16ToString(regBuf[:])
	hostIDLen := len(hostID)
	if hostIDLen != uuidLen {
		return "", fmt.Errorf("HostID incorrect: %q\n", hostID)
	}

	return hostID, nil
}
