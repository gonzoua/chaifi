// +build ignore

package main

import (
    "strings"
    "fmt"
    "bytes"
    "os/exec"
)

type Network struct {
    ssid, bssid, psk string
    wpa bool
}

func genNetworkEntry(network Network) string {
    result := ""
    result += "network={\n"
    if network.ssid != "" {
        result += fmt.Sprintf("    ssid=\"%s\"\n", network.ssid)
    } else if network.bssid != "" {
        result += fmt.Sprintf("    bssid=\"%s\"\n", network.bssid)
    }

    keyMgmt := "NONE"
    if network.wpa {
        keyMgmt = "WPA-PSK"
    }
    result += fmt.Sprintf("    key_mgmt=%s\n", keyMgmt)

    if network.psk != "" {
        result += fmt.Sprintf("    psk=\"%s\"\n", network.ssid)
    }

    result += "}\n"
    return result
}

func listScan(iface string) []Network {
    cmd := exec.Command("ifconfig", "-v", iface, "list", "scan")
    cmdOutput := &bytes.Buffer{}
    cmd.Stdout = cmdOutput
    err := cmd.Run()
    if err != nil {
        return nil
    }
    output := string(cmdOutput.Bytes())
    lines := strings.Split(output, "\n")
    if len(lines) < 1 {
        return nil
    }
    header := lines[0]
    lines = lines[1:]
    ssidEnd := strings.Index(header, "BSSID") - 1
    if ssidEnd < 0 {
        return nil
    }

    result := []Network{}
    for _, line := range lines {
        if len(line) < ssidEnd + 1 {
            continue
        }
        ssid := line[:ssidEnd]
        ssid = strings.Trim(ssid, " ")
        bssid := line[ssidEnd + 1:ssidEnd + 18]
        network := Network {ssid: ssid, bssid: bssid}
        result = append(result, network)

        wpaPos := strings.Index(line, "WPA<") - 1
        if wpaPos > 0 {
            network.wpa = true
        }
    }

    return result
}

func main() {
    networks := listScan("wlan0")
    fmt.Printf("Total networks: %d\n", len(networks))
    for _, net := range networks {
        fmt.Printf("===> [%s] %s\n", net.bssid, net.ssid)
    }

    fmt.Println(genNetworkEntry(networks[0]))
}