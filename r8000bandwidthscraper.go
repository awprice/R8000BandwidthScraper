package main

import (
    "log"
    "github.com/PuerkitoBio/goquery"
    "strings"
    "io"
    "strconv"
    "net/http"
    "os"
    "fmt"
    "os/exec"
    "runtime"
)

type Device struct {
    Connection      string
    Name            string
    Type            string
    SSID            string
    MAC             string
    Download        string
    DownloadPercent float64
    Upload          string
    UploadPercent   float64
    IP              string
}

var clear map[string]func()

func main() {
    log.SetFlags(0)

    if len(os.Args) != 3 {
        panic("Not enough arguments")
    }
    var host string = os.Args[1]
    var key string = os.Args[2]

    message, onMessage := make(chan *http.Response), make(chan string)
    go func() {
        // Create a new client
        client := &http.Client{}
        req, err := http.NewRequest("GET", "http://" + host + "/DEV_device_iQoS.htm", nil)
        if err != nil {
            panic(err.Error())
        }
        req.Header.Set("Authorization", "Basic " + key)
        for {
            resp, err := client.Do(req)
            if err != nil {
                panic(err.Error())
            }
            onMessage <-"Message ready!"
            message <-resp
        }
    }()

    for {
        select {
        case <-onMessage:
            go func() {
                var respMessage *http.Response = <-message
                devices, err := parseHtml(respMessage.Body)
                if err != nil {
                    panic(err.Error())
                }
                printDevices(devices)
            }()
        }
    }
}

func init() {
    clear = make(map[string]func())
    clear["linux"] = func() {
        cmd := exec.Command("clear")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
    clear["darwin"] = func() {
        cmd := exec.Command("clear")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
    clear["windows"] = func() {
        cmd := exec.Command("cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
}

func CallClear() {
    value, ok := clear[runtime.GOOS]
    if ok {
        value()
    } else {
        panic("Your platform is unsupported! I can't clear terminal screen :(")
    }
}

func parseHtml(file io.Reader) ([]Device, error) {
    var devices []Device

    doc, err := goquery.NewDocumentFromReader(file)
    if err != nil {
        return devices, err
    }

    doc.Find("table.sortable").Each(func(i int, s *goquery.Selection) {
        s.Find("tr").Each(func(i int, s *goquery.Selection) {
            if !s.HasClass("table_header") {
                var tempDevice Device
                s.Find("td").Each(func(i int, s *goquery.Selection) {
                    switch i {
                    case 3:
                        tempDevice.Connection = s.Text()
                    case 4:
                        s.Find("table").Each(func(i int, s *goquery.Selection) {
                            // Get the name and IP
                            nameIP, err := s.Find("span").Html()
                            if err != nil {
                                panic(err.Error())
                            }
                            nameIPSlice := strings.Split(nameIP, "<br/>")
                            tempDevice.Name = nameIPSlice[0]
                            if tempDevice.Name == "--" {
                                tempDevice.Name = "(NO NAME)"
                            }
                            tempDevice.IP = nameIPSlice[1]

                            // Get the mac, type and ssid
                            title, _ := s.Attr("title")
                            lines := strings.Split(title, "\n")
                            for _, line := range lines {
                                if strings.Contains(line, "MAC Address") {
                                    tempDevice.MAC = strings.TrimSpace(strings.Replace(line, "MAC Address: ", "", -1))
                                }
                                if strings.Contains(line, "Device Type") {
                                    tempDevice.Type = strings.TrimSpace(strings.Replace(line, "Device Type: ", "", -1))
                                }
                                if strings.Contains(line, "SSID") {
                                    tempDevice.SSID = strings.TrimSpace(strings.Replace(line, "SSID: ", "", -1))
                                }
                            }
                        })
                    case 7:
                        tempDevice.Download = s.Find("p span").Text()
                        tempDevice.DownloadPercent, _ = strconv.ParseFloat(s.Find("span.dev-hide").Text(), 64)
                    case 8:
                        tempDevice.Upload = s.Find("p span").Text()
                        tempDevice.UploadPercent, _ = strconv.ParseFloat(s.Find("span.dev-hide").Text(), 64)
                    }
                })
                if len(tempDevice.Name) != 0 {
                    devices = append(devices, tempDevice)
                }
            }
        })
    })

    return devices, nil
}

func printDevices(devices []Device) {
    CallClear()
    for i := 0; i < len(devices); i++ {
        fmt.Println(devices[i].Name + " (" + devices[i].IP + ")" + ": ↓ " + devices[i].Download + " / ↑ " + devices[i].Upload)
    }
}