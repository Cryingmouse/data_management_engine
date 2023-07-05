package mgmtmodel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/driver"
	"golang.org/x/sync/errgroup"
)

type Host struct {
	ComputerName   string
	IP             string
	Username       string
	Password       string
	StorageType    string
	Caption        string
	OSArchitecture string
	Version        string
	BuildNumber    string
}

func (h *Host) Register() error {
	systemInfo, err := h.getSystemInfo()
	if err != nil {
		return err
	}

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	host := db.Host{
		IP:             h.IP,
		Username:       h.Username,
		Password:       h.Password,
		StorageType:    h.StorageType,
		ComputerName:   systemInfo.ComputerName,
		Caption:        systemInfo.Caption,
		OSArchitecture: systemInfo.OSArchitecture,
		Version:        systemInfo.Version,
		BuildNumber:    systemInfo.BuildNumber,
	}

	h.ComputerName = host.ComputerName
	h.Caption = host.Caption
	h.OSArchitecture = host.OSArchitecture
	h.Version = host.Version
	h.BuildNumber = host.BuildNumber

	return host.Save(engine)
}

func (h *Host) Unregister() error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	directoryList := db.DirectoryList{}

	filter := common.QueryFilter{
		Conditions: struct {
			HostIP string
		}{
			HostIP: h.IP,
		},
	}
	directoryList.Delete(engine, &filter)

	host := db.Host{
		IP: h.IP,
	}

	return host.Delete(engine)
}

func (h *Host) Get() (*Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	host := db.Host{
		ComputerName: h.ComputerName,
		IP:           h.IP,
	}
	if err = host.Get(engine); err != nil {
		return nil, err
	}

	common.CopyStructList(host, h)

	return h, nil
}

type HostList struct {
	Hosts []Host
}

func (hl *HostList) Get(filter *common.QueryFilter) ([]Host, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	hostList := db.HostList{}
	if err := hostList.Get(engine, filter); err != nil {
		return nil, err
	}

	common.CopyStructList(hostList.Hosts, &hl.Hosts)

	return hl.Hosts, nil
}

type PaginationHost struct {
	Hosts      []Host
	Page       int
	Limit      int
	TotalCount int64
}

func (hl *HostList) Pagination(filter *common.QueryFilter) (*PaginationHost, error) {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return nil, err
	}

	hostList := db.HostList{}
	paginationHosts, err := hostList.Pagination(engine, filter)
	if err != nil {
		return nil, err
	}

	paginationHostList := PaginationHost{
		Page:       filter.Pagination.Page,
		Limit:      filter.Pagination.PageSize,
		TotalCount: paginationHosts.TotalCount,
	}

	common.CopyStructList(paginationHosts.Hosts, &paginationHostList.Hosts)

	return &paginationHostList, nil
}

func (hl *HostList) Register() error {
	g, _ := errgroup.WithContext(context.Background())

	results := make([]*common.SystemInfo, len(hl.Hosts))
	var resultErr error

	for i, h := range hl.Hosts {
		index := i // 避免闭包问题
		host := h  // 避免闭包问题
		g.Go(func() error {
			systemInfo, err := host.getSystemInfo()
			if err != nil {
				resultErr = errors.Join(resultErr, err)
				return err
			}
			results[index] = systemInfo // 保存协程的返回值

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if resultErr != nil {
			return resultErr
		}
		return err
	}

	if err := common.CopyStructList(results, &hl.Hosts); err != nil {
		return err
	}

	dbHostList := db.HostList{}

	if err := common.CopyStructList(hl.Hosts, &dbHostList.Hosts); err != nil {
		return err
	}

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	return dbHostList.Save(engine)
}

func (hl *HostList) Unregister() error {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		return err
	}

	hostList := db.HostList{}

	common.CopyStructList(hl.Hosts, &hostList.Hosts)

	return hostList.Delete(engine)
}

func (h *Host) getSystemInfo() (*common.SystemInfo, error) {
	hostContext := common.HostContext{
		IP:       h.IP,
		Username: h.Username,
		Password: h.Password,
	}
	driver := driver.GetDriver(h.StorageType)

	response, err := driver.GetSystemInfo(hostContext)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	defer response.Body.Close()

	// 检查响应状态码
	if response.StatusCode != http.StatusOK {
		fmt.Printf("请求失败，状态码：%d\n", response.StatusCode)
		return nil, fmt.Errorf("Failed")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err, body)
		return nil, err
	}

	var result struct {
		Message    string            `json:"message"`
		SystemInfo common.SystemInfo `json:"system-info"`
	}

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return &result.SystemInfo, nil
}
