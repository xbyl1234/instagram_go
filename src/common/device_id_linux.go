package common

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/ioctl.h>
#include <net/if.h>
#include <unistd.h>

#include <errno.h>
#include <fcntl.h>
#include <linux/hdreg.h>

char * mac_serial(char *iface){
    int fd;
    struct ifreq ifr;
    unsigned char *mac = NULL;
    char *rtn=malloc(21);
	memset(rtn,21,0);
    errno=0;
    memset(&ifr, 0, sizeof(ifr));

    fd = socket(AF_INET, SOCK_DGRAM, 0);

    ifr.ifr_addr.sa_family = AF_INET;
    strncpy(ifr.ifr_name , iface , IFNAMSIZ-1);

	int suc =  ioctl(fd, SIOCGIFHWADDR, &ifr);
    close(fd);

    if (0 == suc) {
        mac = (unsigned char *)ifr.ifr_hwaddr.sa_data;
        sprintf(rtn,"%.2X:%.2X:%.2X:%.2X:%.2X:%.2X" , mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
		return rtn;
	}else{
		free(rtn);
  		return NULL;
	}
}
*/
import "C"
import (
	"os/exec"
	"unsafe"
)

//我的服务器
var DeviceID float64 = -682291321 + 2.1654324

func GetDeviceID() float64 {
	//mac, err := C.mac_serial(C.CString("ens33"))
	mac, err := C.mac_serial(C.CString("enp1s0"))
	if err != nil {
		//fmt.Printf("e1 %v\n", err)
	}
	defer C.free(unsafe.Pointer(mac))

	cmd := exec.Command("lsblk", "--nodeps", "-no", "serial", "/dev/cdrom")
	result, err := cmd.Output()

	goMac := C.GoString(mac)
	hash, _ := Hash([]byte(goMac + string(result)))
	fhash, _ := strconv.ParseFloat(hash, 64)
	return fhash
}
